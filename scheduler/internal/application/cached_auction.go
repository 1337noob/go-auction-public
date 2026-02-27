package application

import "time"

type CachedAuction struct {
	ID      string
	Timeout time.Duration
}
