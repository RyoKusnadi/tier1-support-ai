# Tier-1 Support AI

## Overview
Tier-1 Support AI is a production-ready backend service designed to handle first-level customer support questions using Large Language Models (LLMs). The system implements a complete Retrieval-Augmented Generation (RAG) pipeline with comprehensive reliability, cost control, and observability features.

Built as a multi-tenant and multilingual service, it focuses on safe AI integration for production environments with robust fallback mechanisms, rate limiting, budget controls, and comprehensive monitoring.

## Problem Statement
Customer support teams often receive a large volume of repetitive questions about order status, refunds, delivery times, and basic policies. This service automates Tier-1 support responses using AI, allowing human agents to focus on complex and sensitive cases that require human judgment.

## Key Features

### ✅ Multi-Tenant Architecture
- Isolated tenant configurations and knowledge bases
- Per-tenant rate limiting and budget controls
- Tenant-specific token usage tracking

### ✅ Multilingual Support
- Language-aware prompt generation
- Configurable supported languages (currently: English, Indonesian)
- Language-specific knowledge retrieval

### ✅ Knowledge-Based RAG Pipeline
- In-memory knowledge document storage with tenant isolation
- Keyword-based retrieval with relevance scoring
- Context-aware prompt building with retrieved knowledge
- Safe fallback when no relevant knowledge is found

### ✅ Production-Ready Reliability
- **Rate Limiting**: Per-tenant token bucket rate limiter (5 req/sec, burst 10)
- **Response Caching**: TTL-based in-memory cache (5min default)
- **Budget Controls**: Per-tenant token budget enforcement with sliding windows
- **Retry Logic**: Exponential backoff with configurable retries
- **Graceful Degradation**: Confidence-based fallback responses

### ✅ Comprehensive Observability
- Structured request logging with correlation IDs
- In-process metrics (requests, errors, latency, cache hits/misses)
- Error tracking with standardized error codes
- Token usage monitoring per tenant

### ✅ LLM Integration
- OpenAI API integration with configurable models
- Confidence scoring based on response analysis
- Timeout and retry handling
- Token usage tracking for cost control

## Architecture

### Core Components

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│  Support Handler │───▶│   LLM Client    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │ Knowledge Store  │    │  OpenAI API     │
                       └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │ Reliability Layer│
                       │ • Rate Limiter   │
                       │ • Cache          │
                       │ • Budget Guard   │
                       │ • Usage Tracker  │
                       └──────────────────┘
```

### Request Flow
1. **Request Validation**: Validate tenant_id, language, and question
2. **Rate Limiting**: Check per-tenant rate limits
3. **Cache Check**: Look for cached responses
4. **Budget Validation**: Verify tenant hasn't exceeded token budget
5. **Knowledge Retrieval**: Find relevant documents using keyword matching
6. **LLM Generation**: Generate response using RAG pipeline
7. **Confidence Scoring**: Calculate confidence based on response analysis
8. **Fallback Logic**: Return safe fallback if confidence < 0.7
9. **Usage Tracking**: Record token usage for budget enforcement
10. **Response Caching**: Cache successful responses

## Tech Stack

### Core Technologies
- **Go 1.24.4**: Backend service with Gin web framework
- **OpenAI API**: LLM provider (configurable for other providers)
- **In-Memory Storage**: Knowledge base, caching, and rate limiting
- **JSON**: API serialization and configuration

### Dependencies
- `github.com/gin-gonic/gin`: HTTP web framework
- `go.uber.org/mock`: Testing mocks
- Standard library for HTTP client, JSON, crypto, etc.

## Configuration

The service is configured entirely through environment variables:

### Core Service Configuration
```bash
PORT=8080                    # Server port (default: 8080)
APP_ENV=development          # Environment (default: development)
```

### LLM Configuration
```bash
LLM_PROVIDER=openai          # LLM provider (default: openai)
LLM_API_KEY=sk-...           # OpenAI API key (required)
LLM_BASE_URL=                # Custom API endpoint (optional)
LLM_DEFAULT_MODEL=gpt-3.5-turbo  # Default model (default: gpt-3.5-turbo)
LLM_MAX_TOKENS=500           # Max response tokens (default: 500)
LLM_TEMPERATURE=0.7          # Response creativity (default: 0.7)
LLM_TIMEOUT=30               # Request timeout seconds (default: 30)
LLM_MAX_RETRIES=3            # Max retry attempts (default: 3)
LLM_RETRY_DELAY=100          # Initial retry delay ms (default: 100)
```

### Reliability & Cost Control
```bash
TENANT_RATE_LIMIT_PER_SEC=5.0     # Requests per second per tenant (default: 5.0)
TENANT_RATE_LIMIT_BURST=10        # Burst capacity (default: 10)
RESPONSE_CACHE_TTL_SECONDS=300    # Cache TTL in seconds (default: 300)
TOKEN_USAGE_WINDOW_HOURS=24       # Usage tracking window (default: 24)
TENANT_TOKEN_BUDGET=0             # Per-tenant token budget, 0=disabled (default: 0)
```

### Tenant & Language Configuration
Tenants and supported languages are currently configured in code:

**Supported Tenants** (`internal/config/tenant.go`):
- `shop-123`
- `shop-456`

**Supported Languages** (`internal/config/language.go`):
- `en` (English)
- `id` (Indonesian)

## Development

### Prerequisites
- Go 1.24.4 or later
- OpenAI API key (for LLM integration)

### Quick Start
```bash
# Clone the repository
git clone https://github.com/RyoKusnadi/tier1-support-ai
cd tier1-support-ai

