package llm

import (
	"math"
	"strings"
)

// ConfidenceScorer calculates confidence scores for LLM responses
type ConfidenceScorer struct {
	lowConfidencePhrases []string
}

// NewConfidenceScorer creates a new confidence scorer
func NewConfidenceScorer() *ConfidenceScorer {
	return &ConfidenceScorer{
		lowConfidencePhrases: []string{
			"i don't know",
			"i'm not sure",
			"i cannot",
			"i'm unable",
			"i don't have",
			"no information",
			"not available",
			"unclear",
			"uncertain",
			"might be",
			"could be",
			"possibly",
			"perhaps",
			"maybe",
		},
	}
}

// CalculateConfidence calculates a confidence score for the response
// Returns a score between 0.0 and 1.0
func (cs *ConfidenceScorer) CalculateConfidence(response *Response, knowledgeBase []string) float64 {
	content := strings.ToLower(response.Content)

	// Base confidence starts at 0.7 (moderate confidence)
	confidence := 0.7

	// Check for low confidence phrases
	for _, phrase := range cs.lowConfidencePhrases {
		if strings.Contains(content, phrase) {
			confidence -= 0.2
			break
		}
	}

	// Adjust based on knowledge base availability
	if len(knowledgeBase) == 0 {
		confidence -= 0.3 // Lower confidence if no knowledge base provided
	}

	// Adjust based on response length (very short responses might be incomplete)
	if len(content) < 20 {
		confidence -= 0.1
	}

	// Adjust based on finish reason
	if response.FinishReason == "length" {
		confidence -= 0.1 // Response was truncated
	}

	// Ensure confidence is within bounds
	confidence = math.Max(0.0, math.Min(1.0, confidence))

	return confidence
}

// IsHighConfidence checks if the confidence score is above the threshold
func (cs *ConfidenceScorer) IsHighConfidence(confidence float64, threshold float64) bool {
	return confidence >= threshold
}

