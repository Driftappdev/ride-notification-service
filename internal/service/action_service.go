package service

import (
	"context"
	"errors"
	"log"
	"time"

	"dift_backend_go/notification-service/config"
	"dift_backend_go/notification-service/internal/integration/event"
	mqtt "dift_backend_go/notification-service/internal/integration/mqtt_client"

	eventpb "dift_backend_go/notification-service/proto/pb/notification_eventpb"
	orderpb "dift_backend_go/notification-service/proto/pb/notification_orderpb"
	ui "dift_backend_go/notification-service/proto/pb/notification_uipb"

	"google.golang.org/protobuf/proto"
)

type NotificationService struct {
	cfg           *config.Config
	eventClient   event.Client
	mqttClient    mqtt.Client
	orderProducer event.Publisher
}

func NewNotificationService(
	cfg *config.Config,
	eventClient event.Client,
	mqttClient mqtt.Client,
	orderProducer event.Publisher,
) *NotificationService {
	return &NotificationService{
		cfg:           cfg,
		eventClient:   eventClient,
		mqttClient:    mqttClient,
		orderProducer: orderProducer,
	}
}

func (s *NotificationService) StartListening(ctx context.Context) {
	log.Println("[NotificationService] Listening for matching events...")
	s.eventClient.Subscribe(ctx, s.cfg.EventBusMatchingTopic, s.handleMatchingEvent)
}

func toUILocation(loc *eventpb.Location) *ui.Location {
	if loc == nil {
		return nil
	}
	return &ui.Location{Lat: loc.Lat, Lng: loc.Lng, Address: loc.Address}
}