# Set required environment variables
export LLM_API_KEY=your-openai-api-key

# Run the service
make run

# Or build and run
make build
./bin/tier1-support-ai
```

### Development Commands
```bash
make run     # Run the service locally
make build   # Build binary to bin/tier1-support-ai
make test    # Run all tests
make clean   # Remove build artifacts
```

### Testing
The project includes comprehensive unit tests:
```bash
go test ./...                    # Run all tests
go test -v ./internal/llm/...    # Run LLM package tests with verbose output
go test -cover ./...             # Run tests with coverage
```

Test coverage includes:
- Configuration loading and defaults
- LLM client initialization and provider selection
- Confidence scoring algorithms
- Prompt building with knowledge base integration
- Retry logic with exponential backoff
- Knowledge retrieval with tenant isolation


## API Documentation

### Base URL
```
http://localhost:8080
```

### Authentication
Currently no authentication is required. In production, consider adding API key authentication or OAuth.

### Endpoints

#### Health Check
```http
GET /health
```
Returns service health status.

**Response:**
```
OK
```

#### Metrics
```http
GET /metrics
```
Returns operational metrics in JSON format.

**Response:**
```json
{
  "requests_total": 1250,
  "errors_total": 23,
  "rate_limited_total": 5,
  "budget_blocked_total": 2,
  "cache_hits_total": 340,
  "cache_misses_total": 910,
  "latency_count": 1250,
  "latency_sum_ms": 45000,
  "latency_avg_ms": 36.0
}
```

#### Support Query
```http
POST /v1/support/query
```
Handles Tier-1 customer support questions using AI-assisted responses with RAG pipeline.

**Request Headers:**
```
Content-Type: application/json
X-Request-Id: optional-request-id  # Auto-generated if not provided
```

**Request Body:**
```json
{
  "tenant_id": "shop-123",
  "language": "en",
  "question": "Where is my order?",
  "knowledge_base": ["Optional additional context"]
}
```

| Field          | Type     | Required | Description                           |
|----------------|----------|----------|---------------------------------------|
| tenant_id      | string   | Yes      | Tenant identifier (shop-123, shop-456)|
| language       | string   | Yes      | Language code (en, id)                |
| question       | string   | Yes      | Customer question                     |
| knowledge_base | []string | No       | Additional context documents          |

**Success Response (200 OK):**
```json
{
  "answer": "Your order is on the way and will arrive tomorrow.",
  "confidence": 0.87,
  "tenant_id": "shop-123",
  "language": "en"
}
```

**Fallback Response (200 OK) - Low Confidence:**
```json
{
  "answer": "We are unable to confidently answer your question. Please contact customer support.",
  "confidence": 0.32,
  "tenant_id": "shop-123",
  "language": "en",
  "fallback": true
}
```

**Error Responses:**

**400 Bad Request - Invalid Request:**
```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Invalid request: language is required"
  }
}
```

**429 Too Many Requests - Rate Limited:**
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded for tenant"
  }
}
```

**429 Too Many Requests - Budget Exceeded:**
```json
{
  "error": {
    "code": "BUDGET_EXCEEDED",
    "message": "Token budget exceeded for tenant"
  }
}
```

