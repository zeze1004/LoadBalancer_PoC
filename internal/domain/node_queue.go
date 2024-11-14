package domain

import (
	"github.com/zeze1004/LoadBanlence_PoC/config/clouds"
	"github.com/zeze1004/LoadBanlence_PoC/internal/dto"
	"sync"
)

type NodeQueue struct {
	CloudService *clouds.CloudService
	Nodes        []*NodeAI
	Mutex        sync.Mutex // 큐가 노드에게 요청을 전달할 때 동시성 문제를 해결하기 위한 뮤텍스
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

func (nq *NodeQueue) SendRequestToNode(request *dto.Request) {
	nq.Mutex.Lock()
	defer nq.Mutex.Unlock()
	// TODO: 요청을 처리하는 로직 추가
}
