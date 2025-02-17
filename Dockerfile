FROM golang:1.23.6-bookworm AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/app ./cmd/taskscheduler

FROM debian:bookworm-slim

WORKDIR /usr/local/bin

COPY --from=builder /usr/local/bin/app .


ENV MIN_WORKERS = 1
ENV MAX_WORKERS = 12
ENV BUSY_THRESHOLD = 50
ENV REMINDS_QUEUE_SIZE = 200
ENV SCALE_UP_THRESHOLD = 30
ENV SCALE_DOWN_THRESHOLD = 10
ENV CHECK_REMINDS_INTERVAL = 1
ENV WORKERS_CHECK_INTERVAL = 15

ENTRYPOINT ["./app"]