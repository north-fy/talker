FROM golang:1.26.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./п
RUN go mod download
COPY . .
RUN go build -o main ./cmd/server/
EXPOSE 50052
CMD ["/app/main"]