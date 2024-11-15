package domain

import "time"

type AInode struct {
	ID              string
	BPM             int
	RPM             int
	isActive        bool
	LastHealthCheck time.Time
}

func NewAInode(ID string, BPM, RPM int) *AInode {
	return &AInode{
		ID:       ID,
		BPM:      BPM,
		RPM:      RPM,
		isActive: true,
	}
}
