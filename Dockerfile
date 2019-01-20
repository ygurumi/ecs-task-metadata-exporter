FROM golang:alpine AS builder
COPY . /ecs-task-metadata-exporter
WORKDIR /ecs-task-metadata-exporter
RUN apk add --no-cache git musl-dev gcc g++
RUN go build -ldflags '-s -w'

FROM alpine:latest
WORKDIR /
COPY --from=builder /ecs-task-metadata-exporter/ecs-task-metadata-exporter .
CMD ["/ecs-task-metadata-exporter"]
