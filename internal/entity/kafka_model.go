package entity

import (
	"time"
)

type OrderEvent struct {
	EventType string    `json:"event_type"`
	Payload   *Order    `json:"payload"`
	SentAt    time.Time `json:"sent_at"`
}