**500 Internal Server Error:**
```json
{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "Failed to generate answer"
  }
}
```

### Response Fields

| Field      | Type    | Description                                    |
|------------|---------|------------------------------------------------|
| answer     | string  | AI-generated response or fallback message     |
| confidence | number  | Confidence score (0.0–1.0)                    |
| tenant_id  | string  | Echo of request tenant_id                     |
| language   | string  | Echo of request language                      |
| fallback   | boolean | Present and true when using fallback response |

### Error Codes

| Code               | HTTP Status | Description                    |
|--------------------|-------------|--------------------------------|
| INVALID_REQUEST    | 400         | Request validation failed      |
| RATE_LIMIT_EXCEEDED| 429         | Tenant rate limit exceeded     |
| BUDGET_EXCEEDED    | 429         | Tenant token budget exceeded   |
| INTERNAL_ERROR     | 500         | Unexpected server error        |

## Implementation Details

### LLM Integration

**OpenAI Client** (`internal/llm/openai.go`):
- Configurable API endpoint and model selection
- Automatic retry with exponential backoff
- Context-aware request handling with timeouts
- Token usage tracking for cost monitoring
- Error handling with retryable/non-retryable classification

**Confidence Scoring** (`internal/llm/confidence.go`):
- Base confidence of 0.7 for all responses
- Penalty for uncertainty phrases ("I don't know", "maybe", etc.)
- Penalty for missing knowledge base context (-0.3)
- Penalty for very short responses (-0.1)
- Penalty for truncated responses (-0.1)
- Final score clamped to [0.0, 1.0] range

**Prompt Engineering** (`internal/llm/prompt.go`):
- System prompt with clear instructions for support role
- Knowledge base integration in user message
- Language-specific response instructions
- Structured message building for RAG pipeline

### Knowledge Retrieval

**In-Memory Retriever** (`internal/knowledge/retriever.go`):
- Tenant-isolated document storage
- Language-aware filtering
- Simple keyword-based relevance matching
- Configurable knowledge base with tags and metadata

**Sample Knowledge Base**:
```go
{
    ID:       "order-status-en-1",
    TenantID: "shop-123",
    Language: "en",
    Title:    "Order status",
    Content:  "Customers can track their order status from the Orders page. Most orders ship within 1-2 business days.",
    Tags:     []string{"order", "shipping", "status"},
}
```

### Reliability Layer

**Rate Limiting** (`internal/reliability/limiter.go`):
- Per-tenant token bucket implementation
- Configurable rate (requests/second) and burst capacity
- Thread-safe with mutex protection
- Automatic token refill based on elapsed time

**Response Caching** (`internal/reliability/cache.go`):
- Generic TTL-based in-memory cache
- Thread-safe with read-write mutex
- Automatic expiration cleanup
- Cache key: `tenant_id|language|question`

**Budget Control** (`internal/reliability/budget.go`):
- Per-tenant token budget enforcement
- Sliding window usage tracking
- Pre-call budget validation
- Remaining budget calculation with reset times

**Usage Tracking** (`internal/reliability/usage.go`):
- Per-tenant token consumption monitoring
- Sliding window with configurable duration
- Request count and token usage aggregation
- Automatic window reset and cleanup

### Observability

**Metrics Collection** (`internal/observability/metrics.go`):
- In-process atomic counters for key metrics
- Request latency tracking with count and sum
- Cache hit/miss ratios
- Error and rate limiting counters
- JSON metrics endpoint at `/metrics`

**Request Logging** (`internal/middleware/requestlog.go`):
- Structured logging with correlation IDs
- Request completion logging with latency
- Tenant ID attachment when available
- Client IP and User-Agent tracking

**Error Handling**:
- Standardized error response format
- Categorized error codes for different failure modes
- Comprehensive error logging with context
- Graceful degradation for non-critical failures

### Configuration Management

**Environment-Based Config** (`internal/config/config.go`):
- Type-safe configuration loading
- Sensible defaults for all parameters
- Helper functions for int/float parsing
- Centralized configuration structure

**Static Configuration**:
- Tenant whitelist in `internal/config/tenant.go`
- Supported languages in `internal/config/language.go`
- Easy to extend for dynamic configuration sources

## Example Usage

