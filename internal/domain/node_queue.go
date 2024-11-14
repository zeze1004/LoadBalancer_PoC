package domain

import "github.com/zeze1004/LoadBanlence_PoC/config/clouds"

type NodeQueue struct {
	CloudService *clouds.CloudService
	Nodes        []*NodeAI
	BPM          int
	RPM          int
}

func NewNodeQueue(cloudService *clouds.CloudService, nodes []*NodeAI, BPM, RPM int) *NodeQueue {
	return &NodeQueue{
		CloudService: cloudService,
		Nodes:        nodes,
		BPM:          BPM,
		RPM:          RPM,
	}
}
