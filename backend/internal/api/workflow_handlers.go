package api

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"jijin/backend/internal/agentclient"
	"jijin/backend/internal/holdings"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/persistence"
)

type runWorkflowRequest struct {
	UserID         string `json:"user_id"`
	AttentionLevel string `json:"attention_level"`
	Market         string `json:"market"`
	Symbol         string `json:"symbol"`
}

type workflowTarget struct {
	Market         string `json:"market"`
	Symbol         string `json:"symbol"`
	AttentionLevel string `json:"attention_level"`
	Source         string `json:"source"`
}

type runWorkflowResponse struct {
	Job       persistence.WorkflowJob   `json:"job"`
	Targets   []workflowTarget          `json:"targets"`
	Documents []persistence.RAGDocument `json:"documents"`
}

func (s *Server) handleWorkflowRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	var req runWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	userID := strings.TrimSpace(req.UserID)
	if userID == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "user_id is required", requestID(r))
		return
	}
	attention := holdings.NormalizeAttentionLevel(req.AttentionLevel)
	targets := s.workflowTargets(userID, attention, req.Market, req.Symbol)
	now := time.Now().UTC()
	job := persistence.WorkflowJob{
		ID:             "wf-" + shortHash(userID+attention+req.Market+req.Symbol+now.Format(time.RFC3339Nano)),
		UserID:         userID,
		WorkflowType:   "attention_research",
		AttentionLevel: attention,
		Market:         strings.ToUpper(strings.TrimSpace(req.Market)),
		Symbol:         strings.ToUpper(strings.TrimSpace(req.Symbol)),
		Status:         "running",
		TargetCount:    len(targets),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if s.store != nil {
		if err := s.store.SaveWorkflowJob(job); err != nil {
			WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
			return
		}
	}

	documents := []persistence.RAGDocument{}
	steps := []persistence.WorkflowStep{}
	for _, target := range targets {
		document, targetSteps := s.runResearchTarget(r, job.ID, userID, target)
		documents = append(documents, document)
		steps = append(steps, targetSteps...)
	}
	job.Status = "succeeded"
	job.Summary = fmt.Sprintf("已按%s关注等级处理 %d 只股票，写入 %d 份 RAG/向量记录。", attentionLabel(attention), len(targets), len(documents))
	job.UpdatedAt = time.Now().UTC()
	job.Steps = steps
	if s.store != nil {
		if err := s.store.SaveWorkflowJob(job); err != nil {
			WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
			return
		}
	}
	writeJSON(w, http.StatusOK, runWorkflowResponse{Job: job, Targets: targets, Documents: documents})
}

func (s *Server) handleWorkflows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if userID == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "user_id is required", requestID(r))
		return
	}
	if s.store == nil {
		writeJSON(w, http.StatusOK, []persistence.WorkflowJob{})
		return
	}
	jobs, err := s.store.ListWorkflowJobs(userID, 20)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
		return
	}
	writeJSON(w, http.StatusOK, jobs)
}

func (s *Server) workflowTargets(userID string, attention string, market string, symbol string) []workflowTarget {
	normalizedMarket := strings.ToUpper(strings.TrimSpace(market))
	normalizedSymbol := strings.ToUpper(strings.TrimSpace(symbol))
	if normalizedMarket != "" && normalizedSymbol != "" {
		return []workflowTarget{{Market: normalizedMarket, Symbol: normalizedSymbol, AttentionLevel: attention, Source: "explicit"}}
	}
	out := []workflowTarget{}
	for _, item := range s.listHoldings(userID) {
		itemAttention := holdings.NormalizeAttentionLevel(item.AttentionLevel)
		if itemAttention != attention {
			continue
		}
		out = append(out, workflowTarget{Market: item.Market, Symbol: item.Symbol, AttentionLevel: itemAttention, Source: "holding"})
	}
	return out
}