### Basic Support Query
```bash
curl -X POST http://localhost:8080/v1/support/query \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "shop-123",
    "language": "en",
    "question": "What is your return policy?"
  }'
```

**Response:**
```json
{
  "answer": "Refunds are available within 30 days of delivery for unused items in their original packaging.",
  "confidence": 0.85,
  "tenant_id": "shop-123",
  "language": "en"
}
```

### Query with Additional Context
```bash
curl -X POST http://localhost:8080/v1/support/query \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "shop-123",
    "language": "en",
    "question": "How long does shipping take?",
    "knowledge_base": ["Express shipping: 1-2 days", "Standard shipping: 3-5 days"]
  }'
```

### Low Confidence Fallback
```bash
curl -X POST http://localhost:8080/v1/support/query \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "shop-123",
    "language": "en",
    "question": "Can you help me with my tax filing?"
  }'
```

**Response:**
```json
{
  "answer": "We are unable to confidently answer your question. Please contact customer support.",
  "confidence": 0.1,
  "tenant_id": "shop-123",
  "language": "en",
  "fallback": true
}
```

### Check Service Health
```bash
curl http://localhost:8080/health
# Response: OK

curl http://localhost:8080/metrics
# Response: JSON metrics object
```

## Production Considerations

### Security
- **API Authentication**: Add API key or OAuth authentication
- **Rate Limiting**: Current in-memory rate limiting should be replaced with Redis for multi-instance deployments
- **Input Validation**: Implement additional input sanitization and length limits
- **CORS**: Configure CORS policies for web client access

### Scalability
- **Horizontal Scaling**: Replace in-memory components (cache, rate limiter) with Redis
- **Database**: Replace in-memory knowledge base with PostgreSQL or vector database
- **Load Balancing**: Deploy behind load balancer with health checks
- **Caching**: Add CDN or reverse proxy caching for static responses

### Monitoring & Alerting
- **External Metrics**: Integrate with Prometheus/Grafana or DataDog
- **Distributed Tracing**: Add OpenTelemetry for request tracing
- **Log Aggregation**: Ship logs to ELK stack or similar
- **Alerting**: Set up alerts for error rates, latency, and budget exhaustion

### Cost Optimization
- **Model Selection**: Use cheaper models for simple queries
- **Response Caching**: Increase cache TTL for stable knowledge
- **Batch Processing**: Implement request batching for high-volume tenants
- **Budget Alerts**: Add proactive budget monitoring and alerts

### Knowledge Management
- **Dynamic Knowledge Base**: Implement API for knowledge document management
- **Vector Search**: Replace keyword matching with semantic search
- **Knowledge Versioning**: Add versioning and rollback capabilities
- **Multi-Modal**: Support for images and documents in knowledge base

## Troubleshooting

### Common Issues

**Service won't start:**
```bash
# Check if LLM_API_KEY is set
echo $LLM_API_KEY

# Check port availability
lsof -i :8080

# Check logs for configuration errors
make run
```

**High error rates:**
```bash
# Check metrics endpoint
curl http://localhost:8080/metrics

# Look for rate limiting or budget issues
# Check LLM API key validity and quotas
```

**Poor response quality:**
```bash
# Verify knowledge base content
# Check confidence threshold (default: 0.7)
# Review prompt templates in internal/llm/prompt.go
```

**Performance issues:**
```bash
# Check cache hit rates in metrics
# Monitor LLM API latency
# Review rate limiting configuration
```

### Debugging

**Enable verbose logging:**
```bash
# Add debug logging in internal/logger/logger.go
# Set APP_ENV=development for more detailed logs
```

**Test LLM integration:**
```bash
# Test with minimal request
curl -X POST http://localhost:8080/v1/support/query \
  -H "Content-Type: application/json" \
  -d '{"tenant_id":"shop-123","language":"en","question":"test"}'
```

**Monitor token usage:**
```bash
# Check metrics for budget_blocked_total
# Review token usage in logs after LLM calls
```

## Contributing

### Development Setup
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make changes and add tests
4. Run tests: `make test`
5. Submit a pull request

### Code Style
- Follow Go conventions and `gofmt` formatting
- Add unit tests for new functionality
- Update documentation for API changes
- Use structured logging with appropriate log levels

### Testing
- Maintain test coverage above 80%
- Add integration tests for new endpoints
- Mock external dependencies (LLM APIs)
- Test error conditions and edge cases

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.