.PHONY: gen
gen:
	protoc -I . \
        --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
        pkg/protos/message.proto

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: rebuild
rebuild:
	docker-compose up -d --build

.PHONY: logs
logs:
	docker-compose logs

.PHONY: migrations
migrations:
	migrate -path=/migrations -database="postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" up
