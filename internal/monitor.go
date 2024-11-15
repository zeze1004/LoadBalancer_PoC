package internal

import (
	"fmt"
	"github.com/zeze1004/LoadBanlence_PoC/internal/model"
	"net/http"
	"time"
)

type NodeMonitor struct {
	Nodes []*model.AInode
}

func NewNodeMonitor(nodes []*model.AInode) *NodeMonitor {
	return &NodeMonitor{Nodes: nodes}
}

// StartMonitoring 주기적으로 노드 헬스체크 진행
func (nm *NodeMonitor) StartMonitoring(interval time.Duration) {
	go func() {
		for {
			for _, node := range nm.Nodes {
				nm.checkNodeStatus(node)
			}
			time.Sleep(interval)
		}
	}()
}

func (nm *NodeMonitor) checkNodeStatus(node *model.AInode) {
	resp, err := http.Get(node.URL + "/health")
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// 비활성 처리
		if node.IsActive {
			fmt.Printf("Node %s 는 비활성 상태입니다\n", node.ID)
			node.IsActive = false
		}
	} else {
		// 활성 처리
		if !node.IsActive {
			fmt.Printf("Node %s 는 활성 상태입니다\n", node.ID)
			node.IsActive = true
		}
	}
	node.LastHealthCheck = time.Now()
}
