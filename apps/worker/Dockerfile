FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY apps/worker/go.mod ./go.mod
RUN go mod download

COPY apps/worker/. ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/worker ./cmd/worker

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/worker /usr/local/bin/worker
CMD ["worker"]
