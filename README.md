# LoadBanlence_PoC
로드밸런스 기능 구현 PoC 저장소

## 프로젝트 설명
Go 언어로 AI 프록시 기능을 PoC 하는 로드밸런서 프로젝트를 구현합니다.
- 더 자세한 설명은 [요구사항 정리 문서](requirement.md)를 참고해주세요.

## 프로젝트 설계
이벤트 드리븐 기반 구조를 사용했습니다.
- 요청이 들어오면, 이벤트 대기열에 추가됩니다.
- `Node Worker`가 이벤트를 소비하여 요청을 처리합니다.
- 속도 제한은 이벤트가 대기열에서 꺼내지기 전에 확인됩니다.

## 요청이 들어왔을 때 로드밸런서의 처리 과정

1. 클라이언트의 `POST` 요청 수신
    
   - 서버는 요청의 크기를 측정합니다. 
   - PoC 코드에서는 임의의 requestSize를 지정합니다.
       ```
       http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
       requestSize := 100 // 임의의 requestSize
       event := model.RequestEvent{
           RequestID:   fmt.Sprintf("req-%d", time.Now().UnixNano()),
           RequestSize: requestSize,
           Timestamp:   time.Now(),
       }
       ```

2. 요청 대기열(`Event Queue`)에 저장
   - 서버는 요청을 대기열에 저장합니다.
      ```
       if !eventQueue.Add(event) {
              ...
         }
       ```
     
3. `Node Worker`가 대기열에서 요청을 처리
   - `Node Worker`는 대기열에서 요청을 꺼냅니다.
   - 해당 요청을 처리할 수 있는 노드를 선택하고, 노드가 요청을 처리할 수 있는지 속도제한 여부를 확인하고자 `RateLimiter`를 호출합니다.
   - 요청을 처리할 수 있는 경우 노드로 요청을 전달합니다.
     ```
     func (worker *NodeWorker) ProcessingReqEvent(queue *EventQueue) {
           for event := range queue.Queue {
               if !worker.RateLimiter.AllowRequest(worker.Node.ID, event.RequestSize) {
                   continue // 제한 초과 시 요청 처리 건너뜀
               }
               worker.SendRequest(event) // 노드로 요청 전달
           }
       }
     ```
 
4. `RateLimiter`의 속도 제한 확인
   - `RateLimiter`는 노드의 속도 제한 여부를 확인합니다.
     ```
     func (limiter *RateLimiter) AllowRequest(nodeID string, requestSize int) bool {
         limiter.mu.Lock()
         defer limiter.mu.Unlock()
  
         node, ok := limiter.Nodes[nodeID]
         if !ok {
             return false
         }
  
         if node.BPM < requestSize || node.RPM < 1 {
             return false
         }
  
         node.BPM -= requestSize
         node.RPM--
         return true
     }
     ```

5. 노드로 요청 전달
   - Node Worker는 적절한 노드로 요청을 전달합니다. 
   - 노드의 URL로 요청을 전송합니다.
     ```
       func (worker *NodeWorker) SendRequest(event model.RequestEvent) error {
          client := &http.Client{Timeout: 10 * time.Second}
          req, err := http.NewRequest("POST", worker.Node.URL, nil) // 요청 본문 추가 가능
          if err != nil {
              return err
          }
          
     ```
   
6. 주기적으로 노드 모니터링
   - `NodeMonitor`는 주기적으로 노드의 상태를 확인합니다.
       ```
       func (nm *NodeMonitor) checkNodeStatus(node *model.Node) {
           resp, err := http.Get(node.URL + "/health")
           if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
               node.IsActive = false
           } else {
               node.IsActive = true
           }
           node.LastChecked = time.Now()
       }
       ```

7. 요청 완료 처리

    - 요청이 완료되면 클라이언트에게 응답을 반환합니다.
    - 현재 응답은 요청이 큐에 적재 되었는지에 대한 여부만을 반환합니다.

## 확장 가능한 기능
1. 요청을 처리할 노드가 없을 때, 해당 요청을 가장 우선순위 높게 처리할 수 있는 로직이 추가 될 수 있습니다.
2. Cloud Service 별 API Key를 로드밸런서에서 중앙으로 관리하고 노드에 전달할 수 있습니다.
3. 노드의 상태를 주기적으로 확인할 때, 노드의 상태가 비정상일 때, 해당 노드를 제외하고 요청을 처리할 수 있는 로직이 추가 될 수 있습니다.
