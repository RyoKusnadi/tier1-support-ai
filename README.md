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

## Development Roadmap

### Phase 1 — Service Foundation
- [ ] Initialize Go module
- [ ] Basic HTTP server with graceful shutdown
- [ ] Health check endpoint (`/health`)
- [ ] Config loading (env / config file)
- [ ] Structured logging

### Phase 2 — Core API
- [ ] Support query endpoint (`POST /v1/support/query`)
- [ ] Request validation
- [ ] Tenant resolution (initially static config)
- [ ] Language detection (explicit parameter)

### Phase 3 — LLM Integration
- [ ] LLM client abstraction
- [ ] Timeout & retry handling
- [ ] Prompt template for support answers
- [ ] Confidence scoring strategy

### Phase 4 — Knowledge Retrieval
- [ ] Knowledge document model
- [ ] Basic retrieval (in-memory / stub)
- [ ] Retrieval-Augmented Generation (RAG) flow
- [ ] Fallback when no relevant knowledge found

### Phase 5 — Reliability & Cost Control
- [ ] Rate limiting per tenant
- [ ] Response caching
- [ ] Token usage tracking
- [ ] Budget guardrails

### Phase 6 — Observability & Safety
- [ ] Request logging
- [ ] Latency metrics
- [ ] Error tracking
- [ ] Safe fallback for low-confidence responses