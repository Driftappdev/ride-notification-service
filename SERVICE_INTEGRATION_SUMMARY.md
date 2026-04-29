# Notification Service Summary

## Purpose

Enterprise-grade notification orchestration service for:
- generic notifications
- status/event update notifications
- multi-channel fan-out across multiple providers

Note:
- This service is configured for generic notifications only.
- Ride-hailing realtime notifications are handled by `notification-ride-service`.

## Inbound

- HTTP:
  - `POST /api/v1/notifications/send`
  - `POST /api/v1/notifications/event`
- EventBus consumer topics:
  - `event_bus.generic_topic`
  - `event_bus.event_topics[]`

## Outbound Providers

- EventBus dispatch provider
  - topic: `event_bus.dispatch_topic`
- MQTT provider
  - topics: `mqtt.driver_topic`, `mqtt.passenger_topic`, `mqtt.generic_topic`
- FCM provider (optional)
  - endpoint + server key from config
- Webhook provider (optional)
  - URL + bearer token from config

## Delivery Status Tracking

Per-channel delivery results are published to:
- `event_bus.status_topic`

Each result includes:
- `notification_id`
- `provider`
- `channel`
- `success/error`
- timestamp

## Event/Status Template Support

If `title`/`message` are missing, templates are auto-generated from:
- `event_type`
- `status`

Examples:
- `order_status:completed`
- `driver_status:offline`
- `payment_status:failed`
