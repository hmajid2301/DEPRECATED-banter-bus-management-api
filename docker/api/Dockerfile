ARG GO_VERSION=1.16

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /temp

COPY go.mod go.sum config.yml ./
COPY internal/ ./internal/
COPY cmd/ ./cmd/

RUN go mod download && \
    go build -o ./app ./cmd/banter-bus-management-api/main.go

FROM alpine:latest

EXPOSE 8080

WORKDIR /api
COPY --from=builder /temp/app ./
COPY --from=builder /temp/config.yml ./config.yml

ENTRYPOINT ["./app"]