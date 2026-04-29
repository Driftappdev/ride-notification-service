package mqtt_client

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
)

// Client interface สำหรับ NotificationService
type Client interface {
	SendRaw(ctx context.Context, topic string, data []byte) error
	SendProto(ctx context.Context, topic string, msg proto.Message) error
	Close()
}

// MQTTClient implements Client
type MQTTClient struct {
	client mqtt.Client
	qos    byte
}

// NewMQTTClient สร้าง MQTT client
// broker: tcp://127.0.0.1:1883 หรือ tls://127.0.0.1:8883
// clientID: unique client ID
// qos: 0,1,2
// useTLS: true=เปิด TLS
func NewMQTTClient(broker, clientID string, qos byte, useTLS bool) *MQTTClient {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectTimeout(5 * time.Second)

	if useTLS {
		opts.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[MQTTClient] Failed to connect to broker %s: %v", broker, token.Error())
	}

	return &MQTTClient{
		client: c,
		qos:    qos,
	}
}

// SendRaw ส่ง byte[] ไป MQTT topic ใดก็ได้
func (m *MQTTClient) SendRaw(ctx context.Context, topic string, data []byte) error {
	const maxRetries = 3
	var err error

	for i := 0; i < maxRetries; i++ {
		token := m.client.Publish(topic, m.qos, false, data)

		done := make(chan struct{})
		go func() {
			token.Wait()
			close(done)
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			if token.Error() == nil {
				log.Printf("[MQTTClient] Message published to topic %s", topic)
				return nil
			}
			err = token.Error()
			log.Printf("[MQTTClient] Retry %d/%d failed for topic %s: %v", i+1, maxRetries, topic, err)
		}

		time.Sleep(500 * time.Millisecond)
	}

	return err
}

// SendProto ส่ง Protobuf message ไป MQTT topic
func (m *MQTTClient) SendProto(ctx context.Context, topic string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return m.SendRaw(ctx, topic, data)
}

// Close ปิด connection
func (m *MQTTClient) Close() {
	if m.client.IsConnected() {
		m.client.Disconnect(250)
	}
}
