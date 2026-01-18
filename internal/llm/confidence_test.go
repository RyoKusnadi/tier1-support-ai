package llm

import (
	"testing"
)

func TestConfidenceScorer_CalculateConfidence(t *testing.T) {
	scorer := NewConfidenceScorer()

	tests := []struct {
		name          string
		response      *Response
		knowledgeBase []string
		wantMin       float64
		wantMax       float64
	}{
		{
			name: "high confidence with knowledge base",
			response: &Response{
				Content:      "Based on our policy, you can return items within 30 days.",
				FinishReason: "stop",
			},
			knowledgeBase: []string{"Return policy: 30 days"},
			wantMin:       0.5,
			wantMax:       1.0,
		},
		{
			name: "low confidence without knowledge base",
			response: &Response{
				Content:      "I'm not sure about that.",
				FinishReason: "stop",
			},
			knowledgeBase: []string{},
			wantMin:       0.0,
			wantMax:       0.5,
		},
		{
			name: "low confidence with uncertainty phrase",
			response: &Response{
				Content:      "I don't know the answer to that question.",
				FinishReason: "stop",
			},
			knowledgeBase: []string{"Some knowledge"},
			wantMin:       0.0,
			wantMax:       0.6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confidence := scorer.CalculateConfidence(tt.response, tt.knowledgeBase)
			if confidence < tt.wantMin || confidence > tt.wantMax {
				t.Errorf("CalculateConfidence() = %f, want between %f and %f", confidence, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestConfidenceScorer_IsHighConfidence(t *testing.T) {
	scorer := NewConfidenceScorer()

	tests := []struct {
		name      string
		confidence float64
		threshold  float64
		want       bool
	}{
		{
			name:       "high confidence above threshold",
			confidence: 0.8,
			threshold:  0.7,
			want:       true,
		},
		{
			name:       "low confidence below threshold",
			confidence: 0.5,
			threshold:  0.7,
			want:       false,
		},
		{
			name:       "confidence at threshold",
			confidence: 0.7,
			threshold:  0.7,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scorer.IsHighConfidence(tt.confidence, tt.threshold)
			if got != tt.want {
				t.Errorf("IsHighConfidence() = %v, want %v", got, tt.want)
			}
		})
	}
}

