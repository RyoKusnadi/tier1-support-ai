package knowledge

import (
	"context"
	"strings"
)

// Retriever defines the interface for knowledge retrieval.
type Retriever interface {
	// Retrieve returns relevant knowledge snippets (as strings) for a given question.
	// The returned slice is intended to be passed directly into the LLM prompt as context.
	Retrieve(ctx context.Context, tenantID, language, question string) ([]string, error)
}

// InMemoryRetriever is a simple in-memory implementation of Retriever.
type InMemoryRetriever struct {
	documents []Document
}

// NewInMemoryRetriever creates a retriever with a static in-memory corpus.
// This is intended as a stub for Phase 4 and can be replaced by a real store later.
func NewInMemoryRetriever() *InMemoryRetriever {
	return &InMemoryRetriever{
		documents: []Document{
			{
				ID:       "order-status-en-1",
				TenantID: "shop-123",
				Language: "en",
				Title:    "Order status",
				Content:  "Customers can track their order status from the Orders page. Most orders ship within 1-2 business days.",
				Tags:     []string{"order", "shipping", "status"},
			},
			{
				ID:       "refund-policy-en-1",
				TenantID: "shop-123",
				Language: "en",
				Title:    "Refund policy",
				Content:  "Refunds are available within 30 days of delivery for unused items in their original packaging.",
				Tags:     []string{"refund", "returns", "policy"},
			},
		},
	}
}

// Retrieve performs a very simple keyword-based retrieval over the in-memory corpus.
// This is intentionally naive but sufficient as a Phase 4 stub.
func (r *InMemoryRetriever) Retrieve(_ context.Context, tenantID, language, question string) ([]string, error) {
	questionLower := strings.ToLower(question)

	var results []string

	for _, doc := range r.documents {
		// Filter by tenant and language
		if doc.TenantID != tenantID {
			continue
		}
		if doc.Language != "" && language != "" && doc.Language != language {
			continue
		}

		contentLower := strings.ToLower(doc.Content + " " + doc.Title + " " + strings.Join(doc.Tags, " "))

		// Very simple relevance heuristic: any word overlap (after light punctuation stripping)
		words := strings.Fields(questionLower)
		for _, w := range words {
			if len(w) < 3 { // Ignore very short tokens
				continue
			}
			// Trim common punctuation so "order?" still matches "order"
			clean := strings.Trim(w, ".,?!")
			if len(clean) < 3 {
				continue
			}
			if strings.Contains(contentLower, clean) {
				results = append(results, doc.Content)
				break
			}
		}
	}

	return results, nil
}
