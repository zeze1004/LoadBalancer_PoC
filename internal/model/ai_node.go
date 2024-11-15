package model

import "time"

type AInode struct {
	ID              string
	URL             string
	LimitBPM        int
	LimitRPM        int
	IsActive        bool
	LastHealthCheck time.Time
}
