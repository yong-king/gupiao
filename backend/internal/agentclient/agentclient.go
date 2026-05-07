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
	body, err := json.Marshal(request)
	if err != nil {
		return AnalyzeResult{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/analyze", bytes.NewReader(body))
	if err != nil {
		return AnalyzeResult{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return AnalyzeResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return AnalyzeResult{}, errors.New("agent returned non-200 status")
	}

	var result AnalyzeResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return AnalyzeResult{}, err
	}
	if result.Signal == "" {
		return AnalyzeResult{}, errors.New("agent response missing signal")
	}
	return result, nil
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
