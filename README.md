# Tier-1 Support AI

## Overview
Tier-1 Support AI is a backend service designed to handle first-level customer support questions
(e.g. order status, refunds, delivery time) using Large Language Models (LLMs).

The system is built as a multi-tenant and multilingual service, focusing on reliability,
cost control, and safe AI integration for production use.

## Problem Statement
Customer support teams often receive a large volume of repetitive questions.
This service aims to automate Tier-1 support responses so that human agents can
focus on complex and sensitive cases.

## Scope
- Multi-tenant support
- Multilingual question handling
- Knowledge-based answers (not free-form chat)
- Safe fallback when confidence is low

## Non-Goals
- Replacing human support entirely
- Handling legal or financial disputes
- Training custom AI models

## High-Level Flow
1. Receive customer question
2. Identify tenant and language
3. Retrieve relevant knowledge
4. Generate AI-assisted answer
5. Return answer with confidence score

## Tech Stack (Planned)
- Go (backend service)
- REST API
- External LLM API
- Redis (caching, rate limiting)
- PostgreSQL (metadata)

## API Versioning
This service uses URL-based API versioning.

- Current version: `v1`
- Example endpoint: `POST /v1/support/query`

New versions will be introduced only when breaking changes are required,
while maintaining backward compatibility for existing clients.

## Development

Common development commands:

```bash
make run    # run the service locally
make build  # build binary
make test   # run tests
```

Configuration is provided via environment variables.


## API Contract (v1)

### POST /v1/support/query
Handles Tier-1 customer support questions using AI-assisted responses.

### Request

```json
{
  "tenant_id": "shop-123",
  "language": "en",
  "question": "Where is my order?"
}
```

| Field      | Type   | Required | Description               |
|------------|--------|----------|---------------------------|
| tenant_id  | string | Yes      | Tenant identifier         |
| language   | string | Yes      | Language code (ISO 639-1) |
| question   | string | Yes      | Customer question         |

### Response — Success

```json
{
  "answer": "Your order is on the way and will arrive tomorrow.",
  "confidence": 0.87
}
```

| Field      | Type   | Description                |
|------------|--------|----------------------------|
| answer     | string | AI-generated response      |
| confidence | number | Confidence score (0.0–1.0) |

### Response — Fallback (Low Confidence)

```json
{
  "answer": "We are unable to confidently answer your question. Please contact customer support.",
  "confidence": 0.32,
  "fallback": true
}
```

| Field    | Type    | Description                 |
|----------|---------|-----------------------------|
| fallback | boolean | Indicates fallback response |

### Response — Error

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "language is required"
  }
}
```

| Code             | Description               |
|------------------|---------------------------|
| INVALID_REQUEST  | Request validation failed |
| TENANT_NOT_FOUND | Unknown tenant            |
| INTERNAL_ERROR   | Unexpected server error   |

## Development Roadmap

### Phase 1 — Service Foundation
- [x] Initialize Go module
- [x] Basic HTTP server with graceful shutdown
- [x] Health check endpoint (`/health`)
- [x] Config loading (env / config file)
- [x] Structured logging


### Phase 2 — Core API
- [X] Support query endpoint (`POST /v1/support/query`)
- [X] Request validation
- [X] Tenant resolution (initially static config)
- [X] Language detection (explicit parameter)

### Phase 3 — LLM Integration
- [x] LLM client abstraction
- [x] Timeout & retry handling
- [x] Prompt template for support answers
- [x] Confidence scoring strategy

### Phase 4 — Knowledge Retrieval
- [x] Knowledge document model
- [x] Basic retrieval (in-memory / stub)
- [x] Retrieval-Augmented Generation (RAG) flow
- [x] Fallback when no relevant knowledge found

### Phase 5 — Reliability & Cost Control
- [x] Rate limiting per tenant
- [x] Response caching
- [x] Token usage tracking
- [x] Budget guardrails

### Phase 6 — Observability & Safety
- [x] Request logging
- [x] Latency metrics
- [x] Error tracking
- [x] Safe fallback for low-confidence responses

## Reliability & Cost Control (Phase 5)

Phase 5 adds a first pass of production-oriented safeguards around LLM usage:

- **Per-tenant rate limiting**: simple in-memory token-bucket limiter that returns `429 RATE_LIMIT_EXCEEDED` when a tenant sends too many requests per second.
- **Response caching**: in-memory TTL cache for identical `(tenant_id, language, question)` triples, reducing duplicate LLM calls and latency.
- **Token usage tracking**: per-tenant counters of `TokensUsed` over a sliding window, updated after each successful LLM call.
- **Budget guardrails**: optional per-tenant token budget; once exceeded within the current window, further LLM calls are blocked with `429 BUDGET_EXCEEDED`.

Phase 5 configuration (environment variables):

- `TENANT_RATE_LIMIT_PER_SEC` (float, default `5.0`)
- `TENANT_RATE_LIMIT_BURST` (int, default `10`)
- `RESPONSE_CACHE_TTL_SECONDS` (int, default `300`)
- `TOKEN_USAGE_WINDOW_HOURS` (int, default `24`)
- `TENANT_TOKEN_BUDGET` (int, default `0` = disabled)

## Observability & Safety (Phase 6)

Phase 6 adds basic operational visibility without external dependencies:

- **Request logging**: middleware logs one line per request with `request_id`, `tenant_id` (when available), status, and latency.
- **Latency metrics**: in-process counters track total requests and average latency.
- **Error tracking**: standardized error payloads (with `code` + `message`) plus in-process error counters.
- **Safe fallback**: low-confidence responses are returned as the safe fallback answer (enforced in the support handler).

Metrics are available at `GET /metrics` as JSON.
