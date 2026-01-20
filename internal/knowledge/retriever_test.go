package knowledge

import (
	"context"
	"testing"
)

func TestInMemoryRetriever_Retrieve(t *testing.T) {
	retriever := NewInMemoryRetriever()

	tests := []struct {
		name     string
		tenantID string
		language string
		question string
		wantMin  int
	}{
		{
			name:     "matching order question",
			tenantID: "shop-123",
			language: "en",
			question: "Where is my order?",
			wantMin:  1,
		},
		{
			name:     "unrelated tenant",
			tenantID: "other-tenant",
			language: "en",
			question: "Where is my order?",
			wantMin:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			results, err := retriever.Retrieve(ctx, tt.tenantID, tt.language, tt.question)
			if err != nil {
				t.Fatalf("Retrieve() error = %v", err)
			}
			if len(results) < tt.wantMin {
				t.Fatalf("Retrieve() len(results) = %d, want at least %d", len(results), tt.wantMin)
			}
		})
	}
}
