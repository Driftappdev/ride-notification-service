package service

import (
	"context"
	"log"
	"time"

	ui "dift_backend_go/notification-service/proto/pb/notification_uipb"

	"google.golang.org/protobuf/proto"
)

// ============================
// StatusUpdater interface สำหรับจัดการอัปเดตสถานะ UI
// ============================
type StatusUpdater interface {
	UpdateDriverStatus(ctx context.Context, routeId, driverId string, status ui.DriverStatus)
	UpdatePassengerStatus(ctx context.Context, routeId string, status ui.PassengerStatus, driver *ui.PassengerNotification_DriverInfo)
}

// ============================
// UIStatusService implements StatusUpdater
// ============================
type UIStatusService struct {
	mqttClient interface {
		SendRaw(ctx context.Context, topic string, data []byte) error
	}
	topicDriverUI    string
	topicPassengerUI string
}

// NewUIStatusService สร้าง instance ของ UIStatusService
func NewUIStatusService(
	mqttClient interface {
		SendRaw(ctx context.Context, topic string, data []byte) error
	},
	topicDriverUI, topicPassengerUI string,
) *UIStatusService {
	return &UIStatusService{
		mqttClient:       mqttClient,
		topicDriverUI:    topicDriverUI,
		topicPassengerUI: topicPassengerUI,
	}
}

// UpdateDriverStatus ส่งสถานะไป Driver UI
func (u *UIStatusService) UpdateDriverStatus(ctx context.Context, routeId, driverId string, status ui.DriverStatus) {
	driverUI := &ui.DriverNotification{
		RouteId:   routeId,
		DriverId:  driverId,
		Status:    status,
		Timestamp: time.Now().Unix(),
	}

	data, err := proto.Marshal(driverUI)
	if err != nil {
		log.Printf("[UIStatusService] Failed to marshal DriverNotification: %v", err)
		return
	}

	// ส่งไป topic ที่ config กำหนด
	if err := u.mqttClient.SendRaw(ctx, u.topicDriverUI, data); err != nil {
		log.Printf("[UIStatusService] Failed to send DriverNotification: %v", err)
	}
}

// UpdatePassengerStatus ส่งสถานะไป Passenger UI
func (u *UIStatusService) UpdatePassengerStatus(ctx context.Context, routeId string, status ui.PassengerStatus, driver *ui.PassengerNotification_DriverInfo) {
	passengerUI := &ui.PassengerNotification{
		RouteId: routeId,
		Status:  status,
		Driver:  driver,
	}

	data, err := proto.Marshal(passengerUI)
	if err != nil {
		log.Printf("[UIStatusService] Failed to marshal PassengerNotification: %v", err)
		return
	}

	// ส่งไป topic ที่ config กำหนด
	if err := u.mqttClient.SendRaw(ctx, u.topicPassengerUI, data); err != nil {
		log.Printf("[UIStatusService] Failed to send PassengerNotification: %v", err)
	}
}