func (s *Server) runResearchTarget(r *http.Request, jobID string, userID string, target workflowTarget) (persistence.RAGDocument, []persistence.WorkflowStep) {
	started := time.Now().UTC()
	snapshots := s.listSnapshots(target.Market, target.Symbol)
	if len(snapshots) == 0 {
		if fetched, err := s.refreshes.Provider.FetchQuotes(r.Context(), []marketdata.QuoteRequest{{Market: target.Market, Symbol: target.Symbol}}); err == nil {
			_ = s.refreshes.Snapshots.SaveAll(fetched)
			if s.store != nil {
				_ = s.store.SaveSnapshots(fetched)
			}
			snapshots = s.listSnapshots(target.Market, target.Symbol)
		}
	}
	profile := marketdata.ProfileFromSnapshots(target.Market, target.Symbol, snapshots)
	latest := latestSnapshotForWorkflow(snapshots)
	interval := holdings.AttentionRefreshInterval(target.AttentionLevel)
	if document, steps, ok := s.runAgentResearchTarget(r, jobID, userID, target, profile, latest, len(snapshots), interval); ok {
		return document, steps
	}
	collectorOutput := fmt.Sprintf("%s:%s 行情样本 %d 条，产品 %s。", target.Market, target.Symbol, len(snapshots), strings.Join(profile.Products, "、"))
	summary := workflowSummary(profile, latest, target.AttentionLevel, interval)
	review := "风险审查：输出仅用于研究和提醒，不作为买卖指令；需结合公告、财报和个人仓位人工确认。"
	content := summary + "\n" + review
	document := persistence.RAGDocument{
		ID:         "rag-" + shortHash(userID+target.Market+target.Symbol+time.Now().UTC().Format(time.RFC3339Nano)),
		UserID:     userID,
		Market:     target.Market,
		Symbol:     target.Symbol,
		SourceType: "multi_agent_workflow",
		SourceID:   jobID,
		Content:    content,
		Metadata: map[string]string{
			"attention_level":  target.AttentionLevel,
			"refresh_interval": interval.String(),
			"workflow_job_id":  jobID,
			"agents":           "collector,summarizer,risk_reviewer,rag_writer",
			"embedding_status": "indexed",
		},
		Embedding: localEmbedding(content),
		CreatedAt: time.Now().UTC(),
	}
	if s.store != nil {
		_ = s.store.SaveRAGDocument(document)
	}
	steps := []persistence.WorkflowStep{
		workflowStep(jobID, target, "market_context_collect", "行情上下文采集 Agent", "获取行情、持仓等级和公司产品上下文", collectorOutput, started),
		workflowStep(jobID, target, "company_product_collect", "股票产品信息采集 Agent", "读取公开产品/业务资料和已保存 profile", profile.Business, started.Add(time.Millisecond)),
		workflowStep(jobID, target, "summarize", "归纳整理 Agent", collectorOutput, summary, started.Add(2*time.Millisecond)),
		workflowStep(jobID, target, "risk_review", "风险审查 Agent", summary, review, started.Add(3*time.Millisecond)),
		workflowStep(jobID, target, "rag_vector_write", "RAG/向量写入 Agent", "写入 rag_documents 和 rag_vectors", document.ID, started.Add(4*time.Millisecond)),
	}
	if s.store != nil {
		for _, step := range steps {
			_ = s.store.SaveWorkflowStep(step)
		}
	}
	return document, steps
}

func (s *Server) runAgentResearchTarget(r *http.Request, jobID string, userID string, target workflowTarget, profile marketdata.CompanyProfile, latest marketdata.Snapshot, snapshotsCount int, interval time.Duration) (persistence.RAGDocument, []persistence.WorkflowStep, bool) {
	if strings.TrimSpace(s.cfg.AgentURL) == "" {
		return persistence.RAGDocument{}, nil, false
	}
	ctx := r.Context()
	client := agentclient.NewClient(strings.TrimRight(s.cfg.AgentURL, "/"))
	result, err := client.RunWorkflow(ctx, agentclient.WorkflowRequest{
		UserID:         userID,
		JobID:          jobID,
		Market:         target.Market,
		Symbol:         target.Symbol,
		AttentionLevel: target.AttentionLevel,
		Interval:       interval.String(),
		Profile:        mapFromJSON(profile),
		LatestSnapshot: mapFromJSON(latest),
		SnapshotsCount: snapshotsCount,
	})
	if err != nil {
		return persistence.RAGDocument{}, nil, false
	}
	metadata := map[string]string{}
	for key, value := range result.Metadata {
		metadata[key] = value
	}
	metadata["agent_engine"] = result.Engine
	metadata["workflow_job_id"] = jobID
	if _, ok := metadata["embedding_status"]; !ok {
		metadata["embedding_status"] = "indexed"
	}
	content := strings.TrimSpace(result.Content)
	if content == "" {
		content = workflowSummary(profile, latest, target.AttentionLevel, interval)
	}
	document := persistence.RAGDocument{
		ID:         "rag-" + shortHash(userID+target.Market+target.Symbol+time.Now().UTC().Format(time.RFC3339Nano)),
		UserID:     userID,
		Market:     target.Market,
		Symbol:     target.Symbol,
		SourceType: "langgraph_agent_workflow",
		SourceID:   jobID,
		Content:    content,
		Metadata:   metadata,
		Embedding:  localEmbedding(content),
		CreatedAt:  time.Now().UTC(),
	}
	if s.store != nil {
		_ = s.store.SaveRAGDocument(document)
	}
	steps := make([]persistence.WorkflowStep, 0, len(result.Steps))
	for _, step := range result.Steps {
		started := parseAgentTime(step.StartedAt)
		completed := parseAgentTime(step.CompletedAt)
		output := step.OutputSummary
		if strings.TrimSpace(step.Model) != "" {
			output = output + " [模型: " + step.Model + "; 引擎: " + result.Engine + "]"
		}
		steps = append(steps, persistence.WorkflowStep{
			ID:            jobID + ":" + target.Market + ":" + target.Symbol + ":" + step.StepName,
			JobID:         jobID,
			StepName:      step.StepName,
			AgentName:     step.AgentName,
			Status:        step.Status,
			InputSummary:  step.InputSummary,
			OutputSummary: output,
			StartedAt:     started,
			CompletedAt:   completed,
		})
	}
	if s.store != nil {
		for _, step := range steps {
			_ = s.store.SaveWorkflowStep(step)
		}
	}
	return document, steps, true
}

