package pybridge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// Client is a lightweight HTTP client to the Python validator service.
type Client struct {
	baseURL    string
	httpClient *http.Client
	BaseURL    string
	Timeout    time.Duration
}

type EvalValidationInput struct {
	ScriptID          string `json:"script_id"`
	EvaluatorID       string `json:"evaluator_id"`
	TotalQuestions    int    `json:"total_questions"`
	MarksPerQuestion  int    `json:"marks_per_question"`
	TotalMarks        int    `json:"total_marks"`
	QuestionsAnswered int    `json:"questions_answered"`
	MarksAllotted     []int  `json:"marks_allotted_per_question"`
	MarksScored       []int  `json:"marks_scored"`
	CourseID          string `json:"course_id"`
	Semester          string `json:"semester"`
	CourseCredits     int    `json:"course_credits"`
}

// NewClient creates a new client with sensible timeouts.
func NewClient(baseURL string, timeout time.Duration) *Client {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8082"
	}
	// custom transport with dial timeout to avoid long hangs
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          10,
		IdleConnTimeout:       30 * time.Second,
	}

	trimmed := strings.TrimRight(baseURL, "/")
	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
	return &Client{
		baseURL:    trimmed,
		BaseURL:    trimmed, // exported mirror for visibility if needed elsewhere
		httpClient: httpClient,
		Timeout:    timeout, // exported mirror
	}
}

// Extract calls Python /extract endpoint with a local file path.
func (c *Client) Extract(ctx context.Context, filePath string) (*ExtractResponse, error) {
	reqBody := ExtractRequest{FilePath: filePath}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	url := c.baseURL + "/extract"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("pybridge extract request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Check status code first so we always surface the error detail.
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("python extractor error (status %d): %s", resp.StatusCode, string(body))
	}

	var out ExtractResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode extract response: %w", err)
	}
	return &out, nil
}

// Validate calls Python /validate endpoint with metadata map.
func (c *Client) Validate(ctx context.Context, metadata map[string]string) (*ValidateResponse, error) {
	if c.baseURL == "" {
		c.baseURL = "http://127.0.0.1:8082"
	}
	reqBody := ValidateRequest{Metadata: metadata}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	url := c.baseURL + "/validate"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("pybridge validate request failed: %w", err)
	}
	defer resp.Body.Close()

	var out ValidateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode validate response: %w", err)
	}
	if resp.StatusCode >= 400 {
		if out.Error == "" {
			return &out, fmt.Errorf("python service error status %d", resp.StatusCode)
		}
		return &out, fmt.Errorf("python service error: %s", out.Error)
	}
	return &out, nil
}

func (c *Client) ValidateEvaluation(ctx context.Context, in EvalValidationInput) (*ValidateResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("pybridge client nil")
	}
	// prefer the internal baseURL (always trimmed and non-empty)
	url := c.baseURL + "/validate-evaluation"

	j, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// reuse configured http client (with transport & timeout)
	client := c.httpClient
	if client == nil {
		// fallback
		client = &http.Client{Timeout: c.Timeout}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pybridge validate-evaluation request failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("validator error status %d: %s", resp.StatusCode, string(body))
	}

	var out ValidateResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode validate-evaluation response: %w", err)
	}
	return &out, nil
}
