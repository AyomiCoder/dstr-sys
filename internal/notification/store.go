package notification

import (
	"fmt"
	"sync"
	"time"
)

type Notification struct {
	ID        string    `json:"id"`
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

type Store struct {
	mu            sync.RWMutex
	sequence      uint64
	notifications []Notification
}

func NewStore() *Store {
	return &Store{
		notifications: make([]Notification, 0),
	}
}

func (s *Store) Create(input CreateInput) Notification {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sequence++
	notification := Notification{
		ID:        fmt.Sprintf("ntf-%06d", s.sequence),
		Type:      input.Type,
		Recipient: input.Recipient,
		Message:   input.Message,
		Status:    "queued",
		CreatedAt: time.Now().UTC(),
	}

	s.notifications = append(s.notifications, notification)
	return notification
}

func (s *Store) List() []Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]Notification, len(s.notifications))
	copy(items, s.notifications)
	return items
}
