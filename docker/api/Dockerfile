ARG GO_VERSION=1.15.2

FROM golang:${GO_VERSION}-alpine AS builder
LABEL version="0.1.0"
LABEL maintainer="Haseeb Majid<hello@haseebmajid.dev>"

WORKDIR /api
COPY go.mod go.sum ./
COPY . .

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
RUN go mod download
RUN go build -o ./src ./src/main.go

FROM alpine:latest

WORKDIR /api
COPY --from=builder /api/app .

EXPOSE 8080
ENTRYPOINT ["./app"]