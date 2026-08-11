FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /server /app/server

EXPOSE 50053

ENTRYPOINT ["/app/server"]
