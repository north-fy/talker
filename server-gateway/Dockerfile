FROM golang:1.26.1-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /gateway ./cmd/app

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /gateway /app/gateway

EXPOSE 8080

ENTRYPOINT ["/app/gateway"]