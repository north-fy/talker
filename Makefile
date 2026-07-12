.PHONY: proto
proto:
	protoc --go_out=$(OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GRPC_OUT) --go-grpc_opt=paths=source_relative \
		$(DIR_PROTO)

.PHONY: proto-deps
proto-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpcD@v1.2

.PHONY: gen
gen:
	protoc \
      --go_out=./pkg/protos/chat --go_opt=paths=source_relative \
      --go-grpc_out=./pkg/protos/chat --go-grpc_opt=paths=source_relative \
      --proto_path=./pkg/protos \
      ./pkg/protos/chat.proto