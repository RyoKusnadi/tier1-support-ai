package llm

import "context"

// Request represents a request to the LLM
type Request struct {
	Messages      []Message
	MaxTokens     int
	Temperature   float64
	Model         string
	KnowledgeBase []string // Retrieved knowledge documents for RAG
	Language      string   // Language code (e.g., "en", "es")
	TenantID      string   // Multi-tenant support
}

// Message represents a single message in the conversation
type Message struct {
	Role    string // "system", "user", "assistant"
	Content string
}

// Response represents the LLM response
type Response struct {
	Content        string  // The generated answer
	Confidence     float64 // Confidence score (0.0 to 1.0)
	TokensUsed     int     // Number of tokens consumed
	Model          string  // Model used
	FinishReason   string  // Reason for completion (e.g., "stop", "length")
}

// Client defines the interface for LLM clients
type Client interface {
	// GenerateAnswer generates an answer based on the request
	GenerateAnswer(ctx context.Context, req *Request) (*Response, error)
}

// Config holds LLM client configuration
type Config struct {
	Provider      string  // "openai", "anthropic", etc.
	APIKey        string
	BaseURL       string  // Optional, for custom endpoints
	DefaultModel  string
	MaxTokens     int
	Temperature   float64
	Timeout       int     // Timeout in seconds
	MaxRetries    int     // Maximum number of retries
	RetryDelay    int     // Initial retry delay in milliseconds
}


