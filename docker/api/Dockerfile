ARG GO_VERSION=1.15.2

FROM golang:${GO_VERSION}-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /api
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -directory=/api/src  -log-prefix=false -build="go build -a -installsuffix cgo -o /go/bin/hello" -command="/go/bin/hello"