func (s *NotificationService) handleMatchingEvent(eventBytes []byte) {
	var evt eventpb.MatchingEvent
	if err := proto.Unmarshal(eventBytes, &evt); err != nil {
		log.Printf("[NotificationService] Failed to unmarshal MatchingEvent: %v", err)
		return
	}

	driverID := ""
	driverName := ""
	driverAvatar := ""
	driverCarType := ""
	driverCarPlate := ""
	if evt.Driver != nil {
		driverID = evt.Driver.DriverId
		driverName = evt.Driver.Name
		driverAvatar = evt.Driver.AvatarUrl
		driverCarType = evt.Driver.CarType
		driverCarPlate = evt.Driver.CarPlate
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if evt.Status == "no_driver" || driverID == "" {
		// update order immediately so UI can refresh with not_matched/no_driver.
		go s.publishOrderEvent(s.cfg.EventBusOrderTopic, evt.RouteId, &orderpb.OrderUpdateEvent{
			OrderId:   evt.RouteId,
			Status:    "not_matched",
			Timestamp: time.Now().Unix(),
		})
		log.Printf("[NotificationService] No driver routeId=%s", evt.RouteId)
		return
	}

	go s.sendDriverUI(ctx, &ui.DriverNotification{
		RouteId:    evt.RouteId,
		UserId:     evt.UserId,
		DriverId:   driverID,
		DriverName: driverName,
		Status:     ui.DriverStatus_WAITING_FOR_ACCEPT,
		Pickup:     toUILocation(evt.PickupLocation),
		Dropoff:    toUILocation(evt.DropoffLocation),
		Price:      evt.Price,
		Timestamp:  evt.Timestamp,
	})

	// also update order with candidate assignment state.
	go s.publishOrderEvent(s.cfg.EventBusOrderTopic, evt.RouteId, &orderpb.OrderUpdateEvent{
		OrderId:                 evt.RouteId,
		Status:                  "driver_assigned",
		DriverId:                driverID,
		DriverName:              driverName,
		DriverCarModel:          "",
		DriverAvatarUrl:         driverAvatar,
		CarPlate:                driverCarPlate,
		CarType:                 driverCarType,
		PickupLat:               evt.GetPickupLocation().GetLat(),
		PickupLng:               evt.GetPickupLocation().GetLng(),
		PickupAddress:           evt.GetPickupLocation().GetAddress(),
		DropoffLat:              evt.GetDropoffLocation().GetLat(),
		DropoffLng:              evt.GetDropoffLocation().GetLng(),
		DropoffAddress:          evt.GetDropoffLocation().GetAddress(),
		DistancePickupToDropoff: evt.DistancePickupToDropoff,
		DurationTotal:           evt.DurationTotalSec,
		RoutePolyline:           evt.RoutePolyline,
		Price:                   evt.Price,
		Timestamp:               evt.Timestamp,
	})

	log.Printf("[NotificationService] Assignment distributed routeId=%s driverId=%s", evt.RouteId, driverID)
}

func (s *NotificationService) handleDriverCancelEvent(eventBytes []byte) {
	var evt eventpb.CancellationEvent
	if err := proto.Unmarshal(eventBytes, &evt); err != nil {
		log.Printf("[NotificationService] Failed to unmarshal DriverCancelEvent: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go s.sendPassengerUI(ctx, &ui.PassengerNotification{
		RouteId: evt.RouteId,
		Status:  ui.PassengerStatus_DRIVER_ASSIGNED,
	})

	go s.publishOrderEvent(s.cfg.EventBusOrderTopic, evt.RouteId, &orderpb.OrderUpdateEvent{
		OrderId:   evt.RouteId,
		Status:    "driver_cancelled",
		DriverId:  evt.CancelledBy,
		Timestamp: time.Now().Unix(),
	})

	log.Printf("[NotificationService] Driver cancellation propagated routeId=%s", evt.RouteId)
}

func (s *NotificationService) DriverAcceptedEvent(ctx context.Context, evt *eventpb.DriverAcceptedEvent) error {
	if evt.RouteId == "" || evt.DriverId == "" {
		return errors.New("invalid DriverAcceptedEvent: missing routeId or driverId")
	}

	go s.sendPassengerUI(ctx, &ui.PassengerNotification{
		RouteId: evt.RouteId,
		Status:  ui.PassengerStatus_DRIVER_ASSIGNED,
		Driver: &ui.PassengerNotification_DriverInfo{
			Id:        evt.Driver.DriverId,
			Name:      evt.Driver.Name,
			AvatarUrl: evt.Driver.AvatarUrl,
			CarType:   evt.Driver.CarType,
			CarPlate:  evt.Driver.CarPlate,
		},
	})

	go s.publishOrderEvent(s.cfg.EventBusOrderTopic, evt.RouteId, &orderpb.OrderUpdateEvent{
		OrderId:         evt.RouteId,
		Status:          "matched",
		DriverId:        evt.Driver.DriverId,
		DriverName:      evt.Driver.Name,
		DriverAvatarUrl: evt.Driver.AvatarUrl,
		CarPlate:        evt.Driver.CarPlate,
		CarType:         evt.Driver.CarType,
		Timestamp:       time.Now().Unix(),
	})

	return nil
}

func (s *NotificationService) sendDriverUI(ctx context.Context, msg *ui.DriverNotification) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("[NotificationService] Failed to marshal DriverNotification: %v", err)
		return
	}
	if err := s.mqttClient.SendRaw(ctx, s.cfg.MQTTTopicDriver, data); err != nil {
		log.Printf("[NotificationService] Failed to send Driver UI: %v", err)
	}
}

func (s *NotificationService) sendPassengerUI(ctx context.Context, msg *ui.PassengerNotification) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("[NotificationService] Failed to marshal PassengerNotification: %v", err)
		return
	}
	if err := s.mqttClient.SendRaw(ctx, s.cfg.MQTTTopicPassenger, data); err != nil {
		log.Printf("[NotificationService] Failed to send Passenger UI: %v", err)
	}
}

func (s *NotificationService) publishOrderEvent(topic, key string, msg *orderpb.OrderUpdateEvent) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("[NotificationService] Failed to marshal OrderUpdateEvent: %v", err)
		return
	}
	if err := s.orderProducer.Publish(topic, key, data); err != nil {
		log.Printf("[NotificationService] Failed to publish Order Service: %v", err)
	}
}

func (s *NotificationService) UpdateDriverStatus(ctx context.Context, routeID, driverID string, status ui.DriverStatus) {
	go s.sendDriverUI(ctx, &ui.DriverNotification{
		RouteId:   routeID,
		DriverId:  driverID,
		Status:    status,
		Timestamp: time.Now().Unix(),
	})
}

func (s *NotificationService) UpdatePassengerStatus(ctx context.Context, routeID string, status ui.PassengerStatus, driver *ui.PassengerNotification_DriverInfo) {
	go s.sendPassengerUI(ctx, &ui.PassengerNotification{
		RouteId: routeID,
		Status:  status,
		Driver:  driver,
	})
}
