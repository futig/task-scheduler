FROM golang:1.23.6-bookworm AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/app ./cmd/taskscheduler

FROM debian:bookworm-slim

WORKDIR /usr/local/bin

COPY --from=builder /usr/local/bin/app .

ENTRYPOINT ["./app"]