func workflowStep(jobID string, target workflowTarget, name string, agentName string, input string, output string, at time.Time) persistence.WorkflowStep {
	return persistence.WorkflowStep{
		ID:            jobID + ":" + target.Market + ":" + target.Symbol + ":" + name,
		JobID:         jobID,
		StepName:      name,
		AgentName:     agentName,
		Status:        "succeeded",
		InputSummary:  input,
		OutputSummary: output,
		StartedAt:     at,
		CompletedAt:   at.Add(50 * time.Millisecond),
	}
}

func workflowSummary(profile marketdata.CompanyProfile, snapshot marketdata.Snapshot, attention string, interval time.Duration) string {
	change := "暂无行情样本"
	if snapshot.Symbol != "" {
		change = fmt.Sprintf("最新价 %.2f，涨跌幅 %.2f%%，来源 %s", snapshot.Price, snapshot.ChangePercent, snapshot.Source)
	}
	return fmt.Sprintf("%s:%s 属于%s，重点产品/业务为%s。关注等级%s，建议采集周期%s。行情：%s。归纳：%s",
		profile.Market, profile.Symbol, profile.Sector, strings.Join(profile.Products, "、"), attentionLabel(attention), interval, change, profile.Analysis)
}

func latestSnapshotForWorkflow(snapshots []marketdata.Snapshot) marketdata.Snapshot {
	var latest marketdata.Snapshot
	for _, snapshot := range snapshots {
		if latest.Symbol == "" || snapshot.DataTime.After(latest.DataTime) {
			latest = snapshot
		}
	}
	return latest
}

func mapFromJSON(value interface{}) map[string]interface{} {
	data, err := json.Marshal(value)
	if err != nil {
		return map[string]interface{}{}
	}
	out := map[string]interface{}{}
	_ = json.Unmarshal(data, &out)
	return out
}

func parseAgentTime(value string) time.Time {
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed
	}
	return time.Now().UTC()
}

func localEmbedding(text string) []float64 {
	vector := make([]float64, 16)
	words := strings.Fields(strings.ToLower(text))
	if len(words) == 0 {
		return vector
	}
	for _, word := range words {
		hash := sha1.Sum([]byte(word))
		index := int(binary.BigEndian.Uint32(hash[:4]) % uint32(len(vector)))
		vector[index] += 1
	}
	norm := 0.0
	for _, value := range vector {
		norm += value * value
	}
	norm = math.Sqrt(norm)
	if norm == 0 {
		return vector
	}
	for i := range vector {
		vector[i] = math.Round((vector[i]/norm)*10000) / 10000
	}
	return vector
}

func shortHash(value string) string {
	hash := sha1.Sum([]byte(value))
	return hex.EncodeToString(hash[:])[:20]
}

func attentionLabel(level string) string {
	switch holdings.NormalizeAttentionLevel(level) {
	case "high":
		return "高"
	case "low":
		return "低"
	default:
		return "中"
	}
}

