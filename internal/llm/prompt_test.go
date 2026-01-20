package llm

import (
	"testing"
)

func TestPromptBuilder_BuildMessages(t *testing.T) {
	builder := NewPromptBuilder()

	tests := []struct {
		name     string
		request  *Request
		wantLen  int
		wantRole string
	}{
		{
			name: "with knowledge base",
			request: &Request{
				Messages: []Message{
					{Role: "user", Content: "What is the return policy?"},
				},
				KnowledgeBase: []string{"Return policy: 30 days", "Refund: Full refund"},
				Language:      "en",
			},
			wantLen:  2,
			wantRole: "system",
		},
		{
			name: "without knowledge base",
			request: &Request{
				Messages: []Message{
					{Role: "user", Content: "What is the return policy?"},
				},
				Language: "en",
			},
			wantLen:  2,
			wantRole: "system",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := builder.BuildMessages(tt.request)
			if len(messages) != tt.wantLen {
				t.Errorf("BuildMessages() got %d messages, want %d", len(messages), tt.wantLen)
			}
			if len(messages) > 0 && messages[0].Role != tt.wantRole {
				t.Errorf("BuildMessages() first message role = %s, want %s", messages[0].Role, tt.wantRole)
			}
		})
	}
}


