package fcm

import (
	"context"
	"errors"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/driftappdev/libpackage/notification/domain"
	"github.com/driftappdev/libpackage/notification/port"
	"google.golang.org/api/option"
)

type Client struct {
	cfg Config
	msg *messaging.Client
}

var _ port.Sender = (*Client)(nil)

func New(ctx context.Context, cfg Config) (*Client, error) {
	cfg = cfg.Normalize()

	var opts []option.ClientOption
	if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: cfg.ProjectID}, opts...)
	if err != nil {
		return nil, err
	}

	msg, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{cfg: cfg, msg: msg}, nil
}

func NewWithMessagingClient(cfg Config, msg *messaging.Client) *Client {
	return &Client{cfg: cfg.Normalize(), msg: msg}
}

func (c *Client) Send(ctx context.Context, n domain.Notification) (domain.DeliveryResult, error) {
	if err := n.Validate(); err != nil {
		return failed(n, "VALIDATION_ERROR", err, false), err
	}
	if n.Channel != domain.ChannelFCM {
		err := errors.New("fcm client only supports fcm channel")
		return failed(n, "UNSUPPORTED_CHANNEL", err, false), err
	}

	switch n.Audience.Type {
	case domain.AudienceToken:
		return c.sendOne(ctx, n, n.Audience.Token)
	case domain.AudienceTokens:
		return c.sendTokens(ctx, n, n.Audience.Tokens)
	case domain.AudienceTopic, domain.AudienceCondition:
		return c.sendTopicOrCondition(ctx, n)
	default:
		err := errors.New("audience must be resolved before fcm sender")
		return failed(n, "UNRESOLVED_AUDIENCE", err, false), err
	}
}

func (c *Client) sendOne(ctx context.Context, n domain.Notification, token string) (domain.DeliveryResult, error) {
	messageID, err := c.msg.SendDryRun(ctx, toMessage(n, token, c.cfg))
	if !c.cfg.DryRun {
		messageID, err = c.msg.Send(ctx, toMessage(n, token, c.cfg))
	}
	if err != nil {
		cl := ClassifyError(err)
		return domain.DeliveryResult{
			NotificationID: n.ID,
			Provider:       "fcm",
			Status:         domain.DeliveryFailed,
			SentAt:         time.Now().UTC(),
			ErrorCode:      cl.Code,
			ErrorMessage:   err.Error(),
			Retryable:      cl.Retryable,
			TokenResults: []domain.TokenResult{{
				Token: token, OK: false, ErrorCode: cl.Code,
				ErrorMessage: err.Error(), Retryable: cl.Retryable, RemoveToken: cl.RemoveToken,
			}},
		}, err
	}

	return domain.DeliveryResult{
		NotificationID: n.ID,
		Provider:       "fcm",
		Status:         domain.DeliverySent,
		MessageID:      messageID,
		SentAt:         time.Now().UTC(),
		TokenResults: []domain.TokenResult{{
			Token: token, OK: true, MessageID: messageID,
		}},
	}, nil
}

func (c *Client) sendTokens(ctx context.Context, n domain.Notification, tokens []string) (domain.DeliveryResult, error) {
	now := time.Now().UTC()
	var all []domain.TokenResult
	var okCount int
	var retryable bool
	var firstErr error

	for _, batch := range chunk(tokens, c.cfg.MaxTokensPerBatch) {
		for _, token := range batch {
			res, err := c.sendOne(ctx, n, token)
			if err != nil && firstErr == nil { firstErr = err }
			if res.Retryable { retryable = true }
			for _, tr := range res.TokenResults {
				if tr.OK { okCount++ }
				all = append(all, tr)
			}
		}
	}

	status := domain.DeliveryFailed
	if okCount == len(tokens) { status = domain.DeliverySent }
	if okCount > 0 && okCount < len(tokens) { status = domain.DeliveryPartial }

	return domain.DeliveryResult{
		NotificationID: n.ID,
		Provider:       "fcm",
		Status:         status,
		SentAt:         now,
		TokenResults:   all,
		Retryable:      retryable,
	}, firstErr
}

func (c *Client) sendTopicOrCondition(ctx context.Context, n domain.Notification) (domain.DeliveryResult, error) {
	msg := toTopicOrConditionMessage(n, c.cfg)
	messageID, err := c.msg.SendDryRun(ctx, msg)
	if !c.cfg.DryRun {
		messageID, err = c.msg.Send(ctx, msg)
	}
	if err != nil {
		cl := ClassifyError(err)
		return failed(n, cl.Code, err, cl.Retryable), err
	}
	return domain.DeliveryResult{
		NotificationID: n.ID,
		Provider:       "fcm",
		Status:         domain.DeliverySent,
		MessageID:      messageID,
		SentAt:         time.Now().UTC(),
	}, nil
}

func failed(n domain.Notification, code string, err error, retryable bool) domain.DeliveryResult {
	msg := ""
	if err != nil { msg = err.Error() }
	return domain.DeliveryResult{
		NotificationID: n.ID,
		Provider:       "fcm",
		Status:         domain.DeliveryFailed,
		SentAt:         time.Now().UTC(),
		ErrorCode:      code,
		ErrorMessage:   msg,
		Retryable:      retryable,
	}
}

func chunk[T any](in []T, size int) [][]T {
	if size <= 0 { size = 500 }
	var out [][]T
	for len(in) > 0 {
		n := size
		if len(in) < n { n = len(in) }
		out = append(out, in[:n])
		in = in[n:]
	}
	return out
}
