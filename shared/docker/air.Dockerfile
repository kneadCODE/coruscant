FROM golang:1.25-alpine3.22

ENV PATH="/go/bin:${PATH}"

RUN apk add --no-cache gcc git musl-dev && \
    go install github.com/air-verse/air@v1.62.0

# No ENTRYPOINT - we'll pass commands via docker-compose run
