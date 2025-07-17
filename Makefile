generate_proto:
	protoc --go_out=. --go-grpc_out=. --proto_path=proto proto/exchange.proto

docker-up:
	docker-compose -f docker-compose.yml up -d

run:
	go run cmd/main.go

test:
	go test internal/handler/exchange_test.go