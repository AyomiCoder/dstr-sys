package notification

import (
	"context"
	"time"
)

type Notification struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Recipient string    `json:"recipient"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateInput struct {
	Type      string
	Recipient string
	Message   string
}

type Repository interface {
	Create(ctx context.Context, input CreateInput) (Notification, error)
	List(ctx context.Context) ([]Notification, error)
}
