package model

import "time"

type RequestEvent struct {
	RequestID   string
	RequestSize int
	Timestamp   time.Time
	ResponseCh  chan *ResponseEvent
}

type ResponseEvent struct {
	RequestID string
	IsSuccess bool
}
