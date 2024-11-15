package internal

import (
	"fmt"
	"github.com/zeze1004/LoadBanlence_PoC/internal/model"
	"io"
	"net/http"
	"time"
)

type AInodeWorker struct {
	AInode      *model.AInode
	RateLimiter *RateLimiter
}

func NewAInodeWorker(node *model.AInode, rateLimiter *RateLimiter) *AInodeWorker {
	return &AInodeWorker{
		AInode:      node,
		RateLimiter: rateLimiter,
	}
}

// ProcessingReqEvent AInodeWorker에 있는 노드가 요청 이벤트를 처리하도록
func (nw *AInodeWorker) ProcessingReqEvent(eventQueue *RequestEventQueue) {
	for event := range eventQueue.ReqQueue {
		if !nw.RateLimiter.AllowRequest(nw.AInode.ID, event.RequestSize) {
			fmt.Printf("Node %s의 속도 제한이 초과되어 %s 요청을 처리할 수 없습니다\n", nw.AInode.ID, event.RequestID)
			continue
		}

		err := nw.SendRequest(event)
		if err != nil {
			fmt.Printf("Node %s가 %s요청 처리를 실패했습니\n err: %v\n", nw.AInode.ID, event.RequestID, err)
		}
	}
}

func (nw *AInodeWorker) SendRequest(event model.RequestEvent) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", nw.AInode.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", event.RequestSize))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("응답 body를 닫는데 실패했습니다: %v\n", err)
		}
	}(resp.Body)

	fmt.Printf("Node %s가 성공적으로 %s 요청을 처리했습니다\n", nw.AInode.ID, event.RequestID)
	return nil
}
