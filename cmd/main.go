package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dift_backend_go/notification-service/config"
	httpadapter "dift_backend_go/notification-service/internal/adapter/http"
	"dift_backend_go/notification-service/internal/integration/event"
	mqtt "dift_backend_go/notification-service/internal/integration/mqtt_client"
	"dift_backend_go/notification-service/internal/service"
	"dift_backend_go/notification-service/internal/servicecore"
	eventing "dift_backend_go/notification-service/internal/servicecore/eventing"
)

func main() {
	cfg := config.LoadConfig()
	health := servicecore.HealthController("notification-service", "v2-generic")
	log.Printf("service=%s version=%s status=%s", health.Service, health.Version, health.Status)

	rawEventBusConsumer := event.NewEventBusConsumer(cfg.EventBusBrokers, cfg.EventBusGroup)
	eventBusConsumer, err := eventing.NewReliableConsumer(rawEventBusConsumer)
	if err != nil {
		log.Fatalf("init reliable event consumer failed: %v", err)
	}
	rawEventPublisher := event.NewEventBusPublisher(cfg.EventBusBrokers)
	eventPublisher, err := eventing.NewReliablePublisher(rawEventPublisher)
	if err != nil {
		log.Fatalf("init reliable event publisher failed: %v", err)
	}
	mqttClient := mqtt.NewMQTTClient(cfg.MQTTBroker, "notification-service", 1, cfg.MQTTSecure)

	svc := service.NewGenericNotificationService(cfg, eventBusConsumer, mqttClient, eventPublisher)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go svc.StartListening(ctx)

	mux := http.NewServeMux()
	h := httpadapter.NewGenericHandler(svc.Dispatch)
	h.Register(mux)

	httpSrv := &http.Server{
		Addr:              ":" + itoa(cfg.HTTPPort),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	go func() {
		log.Printf("[notification-service] http listening on %s", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[notification-service] http server error: %v", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = httpSrv.Shutdown(shutdownCtx)

	done := make(chan struct{})
	go func() {
		_ = eventPublisher.Close()
		_ = eventBusConsumer.Close()
		mqttClient.Close()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		log.Println("timeout waiting for shutdown")
	}
}

func itoa(v int) string {
	if v <= 0 {
		return "2222"
	}
	return fmt.Sprintf("%d", v)
}
