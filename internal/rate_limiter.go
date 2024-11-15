package internal

import (
	"github.com/zeze1004/LoadBanlence_PoC/internal/model"
	"sync"
	"time"
)

type RateLimiter struct {
	NodeLimits map[string]*NodeLimit    // 노드 ID별 속도 제한 관리
	Nodes      map[string]*model.AInode // 노드 ID별 노드 정보 관리
	Mutex      sync.Mutex
}

type NodeLimit struct {
	CurrentBPM int
	CurrentRPM int
	LastReset  time.Time // 마지막 초기화 시간
	BPM        int
	RPM        int
}

func NewRateLimiter(nodes []*model.AInode) *RateLimiter {
	nodeLimits := make(map[string]*NodeLimit)
	nodeMap := make(map[string]*model.AInode)

	for _, node := range nodes {
		nodeLimits[node.ID] = &NodeLimit{
			BPM:       node.LimitBPM,
			RPM:       node.LimitRPM,
			LastReset: time.Now(),
		}
		nodeMap[node.ID] = node
	}

	return &RateLimiter{
		NodeLimits: nodeLimits,
		Nodes:      nodeMap,
	}
}

// AllowRequest 요청 허용 여부 확인
func (rl *RateLimiter) AllowRequest(nodeID string, requestSize int) bool {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	// 노드 정보와 속도 제한 정보 확인
	node, nodeExists := rl.Nodes[nodeID]
	limit, limitExists := rl.NodeLimits[nodeID]

	// 노드가 존재하지 않거나 비활성 상태이면 요청 불가
	if !nodeExists || !limitExists || !node.IsActive {
		return false
	}

	// 1분 경과 시 속도 제한 초기화
	now := time.Now()
	if now.Sub(limit.LastReset) >= time.Minute {
		limit.CurrentBPM = 0
		limit.CurrentRPM = 0
		limit.LastReset = now
	}

	// LimitBPM 및 LimitRPM 제한 확인
	if limit.CurrentBPM+requestSize > limit.BPM || limit.CurrentRPM+1 > limit.RPM {
		return false
	}

	// 요청 허용: 카운터 갱신
	limit.CurrentBPM += requestSize
	limit.CurrentRPM++
	return true
}
