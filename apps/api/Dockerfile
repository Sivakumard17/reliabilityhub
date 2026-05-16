# ── Stage 1: Build ────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go.mod first (Docker cache optimization)
COPY apps/api/go.mod apps/api/go.sum ./
RUN GONOSUMDB="*" go mod download

# Copy source
COPY apps/api/ .

# Build static binary
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-X main.Version=${VERSION} -w -s" \
    -o /server ./cmd/server/...

# ── Stage 2: Runtime ──────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

WORKDIR /app
COPY --from=builder /server .

USER appuser

EXPOSE 9090

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:9090/healthz || exit 1

ENTRYPOINT ["./server"]
