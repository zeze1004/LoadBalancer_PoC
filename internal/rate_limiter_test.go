package internal

import (
	"github.com/zeze1004/LoadBanlence_PoC/internal/model"
	"testing"
	"time"
)

func TestRateLimiter_AllowRequest(t *testing.T) {
	t.Run("노드가 LimitBPM 초과했을 때", func(t *testing.T) {
		nodes := []*model.AInode{
			{ID: "node1", LimitBPM: 1000, LimitRPM: 10, IsActive: true},
		}
		rateLimiter := NewRateLimiter(nodes)

		nodeLimit := rateLimiter.NodeLimits["node1"]
		nodeLimit.LastReset = time.Now()

		// 첫 번째 요청: 900 바이트
		allowed1 := rateLimiter.AllowRequest("node1", 900)
		if !allowed1 {
			t.Errorf("node1은 요청 처리가 허락되지만 거부되었습니다")
		}

		// 두 번째 요청: 200 바이트 (총 1100 바이트, LimitBPM 초과)
		allowed2 := rateLimiter.AllowRequest("node1", 200)
		if allowed2 {
			t.Errorf("node1은 LimitBPM을 초과했으나, 요청 처리가 허락됐습니다")
		}
	})

	t.Run("노드가 LimitRPM 초과했을 때", func(t *testing.T) {
		nodes := []*model.AInode{
			{ID: "node1", LimitBPM: 1000, LimitRPM: 3, IsActive: true},
		}
		rateLimiter := NewRateLimiter(nodes)

		// LimitRPM 요청 반복
		for i := 0; i < 3; i++ {
			allowed := rateLimiter.AllowRequest("node1", 100)
			if !allowed {
				t.Errorf("node1은 요청 처리가 허락되지만 거부되었습니다")
			}
		}

		// 네 번째 요청은 LimitRPM 초과
		allowed := rateLimiter.AllowRequest("node1", 100)
		if allowed {
			t.Errorf("node1은 LimitRPM을 초과했으나, 요청 처리가 허락됐습니다")
		}
	})
}
