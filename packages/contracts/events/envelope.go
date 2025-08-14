package events

import "time"

// Envelope is the standard event wrapper.
type Envelope[T any] struct {
	EventName      string    `json:"eventName"`
	EventVersion   string    `json:"eventVersion"`
	EventID        string    `json:"eventId"`
	OccurredAt     time.Time `json:"occurredAt"`
	CorrelationID  string    `json:"correlationId"`
	CausationID    string    `json:"causationId"`
	IdempotencyKey string    `json:"idempotencyKey"`
	Producer       string    `json:"producer"`
	TenantID       string    `json:"tenantId"`
	Payload        T         `json:"payload"`
}

