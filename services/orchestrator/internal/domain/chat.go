package domain

import "time"

type ChatMessage struct {
	ID        int
	SessionID string
	Role      string
	Content   string
	CreatedAt time.Time
}
