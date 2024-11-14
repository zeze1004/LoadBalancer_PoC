package domain

import "github.com/zeze1004/LoadBanlence_PoC/config/clouds"

type NodeAI struct {
	CloudService *clouds.CloudService
	BPM          int
	RPM          int
}

func NewNodeAI(cloudService *clouds.CloudService, BPM, RPM int) *NodeAI {
	return &NodeAI{
		CloudService: cloudService,
		BPM:          BPM,
		RPM:          RPM,
	}
}
