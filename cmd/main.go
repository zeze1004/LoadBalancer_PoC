package main

import (
	"fmt"
	"github.com/zeze1004/LoadBanlence_PoC/internal"
	"github.com/zeze1004/LoadBanlence_PoC/internal/model"
	"net/http"
	"time"
)

func main() {
	// 더미용 노드 생성
	nodes := []*model.AInode{
		{ID: "node1", URL: "http://localhost:9001", LimitBPM: 1000000, LimitRPM: 100, IsActive: true},
		{ID: "node2", URL: "http://localhost:9002", LimitBPM: 2000000, LimitRPM: 150, IsActive: true},
	}
	rateLimiter := internal.NewRateLimiter(nodes)

	// 이벤트 대기열 생성
	eventQueue := internal.NewEventQueue(100)

	// 노드 워커 시작
	for _, node := range nodes {
		nodeWorker := internal.NewAInodeWorker(node, rateLimiter)
		go nodeWorker.ProcessingReqEvent(eventQueue)
	}

	// 노드 상태 모니터링 시작
	nodeMonitor := internal.NewNodeMonitor(nodes)
	go nodeMonitor.StartMonitoring(10 * time.Second)

	// HTTP 요청 처리
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestSize := 1000
		event := model.RequestEvent{
			RequestID:   fmt.Sprintf("req-%d", time.Now().UnixNano()),
			RequestSize: requestSize,
			Timestamp:   time.Now(),
			ResponseCh:  make(chan *model.ResponseEvent),
		}

		// 요청이 처리 됐을 때 응답을 보내는 로직은 구현이 되어 있지 않습니다
		// 현재 구현된 응답은 요청에 대한 대기열 처리 여부에 대해만 반환됩니다
		select {
		case response := <-event.ResponseCh:
			if response.IsSuccess {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("요청이 성공적으로 처리되었습니다"))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("요청 처리에 실패했습니다"))
			}
		default:
			if !eventQueue.Add(event) {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte("요청 대기열이 가득 찼습니다"))
			} else {
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write([]byte("요청이 대기열에 추가되었습니다"))
			}
		}
	})

	fmt.Println("Server is running on port 8080...")
	_ = http.ListenAndServe(":8080", nil)
}
