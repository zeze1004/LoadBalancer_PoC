package internal

import "github.com/zeze1004/LoadBanlence_PoC/internal/model"

// RequestEventQueue 요청을 담는 큐
type RequestEventQueue struct {
	ReqQueue chan model.RequestEvent
}

func NewEventQueue(size int) *RequestEventQueue {
	return &RequestEventQueue{
		ReqQueue: make(chan model.RequestEvent, size),
	}
}

func (eq *RequestEventQueue) Add(event model.RequestEvent) bool {
	select {
	case eq.ReqQueue <- event:
		return true
	default:
		return false // 큐가 가득 찬 경우
	}
}

func (eq *RequestEventQueue) Pop() *model.RequestEvent {
	select {
	case event := <-eq.ReqQueue:
		return &event
	default:
		return nil // 큐가 비어 있는 경우
	}
}
