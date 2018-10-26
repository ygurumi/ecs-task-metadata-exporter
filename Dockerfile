FROM golang:alpine AS builder

COPY . /go/src/github.com/ygurumi/ecs-task-metadata-exporter
WORKDIR /go/src/github.com/ygurumi/ecs-task-metadata-exporter
RUN apk add --no-cache git
RUN go get -u github.com/golang/dep/cmd/dep && \
    dep ensure && \
    go build -ldflags '-s -w'

FROM alpine:latest
WORKDIR /
COPY --from=builder /go/src/github.com/ygurumi/ecs-task-metadata-exporter/ecs-task-metadata-exporter .
CMD ["/ecs-task-metadata-exporter"]
