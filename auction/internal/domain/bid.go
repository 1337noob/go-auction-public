package domain

import "time"

type Bid struct {
	ID        string
	UserID    string
	Amount    int
	CreatedAt time.Time
}
