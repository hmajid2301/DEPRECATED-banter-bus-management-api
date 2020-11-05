ARG GO_VERSION=1.15.2

FROM golang:${GO_VERSION}-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    BANTER_BUS_CONFIG_PATH=/go/src/src/config.yml

WORKDIR /go/src/
COPY go.mod go.sum docker/api/install.sh ./
RUN go mod download && \
    chmod +x /go/src/install.sh

EXPOSE 8080 2345
RUN go get github.com/go-delve/delve/cmd/dlv && \
    go get github.com/githubnemo/CompileDaemon 
ADD . /go/src/

WORKDIR /go/src/src/
ENTRYPOINT /go/src/install.sh