type assistantChatRequest struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Market    string `json:"market"`
	Symbol    string `json:"symbol"`
	Question  string `json:"question"`
}

type assistantChatResponse struct {
	Market         string   `json:"market"`
	Symbol         string   `json:"symbol"`
	SessionID      string   `json:"session_id"`
	Answer         string   `json:"answer"`
	ContextSummary string   `json:"context_summary"`
	RAGDocumentIDs []string `json:"rag_document_ids"`
	Model          string   `json:"model"`
}

func (s *Server) handleAssistantChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	var req assistantChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	userID := strings.TrimSpace(req.UserID)
	if userID == "" || strings.TrimSpace(req.Symbol) == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "user_id and symbol are required", requestID(r))
		return
	}
	response, err := s.buildAssistantChatResponse(r, req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "assistant_error", err.Error(), requestID(r))
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) buildAssistantChatResponse(r *http.Request, req assistantChatRequest) (assistantChatResponse, error) {
	userID := strings.TrimSpace(req.UserID)
	sessionID := strings.TrimSpace(req.SessionID)
	if sessionID == "" {
		sessionID = "default"
	}
	market, symbol := s.resolveAssistantTarget(userID, req.Market, req.Symbol)
	snapshots := s.listSnapshots(market, symbol)
	profile := marketdata.ProfileFromSnapshots(market, symbol, snapshots)
	holding := s.findHolding(userID, market, symbol)
	inPool := s.stockInAnyWatchlist(userID, market, symbol)
	ragDocs := []persistence.RAGDocument{}
	if s.store != nil {
		if docs, err := s.store.ListRAGDocuments(userID, market, symbol, 5); err == nil {
			ragDocs = docs
		}
	}
	ragIDs := make([]string, 0, len(ragDocs))
	ragText := "暂无 RAG 历史总结"
	if len(ragDocs) > 0 {
		ragText = ragDocs[0].Content
		for _, doc := range ragDocs {
			ragIDs = append(ragIDs, doc.ID)
		}
	}
	contextSummary := assistantContextSummary(holding, inPool, len(ragDocs), len(snapshots))
	history := []persistence.AssistantMessage{}
	if s.store != nil {
		if messages, err := s.store.ListAssistantMessages(userID, sessionID, market, symbol, 8); err == nil {
			history = messages
		}
	}
	answer := fmt.Sprintf("针对 %s:%s：%s。你问的是：%s。结合当前上下文：%s。RAG 参考：%s。结论：继续作为研究对象观察，重点核对产品/业务变化、价格波动和你的持仓成本；这不是买卖指令。",
		market, symbol, profile.Analysis, strings.TrimSpace(req.Question), contextSummary, ragText)
	model := "local-fallback"
	if strings.TrimSpace(s.cfg.AgentURL) != "" {
		client := agentclient.NewClient(strings.TrimRight(s.cfg.AgentURL, "/"))
		if result, err := client.Chat(r.Context(), agentclient.ChatRequest{
			UserID:         userID,
			SessionID:      sessionID,
			Market:         market,
			Symbol:         symbol,
			Question:       req.Question,
			ContextSummary: contextSummary,
			History:        assistantHistoryMaps(history),
			RAGDocuments:   ragDocumentMaps(ragDocs),
			Profile:        mapFromJSON(profile),
		}); err == nil {
			answer = result.Answer
			model = result.Model
		}
	}
	if s.store != nil {
		_ = s.store.SaveAssistantMessage(persistence.AssistantMessage{
			ID:             "chat-" + shortHash(userID+market+symbol+req.Question+time.Now().UTC().Format(time.RFC3339Nano)),
			UserID:         userID,
			SessionID:      sessionID,
			Market:         market,
			Symbol:         symbol,
			Question:       req.Question,
			Answer:         answer,
			ContextSummary: contextSummary,
			RAGDocumentIDs: ragIDs,
			CreatedAt:      time.Now().UTC(),
		})
	}
	return assistantChatResponse{Market: market, Symbol: symbol, SessionID: sessionID, Answer: answer, ContextSummary: contextSummary, RAGDocumentIDs: ragIDs, Model: model}, nil
}

