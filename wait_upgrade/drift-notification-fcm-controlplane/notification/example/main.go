package main

import (
	"context"
	"fmt"
	"time"

	"github.com/driftappdev/libpackage/notification/controlplane"
	"github.com/driftappdev/libpackage/notification/domain"
	"github.com/driftappdev/libpackage/notification/fcm"
	"github.com/driftappdev/libpackage/notification/memory"
	"github.com/driftappdev/libpackage/notification/port"
)

func main() {
	ctx := context.Background()

	fcmClient, err := fcm.New(ctx, fcm.Config{
		ProjectID:       "your-firebase-project-id",
		CredentialsFile: "service-account.json",
		DryRun:          true,
	})
	if err != nil {
		panic(err)
	}

	router := controlplane.NewStaticRouter(map[domain.Channel]port.Sender{
		domain.ChannelFCM: fcmClient,
	})

	cp := controlplane.New(
		controlplane.Config{},
		router,
		controlplane.WithIdempotency(memory.NewIdempotencyStore()),
		controlplane.WithRateLimiter(memory.NewBucketLimiter(1000, 1000, time.Minute)),
	)

	res, err := cp.Dispatch(ctx, domain.Notification{
		ID:             "notif_001",
		TenantID:       "tenant_a",
		IdempotencyKey: "order_123:paid_push",
		Channel:        domain.ChannelFCM,
		Priority:       domain.PriorityHigh,
		Audience:       domain.Audience{Type: domain.AudienceToken, Token: "device-token"},
		Title:          "Order paid",
		Body:           "Your order has been paid successfully.",
		Data:           map[string]string{"order_id": "123"},
	})
	fmt.Println(res, err)
}
