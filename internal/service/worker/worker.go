package service

import (
	"context"
	"log"
	"time"

	driverpb "dift_backend_go/notification-service/proto/pb/notification_driverpb" // notification_to_driver
	orderpb "dift_backend_go/notification-service/proto/pb/notification_orderpb"   // notification_to_order
	ui "dift_backend_go/notification-service/proto/pb/notification_uipb"           // notification_ui

	"google.golang.org/protobuf/proto"
)

// ============================
// StatusUpdater interface
// ============================
type StatusUpdater interface {
	UpdateDriverStatus(ctx context.Context, routeId, driverId string, status ui.DriverStatus)
	UpdatePassengerStatus(ctx context.Context, routeId string, status ui.PassengerStatus, driver *ui.PassengerNotification_DriverInfo)
}

// ============================
// NotificationWorker สำหรับส่ง message
// ============================
type NotificationWorker struct {
	uiUpdater      StatusUpdater
	driverProducer interface {
		Publish(topic, key string, value []byte) error
	}
	orderProducer interface {
		Publish(topic, key string, value []byte) error
	}
	driverTopic string
	orderTopic  string
	retryCount  int
	retryDelay  time.Duration
	timeout     time.Duration
}

// NewNotificationWorker สร้าง worker
func NewNotificationWorker(
	uiUpdater StatusUpdater,
	driverProducer interface {
		Publish(topic, key string, value []byte) error
	},
	orderProducer interface {
		Publish(topic, key string, value []byte) error
	},
	driverTopic, orderTopic string,
) *NotificationWorker {
	return &NotificationWorker{
		uiUpdater:      uiUpdater,
		driverProducer: driverProducer,
		orderProducer:  orderProducer,
		driverTopic:    driverTopic,
		orderTopic:     orderTopic,
		retryCount:     3,
		retryDelay:     500 * time.Millisecond,
		timeout:        5 * time.Second,
	}
}

// ============================
// Worker functions (non-blocking)
// ============================

// SendToDriverUI ส่งไป Driver UI แบบ non-blocking
func (w *NotificationWorker) SendToDriverUI(routeId, driverId string, status ui.DriverStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()
	go w.uiUpdater.UpdateDriverStatus(ctx, routeId, driverId, status)
}

// SendToPassengerUI ส่ง notification ไป Passenger UI แบบ non-blocking
func (w *NotificationWorker) SendToPassengerUI(routeId string, driver *ui.PassengerNotification_DriverInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()
	go w.uiUpdater.UpdatePassengerStatus(ctx, routeId, ui.PassengerStatus_DRIVER_ASSIGNED, driver)
}

// SendToDriverService ส่ง event ไป Driver Service แบบ non-blocking
func (w *NotificationWorker) SendToDriverService(routeId, driverId string) {
	event := &driverpb.DriverAssignmentEvent{
		RouteId:  routeId,
		DriverId: driverId,
		Status:   "assigned",
	}
	go w.publishDriverEventWithRetry(w.driverTopic, routeId, event)
}

// SendToOrderService ส่ง event ไป Order Service แบบ non-blocking
func (w *NotificationWorker) SendToOrderService(routeId, driverId string) {
	event := &orderpb.OrderUpdateEvent{
		OrderId:   routeId, // ใช้ OrderId แทน RouteId
		Status:    "matched",
		DriverId:  driverId,
		Timestamp: time.Now().Unix(),
	}
	go w.publishOrderEventWithRetry(w.orderTopic, routeId, event)
}

// ============================
// Internal helper functions
// ============================

// marshal & publish event to Driver Service with retry
func (w *NotificationWorker) publishDriverEventWithRetry(topic, key string, event *driverpb.DriverAssignmentEvent) {
	data, err := proto.Marshal(event)
	if err != nil {
		log.Printf("[NotificationWorker] Failed to marshal DriverAssignmentEvent: %v", err)
		return
	}

	for i := 0; i < w.retryCount; i++ {
		if err := w.driverProducer.Publish(topic, key, data); err != nil {
			log.Printf("[NotificationWorker] Attempt %d: Failed to send to Driver Service: %v", i+1, err)
			time.Sleep(w.retryDelay)
			continue
		}
		log.Printf("[NotificationWorker] Published DriverAssignmentEvent successfully: routeId=%s", key)
		return
	}
	log.Printf("[NotificationWorker] All retries failed for Driver Service event: routeId=%s", key)
}

// marshal & publish event to Order Service with retry
func (w *NotificationWorker) publishOrderEventWithRetry(topic, key string, event *orderpb.OrderUpdateEvent) {
	data, err := proto.Marshal(event)
	if err != nil {
		log.Printf("[NotificationWorker] Failed to marshal OrderUpdateEvent: %v", err)
		return
	}

	for i := 0; i < w.retryCount; i++ {
		if err := w.orderProducer.Publish(topic, key, data); err != nil {
			log.Printf("[NotificationWorker] Attempt %d: Failed to send to Order Service: %v", i+1, err)
			time.Sleep(w.retryDelay)
			continue
		}
		log.Printf("[NotificationWorker] Published OrderUpdateEvent successfully: orderId=%s", key)
		return
	}
	log.Printf("[NotificationWorker] All retries failed for Order Service event: orderId=%s", key)
}
