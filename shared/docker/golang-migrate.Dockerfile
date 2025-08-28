# syntax=docker/dockerfile:1.17

ARG RUNTIME_IMAGE_TAG=nonroot

# ----------------------------------------------------------------
# Builder stage
# ----------------------------------------------------------------
FROM golang:1.25-alpine3.22 AS builder

ARG SYSTEM

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go install -v github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

COPY systems/${SYSTEM}/db-migrations /db-migrations

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

ARG ORG
ARG REPO_NAME
ARG SYSTEM
ARG VERSION
ARG COMMIT_SHA
ARG BUILD_TIME
ARG REPO_URL
ARG LICENSE

ENV TZ=UTC

COPY --from=deps /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=deps /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=deps /etc/timezone /etc/timezone
COPY --from=deps /etc/localtime /etc/localtime

COPY --from=builder /go/bin/linux_amd64/migrate /usr/local/bin/
COPY --from=builder /db-migrations /migrations

# OCI standard labels (https://github.com/opencontainers/image-spec/blob/main/annotations.md)
LABEL org.opencontainers.image.title="${REPO_NAME} ${SYSTEM} Service" \
    org.opencontainers.image.description="${REPO_NAME} ${SYSTEM} system golang-migrate" \
    org.opencontainers.image.url=${REPO_URL} \
    org.opencontainers.image.source=${REPO_URL} \
    org.opencontainers.image.documentation="${REPO_URL}/tree/main/README.md" \
    org.opencontainers.image.version=${VERSION} \
    org.opencontainers.image.revision=${COMMIT_SHA} \
    org.opencontainers.image.created=${BUILD_TIME} \
    org.opencontainers.image.authors=${ORG} \
    org.opencontainers.image.licenses=${LICENSE} \
    org.opencontainers.image.vendor=${ORG} \
    org.opencontainers.image.base.name="gcr.io/distroless/static-debian12:${RUNTIME_IMAGE_TAG}" \
    org.opencontainers.image.ref.name="${ORG}/${SYSTEM}:${VERSION}"

USER 10001:10001

ENTRYPOINT ["/usr/local/bin/migrate"]
