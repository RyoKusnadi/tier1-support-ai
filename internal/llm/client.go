package llm

import (
	"fmt"
)

// NewClient creates a new LLM client based on the provider
func NewClient(config Config) (Client, error) {
	switch config.Provider {
	case "openai":
		return NewOpenAIClient(config), nil
	case "":
		// Default to OpenAI if not specified
		return NewOpenAIClient(config), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", config.Provider)
	}
}
