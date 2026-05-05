package messaging

import (
	"context"
	"log"
	"ss-catalog-service/internal/worker"
)

type logBroker struct{}

func NewLogBroker() worker.MessageBroker {
	return &logBroker{}
}

func (b *logBroker) Publish(ctx context.Context, eventType string, payload []byte) error {
	log.Printf(" [BROKER] Publishing event: %s | Payload: %s", eventType, string(payload))
	return nil
}
