FROM golang:1.26.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/server/
EXPOSE 50051
CMD ["/app/main"]