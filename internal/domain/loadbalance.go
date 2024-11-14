package domain

import (
	"github.com/zeze1004/LoadBanlence_PoC/internal/dto"
)

type LoadBalancer struct {
	Queues    []*NodeQueue
	waitQueue chan *dto.Request
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		Queues:    []*NodeQueue{},
		waitQueue: make(chan *dto.Request, 10), // 대기 큐 버퍼 사이즈 10
	}
}

// 적합한 큐 선택 함수
func (lb *LoadBalancer) selectQueue(request dto.Request) *NodeQueue {
	var selectedQueue *NodeQueue
	if len(lb.Queues) == 0 {
		return nil
	}

	for _, queue := range lb.Queues {
		if queue.BPM >= request.BPM && queue.RPM >= request.RPM {
			selectedQueue = queue
			break
		}
	}

	// 적합한 큐가 없는 경우 대기 큐에 추가
	if selectedQueue == nil {
		lb.waitQueue <- &request
		return nil
	}

	return selectedQueue
}

// 요청을 처리하는 ServeHTTP 함수 수정
func (lb *LoadBalancer) ServeHTTP(request dto.Request) {
	queue := lb.selectQueue(request)

	if queue != nil {
		queue.SendRequestToNode(&request)
	}
}

// 대기 큐의 요청을 처리하는 워커 함수
func (lb *LoadBalancer) processQueue() {
	for req := range lb.waitQueue {
		processed := false
		selectedQueue := lb.selectQueue(*req)

		if selectedQueue != nil {
			selectedQueue.SendRequestToNode(req)
			processed = true
		}

		if !processed {
			// TODO: 요청을 다시 대기 큐에 추가하거나 다른 처리 로직
		}
	}
}
