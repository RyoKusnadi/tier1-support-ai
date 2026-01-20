package knowledge

// Document represents a knowledge base document that can be used for retrieval.
type Document struct {
	ID       string   // Unique identifier
	TenantID string   // Tenant that owns this document
	Language string   // ISO 639-1 language code (e.g., "en")
	Title    string   // Optional title or short label
	Content  string   // Main text content
	Tags     []string // Optional tags / categories
}
