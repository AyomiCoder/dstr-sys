package notification

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (s *PostgresRepository) Create(ctx context.Context, input CreateInput) (Notification, error) {
	const query = `
INSERT INTO notifications (type, recipient, message)
VALUES ($1, $2, $3)
RETURNING id, type, recipient, message, status, created_at;`

	var n Notification
	if err := s.pool.QueryRow(ctx, query, input.Type, input.Recipient, input.Message).Scan(
		&n.ID,
		&n.Type,
		&n.Recipient,
		&n.Message,
		&n.Status,
		&n.CreatedAt,
	); err != nil {
		return Notification{}, fmt.Errorf("insert notification: %w", err)
	}

	return n, nil
}

func (s *PostgresRepository) List(ctx context.Context) ([]Notification, error) {
	const query = `
SELECT id, type, recipient, message, status, created_at
FROM notifications
ORDER BY id DESC;`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query notifications: %w", err)
	}
	defer rows.Close()

	notifications := make([]Notification, 0)
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.Type, &n.Recipient, &n.Message, &n.Status, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan notification row: %w", err)
		}
		notifications = append(notifications, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notification rows: %w", err)
	}

	return notifications, nil
}
