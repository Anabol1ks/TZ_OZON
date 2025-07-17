generate_proto:
	protoc --go_out=. --go-grpc_out=. --proto_path=proto proto/exchange.proto
