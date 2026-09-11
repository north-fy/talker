# Build context MUST be the repository root (run: docker build -f services/chat/Dockerfile -t talker/chat:dev .)
# Uses go.work to resolve the shared github.com/north-fy/talker module (pkg/protos).
FROM golang:1.26.1-bookworm AS builder
WORKDIR /src
COPY . .
WORKDIR /src/services/chat
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o /out/main ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /out/main /main
EXPOSE 50053
ENTRYPOINT ["/main"]
