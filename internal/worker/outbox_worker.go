package worker

import (
	"context"
	"log"
	"time"

	"ss-catalog-service/internal/domain"
)

// MessageBroker defines the contract for publishing events.
type MessageBroker interface {
	Publish(ctx context.Context, eventType string, payload []byte) error
}

type OutboxWorker struct {
	outboxRepo domain.OutboxRepository
	broker     MessageBroker
	interval   time.Duration
}

func NewOutboxWorker(repo domain.OutboxRepository, broker MessageBroker, interval time.Duration) *OutboxWorker {
	if interval == 0 {
		interval = 5 * time.Second
	}
	return &OutboxWorker{
		outboxRepo: repo,
		broker:     broker,
		interval:   interval,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("Outbox worker started with interval %v", w.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Outbox worker shutting down...")
			return
		case <-ticker.C:
			w.processEvents(ctx)
		}
	}
}

func (w *OutboxWorker) processEvents(ctx context.Context) {
	// Fetch a batch of pending events
	events, err := w.outboxRepo.FetchPending(ctx, 50)
	if err != nil {
		log.Printf("Error fetching pending events: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	log.Printf("Processing %d outbox events", len(events))

	for _, event := range events {
		// Check context for graceful shutdown mid-batch
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Publish to broker
		if err := w.broker.Publish(ctx, event.EventType, event.Payload); err != nil {
			log.Printf("Failed to publish event ID %d: %v", event.ID, err)
			
			if err := w.outboxRepo.MarkAsFailed(ctx, event.ID, err.Error()); err != nil {
				log.Printf("Failed to mark event ID %d as failed: %v", event.ID, err)
			}
			
			// Best Practice: If broker is down, stop processing this batch to avoid I/O waste
			log.Printf("Broker might be down, skipping rest of the batch...")
			return
		}

		// Mark as published
		if err := w.outboxRepo.MarkAsPublished(ctx, event.ID); err != nil {
			log.Printf("Failed to mark event ID %d as published: %v", event.ID, err)
		}
	}
}
