package fcm

import (
	"strconv"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/driftappdev/libpackage/notification/domain"
)

func toMessage(n domain.Notification, token string, cfg Config) *messaging.Message {
	ttl := n.TTL
	if ttl <= 0 { ttl = cfg.DefaultTTL }

	msg := &messaging.Message{
		Token: token,
		Data:  n.Data,
		Notification: &messaging.Notification{
			Title:    n.Title,
			Body:     n.Body,
			ImageURL: n.Image,
		},
		Android: &messaging.AndroidConfig{
			Priority:    androidPriority(n.Priority),
			TTL:         &ttl,
			CollapseKey: n.CollapseKey,
			Notification: &messaging.AndroidNotification{
				ChannelID: cfg.DefaultAndroidChannel,
				ImageURL:  n.Image,
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": apnsPriority(n.Priority),
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: n.Title,
						Body:  n.Body,
					},
					Sound: sound(n.Priority),
				},
			},
		},
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: n.Title,
				Body:  n.Body,
				Image: n.Image,
			},
			Headers: map[string]string{
				"TTL": strconv.Itoa(int(ttl.Seconds())),
			},
		},
	}
	return msg
}

func toTopicOrConditionMessage(n domain.Notification, cfg Config) *messaging.Message {
	msg := toMessage(n, "", cfg)
	msg.Token = ""
	switch n.Audience.Type {
	case domain.AudienceTopic:
		msg.Topic = n.Audience.Topic
	case domain.AudienceCondition:
		msg.Condition = n.Audience.Condition
	}
	return msg
}

func androidPriority(p domain.Priority) string {
	if p >= domain.PriorityHigh { return "high" }
	return "normal"
}

func apnsPriority(p domain.Priority) string {
	if p >= domain.PriorityHigh { return "10" }
	return "5"
}

func sound(p domain.Priority) string {
	if p >= domain.PriorityHigh { return "default" }
	return ""
}

func ttlOrDefault(ttl time.Duration, def time.Duration) time.Duration {
	if ttl > 0 { return ttl }
	return def
}
