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

WORKDIR /workspace

COPY go.work ./

COPY shared/golib shared/golib/

COPY systems/${SYSTEM} systems/${SYSTEM}/

WORKDIR /workspace/systems/${SYSTEM}

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && \
    go mod verify

# Build all binaries from cmd directory
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    for cmd in $(ls cmd); do \
    if [ -f "cmd/${cmd}/main.go" ]; then \
    echo "Building ${cmd}..." && \
    go build -v -ldflags=" \
    -X main.version=${VERSION} \
    -X main.commitSha=${COMMIT_SHA} \
    -X main.buildTime=${BUILD_TIME}" \
    -o "/workspace/bin/${cmd}" "./cmd/${cmd}"; \
    fi \
    done

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

COPY --from=builder /workspace/bin/ /usr/local/bin/

# OCI labels are automatically added by GitHub Actions docker/metadata-action

USER 10001:10001

# No default CMD as this is a multi-binary image.
# The specific binary to run should be specified in the k8s deployment
