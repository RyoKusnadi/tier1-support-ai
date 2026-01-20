package llm

import (
	"fmt"
	"strings"
)

// PromptBuilder builds prompts for support queries
type PromptBuilder struct {
	systemPrompt string
}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		systemPrompt: `You are a helpful Tier-1 customer support assistant. Your role is to answer customer questions based on the provided knowledge base.

Guidelines:
- Provide clear, concise, and accurate answers
- Base your answers ONLY on the provided knowledge base
- If the knowledge base doesn't contain relevant information, politely indicate that you don't have enough information
- Use a friendly and professional tone
- Keep answers brief and focused
- Do not make up information or speculate beyond what's in the knowledge base`,
	}
}

// BuildMessages constructs the message array for the LLM request
func (pb *PromptBuilder) BuildMessages(req *Request) []Message {
	messages := make([]Message, 0)

	// System message with instructions
	systemContent := pb.systemPrompt

	// Add language-specific instruction if provided
	if req.Language != "" {
		systemContent += fmt.Sprintf("\n\nPlease respond in %s.", getLanguageName(req.Language))
	}

	messages = append(messages, Message{
		Role:    "system",
		Content: systemContent,
	})

	// Add knowledge base context if provided
	if len(req.KnowledgeBase) > 0 {
		knowledgeText := strings.Join(req.KnowledgeBase, "\n\n")
		messages = append(messages, Message{
			Role:    "user",
			Content: fmt.Sprintf("Knowledge Base:\n%s\n\nCustomer Question: %s", knowledgeText, req.Messages[0].Content),
		})
	} else {
		// If no knowledge base, just pass the user question
		if len(req.Messages) > 0 {
			messages = append(messages, req.Messages[0])
		}
	}

	return messages
}

// getLanguageName returns a human-readable language name
func getLanguageName(code string) string {
	langMap := map[string]string{
		"en": "English",
		"es": "Spanish",
		"fr": "French",
		"de": "German",
		"it": "Italian",
		"pt": "Portuguese",
		"ja": "Japanese",
		"ko": "Korean",
		"zh": "Chinese",
		"ar": "Arabic",
		"hi": "Hindi",
		"ru": "Russian",
	}

	if name, ok := langMap[strings.ToLower(code)]; ok {
		return name
	}
	return code
}