func (s *Server) handleAssistantChatStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	var req assistantChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.Symbol) == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "user_id and symbol are required", requestID(r))
		return
	}
	response, err := s.buildAssistantChatResponse(r, req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "assistant_error", err.Error(), requestID(r))
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, _ := w.(http.Flusher)
	for _, chunk := range streamChunks(response.Answer, 18) {
		payload, _ := json.Marshal(map[string]string{"delta": chunk})
		_, _ = fmt.Fprintf(w, "data: %s\n\n", payload)
		if flusher != nil {
			flusher.Flush()
		}
	}
	final, _ := json.Marshal(map[string]interface{}{"done": true, "response": response})
	_, _ = fmt.Fprintf(w, "data: %s\n\n", final)
	if flusher != nil {
		flusher.Flush()
	}
}

func (s *Server) resolveAssistantTarget(userID string, market string, symbol string) (string, string) {
	normalizedSymbol := strings.ToUpper(strings.TrimSpace(symbol))
	normalizedMarket := strings.ToUpper(strings.TrimSpace(market))
	if normalizedMarket != "" {
		return normalizedMarket, normalizedSymbol
	}
	for _, holding := range s.listHoldings(userID) {
		if holding.Symbol == normalizedSymbol {
			return holding.Market, normalizedSymbol
		}
	}
	for _, watchlist := range s.listWatchlists(userID) {
		for _, item := range watchlist.Symbols {
			if item.Symbol == normalizedSymbol {
				return item.Market, normalizedSymbol
			}
		}
	}
	return "CN", normalizedSymbol
}

func (s *Server) findHolding(userID string, market string, symbol string) *holdings.Holding {
	for _, item := range s.listHoldings(userID) {
		if item.Market == market && item.Symbol == symbol {
			return &item
		}
	}
	return nil
}

func (s *Server) stockInAnyWatchlist(userID string, market string, symbol string) bool {
	for _, watchlist := range s.listWatchlists(userID) {
		for _, item := range watchlist.Symbols {
			if item.Market == market && item.Symbol == symbol {
				return true
			}
		}
	}
	return false
}

func (s *Server) listWatchlists(userID string) []struct {
	ID      string
	Name    string
	Symbols []struct {
		Market string
		Symbol string
	}
} {
	out := []struct {
		ID      string
		Name    string
		Symbols []struct {
			Market string
			Symbol string
		}
	}{}
	var source []struct {
		ID      string
		Name    string
		Symbols []struct{ Market, Symbol string }
	}
	_ = source
	items := s.watchlists.ListByUser(userID)
	if s.store != nil {
		if stored, err := s.store.ListWatchlistsByUser(userID); err == nil {
			items = stored
		}
	}
	for _, item := range items {
		row := struct {
			ID      string
			Name    string
			Symbols []struct {
				Market string
				Symbol string
			}
		}{ID: item.ID, Name: item.Name}
		for _, symbol := range item.Symbols {
			row.Symbols = append(row.Symbols, struct {
				Market string
				Symbol string
			}{Market: symbol.Market, Symbol: symbol.Symbol})
		}
		out = append(out, row)
	}
	return out
}

func assistantContextSummary(holding *holdings.Holding, inPool bool, ragCount int, snapshotCount int) string {
	parts := []string{}
	if holding != nil {
		parts = append(parts, fmt.Sprintf("持仓数量 %.4f，成本 %.2f，关注等级%s", holding.Quantity, holding.CostBasis, attentionLabel(holding.AttentionLevel)))
	}
	if inPool {
		parts = append(parts, "该股票在股票池中")
	}
	parts = append(parts, "RAG 文档 "+strconv.Itoa(ragCount)+" 份")
	parts = append(parts, "行情样本 "+strconv.Itoa(snapshotCount)+" 条")
	return strings.Join(parts, "；")
}

func assistantHistoryMaps(messages []persistence.AssistantMessage) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(messages))
	for _, message := range messages {
		out = append(out, map[string]interface{}{
			"question": message.Question,
			"answer":   message.Answer,
			"created":  message.CreatedAt.Format(time.RFC3339),
		})
	}
	return out
}

func ragDocumentMaps(documents []persistence.RAGDocument) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(documents))
	for _, document := range documents {
		out = append(out, map[string]interface{}{
			"id":      document.ID,
			"content": document.Content,
			"source":  document.SourceType,
		})
	}
	return out
}

func streamChunks(text string, size int) []string {
	if size <= 0 {
		size = 24
	}
	runes := []rune(text)
	chunks := []string{}
	for start := 0; start < len(runes); start += size {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	if len(chunks) == 0 {
		return []string{""}
	}
	return chunks
}
