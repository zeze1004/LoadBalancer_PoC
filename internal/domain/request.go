package domain

import "time"

type RequestEvent struct {
	RequestID   string
	RequestSize int
	Timestamp   time.Time
}
