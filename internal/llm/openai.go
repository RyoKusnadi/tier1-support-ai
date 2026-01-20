package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIClient implements the Client interface for OpenAI
type OpenAIClient struct {
	config      Config
	httpClient  *http.Client
	promptBuilder *PromptBuilder
	confidenceScorer *ConfidenceScorer
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(config Config) *OpenAIClient {
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second // Default timeout
	}

	return &OpenAIClient{
		config: config,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		promptBuilder: NewPromptBuilder(),
		confidenceScorer: NewConfidenceScorer(),
	}
}

// openAIRequest represents the OpenAI API request
type openAIRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIResponse represents the OpenAI API response
type openAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int    `json:"index"`
		Message      message `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// GenerateAnswer generates an answer using OpenAI API
func (c *OpenAIClient) GenerateAnswer(ctx context.Context, req *Request) (*Response, error) {
	// Build messages using prompt builder
	messages := c.promptBuilder.BuildMessages(req)

	// Convert to OpenAI message format
	openAIMessages := make([]message, len(messages))
	for i, msg := range messages {
		openAIMessages[i] = message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Set defaults
	model := req.Model
	if model == "" {
		model = c.config.DefaultModel
		if model == "" {
			model = "gpt-3.5-turbo"
		}
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = c.config.MaxTokens
		if maxTokens == 0 {
			maxTokens = 500
		}
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = c.config.Temperature
		if temperature == 0 {
			temperature = 0.7
		}
	}

	// Create request payload
	payload := openAIRequest{
		Model:       model,
		Messages:    openAIMessages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Determine API endpoint
	baseURL := c.config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	url := fmt.Sprintf("%s/chat/completions", baseURL)

	// Create HTTP request with context
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	// Execute request with retry logic
	var resp *Response
	retryConfig := DefaultRetryConfig()
	retryConfig.MaxRetries = c.config.MaxRetries
	if c.config.RetryDelay > 0 {
		retryConfig.InitialDelay = time.Duration(c.config.RetryDelay) * time.Millisecond
	}

	err = Retry(ctx, func() error {
		httpResp, err := c.httpClient.Do(httpReq)
		if err != nil {
			return &RetryableError{Err: err, Retryable: true}
		}
		defer httpResp.Body.Close()

		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return &RetryableError{Err: fmt.Errorf("failed to read response: %w", err), Retryable: true}
		}

		if httpResp.StatusCode != http.StatusOK {
			// Check if it's a retryable error
			retryable := httpResp.StatusCode >= 500 || httpResp.StatusCode == 429
			return &RetryableError{
				Err:       fmt.Errorf("API error: %d - %s", httpResp.StatusCode, string(body)),
				Retryable: retryable,
			}
		}

		var openAIResp openAIResponse
		if err := json.Unmarshal(body, &openAIResp); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		if openAIResp.Error != nil {
			return fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
		}

		if len(openAIResp.Choices) == 0 {
			return fmt.Errorf("no choices in response")
		}

		choice := openAIResp.Choices[0]
		resp = &Response{
			Content:      choice.Message.Content,
			TokensUsed:   openAIResp.Usage.TotalTokens,
			Model:        openAIResp.Model,
			FinishReason: choice.FinishReason,
		}

		// Calculate confidence score
		resp.Confidence = c.confidenceScorer.CalculateConfidence(resp, req.KnowledgeBase)

		return nil
	}, retryConfig)

	if err != nil {
		return nil, err
	}

	return resp, nil
}


