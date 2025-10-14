FROM golang:1.25.3-alpine3.22

ENV PATH="/go/bin:${PATH}"

RUN apk add --no-cache gcc git musl-dev && \
    go install github.com/air-verse/air@v1.63.0

# Create non-root user for security
RUN adduser -D -u 1001 developer && \
    chown -R developer:developer /go

USER developer

# No ENTRYPOINT - we'll pass commands via docker-compose run
