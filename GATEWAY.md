# Gateway Integration Strategy: SS-CatalogService

This document outlines the architectural decisions for integrating the Catalog Service with the SAMSTORE API Gateway.

## Routing Strategy
- **Path Pattern**: `/api/catalog/v1/*`
- **Prefix Strategy**: **Full Path Forwarding**. The Gateway forwards the complete path starting with `/api/catalog`. The service handles this prefix natively in its router.
- **Base URL**: `https://<gateway-domain>/api/catalog/v1`

## Authentication (Option A)
- **Method**: **Local JWT Validation** (RS256).
- **Public Key**: The service validates JWTs using the same RSA Public Key as the Auth Service.
- **Trust Model**: The service does NOT trust `X-User-*` headers from the gateway. It strictly validates the incoming `Authorization: Bearer <JWT>` header. This approach provides stronger zero-trust security.

## Observability
- **Trace Propagation**: Supports W3C `traceparent`.
- **Correlation ID**: Propagates `X-Correlation-Id` across requests.
- **Health Checks**: Exposes `GET /health` for active load balancer / gateway health monitoring.

## Caching
- **Public Routes**: Emits `Cache-Control` and `ETag`.
- **Conditional Requests**: Supports `If-None-Match` for `304 Not Modified` responses.
- **Vary Header**: Includes `Vary: Accept-Language, Authorization` to prevent cache poisoning across different user contexts.
