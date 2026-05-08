package agentclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

type AnalyzeRequest struct {
	Symbol         string                   `json:"symbol"`
	Market         string                   `json:"market"`
	Prices         []float64                `json:"prices"`
	TriggeredRules []string                 `json:"triggered_rules"`
	Texts          []string                 `json:"texts"`
	SourceRefs     []map[string]interface{} `json:"source_refs"`
}

type AnalyzeResult struct {
	Signal            string                   `json:"signal"`
	Confidence        float64                  `json:"confidence"`
	RiskLevel         string                   `json:"risk_level"`
	TriggeredRules    []string                 `json:"triggered_rules"`
	Summary           string                   `json:"summary"`
	Reasoning         []string                 `json:"reasoning"`
	DataTime          string                   `json:"data_time"`
	SourceRefs        []map[string]interface{} `json:"source_refs"`
	MissingData       []string                 `json:"missing_data"`
	RecommendedAction string                   `json:"recommended_action"`
	Indicators        map[string]interface{}   `json:"indicators"`
}

type WorkflowRequest struct {
	UserID         string                 `json:"user_id"`
	JobID          string                 `json:"job_id"`
	Market         string                 `json:"market"`
	Symbol         string                 `json:"symbol"`
	AttentionLevel string                 `json:"attention_level"`
	Interval       string                 `json:"interval"`
	Profile        map[string]interface{} `json:"profile"`
	LatestSnapshot map[string]interface{} `json:"latest_snapshot"`
	SnapshotsCount int                    `json:"snapshots_count"`
}

type WorkflowStep struct {
	StepName      string `json:"step_name"`
	AgentName     string `json:"agent_name"`
	Status        string `json:"status"`
	InputSummary  string `json:"input_summary"`
	OutputSummary string `json:"output_summary"`
	Model         string `json:"model"`
	StartedAt     string `json:"started_at"`
	CompletedAt   string `json:"completed_at"`
}

type WorkflowResult struct {
	Engine   string            `json:"engine"`
	Market   string            `json:"market"`
	Symbol   string            `json:"symbol"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
	Steps    []WorkflowStep    `json:"steps"`
}

type ChatRequest struct {
	UserID         string                   `json:"user_id"`
	SessionID      string                   `json:"session_id"`
	Market         string                   `json:"market"`
	Symbol         string                   `json:"symbol"`
	Question       string                   `json:"question"`
	ContextSummary string                   `json:"context_summary"`
	History        []map[string]interface{} `json:"history"`
	RAGDocuments   []map[string]interface{} `json:"rag_documents"`
	Profile        map[string]interface{}   `json:"profile"`
}

type ChatResult struct {
	Market    string `json:"market"`
	Symbol    string `json:"symbol"`
	Answer    string `json:"answer"`
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
}

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) Analyze(ctx context.Context, request AnalyzeRequest) (AnalyzeResult, error) {
	var result AnalyzeResult
	if err := c.post(ctx, "/analyze", request, &result); err != nil {
		return AnalyzeResult{}, err
	}
	if result.Signal == "" {
		return AnalyzeResult{}, errors.New("agent response missing signal")
	}
	return result, nil
}

func (c *Client) RunWorkflow(ctx context.Context, request WorkflowRequest) (WorkflowResult, error) {
	var result WorkflowResult
	if err := c.post(ctx, "/workflow/research", request, &result); err != nil {
		return WorkflowResult{}, err
	}
	if result.Content == "" || len(result.Steps) == 0 {
		return WorkflowResult{}, errors.New("agent workflow response missing content or steps")
	}
	return result, nil
}

func (c *Client) Chat(ctx context.Context, request ChatRequest) (ChatResult, error) {
	var result ChatResult
	if err := c.post(ctx, "/assistant/chat", request, &result); err != nil {
		return ChatResult{}, err
	}
	if result.Answer == "" {
		return ChatResult{}, errors.New("agent chat response missing answer")
	}
	return result, nil
}

func (c *Client) post(ctx context.Context, path string, request interface{}, result interface{}) error {
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("agent returned non-200 status")
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}

type Run struct {
	ID            string
	AgentVersion  string
	PromptVersion string
	Status        string
	Input         AnalyzeRequest
	Output        AnalyzeResult
	Error         string
	CreatedAt     time.Time
}

type RunRepository struct {
	mu   sync.RWMutex
	runs map[string]Run
	now  func() time.Time
}

func NewRunRepository() *RunRepository {
	return &RunRepository{runs: make(map[string]Run), now: func() time.Time { return time.Now().UTC() }}
}

func (r *RunRepository) Save(run Run) error {
	if run.ID == "" {
		return errors.New("run id is required")
	}
	if run.CreatedAt.IsZero() {
		run.CreatedAt = r.now()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runs[run.ID] = run
	return nil
}

func (r *RunRepository) FindByID(id string) (Run, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	run, ok := r.runs[id]
	return run, ok
}
