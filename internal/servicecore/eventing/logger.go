package eventing

import (
	"log"

	"dift_backend_go/notification-service/internal/servicecore/queue"
)

type queueLogger struct {
	component string
}

func newQueueLogger(component string) queue.Logger {
	return &queueLogger{component: component}
}

func (l *queueLogger) Info(msg string, args ...any) {
	log.Printf("[%s] %s %v", l.component, msg, args)
}

func (l *queueLogger) Warn(msg string, args ...any) {
	log.Printf("[%s] WARN %s %v", l.component, msg, args)
}

func (l *queueLogger) Error(msg string, args ...any) {
	log.Printf("[%s] ERROR %s %v", l.component, msg, args)
}
