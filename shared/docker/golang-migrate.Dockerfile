# syntax=docker/dockerfile:1.18

ARG RUNTIME_IMAGE_TAG=nonroot

# ----------------------------------------------------------------
# Builder stage
# ----------------------------------------------------------------
FROM golang:1.25.3-alpine3.22 AS builder

ARG SYSTEM

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go install -v github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

# ----------------------------------------------------------------
# Dependencies stage
# ----------------------------------------------------------------
FROM debian:bookworm-slim AS deps

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata && \
    rm -rf /var/lib/apt/lists/*

# ----------------------------------------------------------------
# Runtime stage
# ----------------------------------------------------------------

FROM gcr.io/distroless/static-debian12:${RUNTIME_IMAGE_TAG} AS runtime

ARG SYSTEM

ENV TZ=UTC

COPY --from=deps /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=deps /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=deps /etc/timezone /etc/timezone
COPY --from=deps /etc/localtime /etc/localtime

COPY --from=builder /go/bin/migrate /usr/local/bin/
COPY systems/${SYSTEM}/db-migrations /db-migrations

# OCI labels are automatically added by GitHub Actions docker/metadata-action

USER 10001:10001

ENTRYPOINT ["/usr/local/bin/migrate"]
