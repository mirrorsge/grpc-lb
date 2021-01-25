.PHONY: proto
proto:
	protoc -I ./examples/proto \
	--go_out ./examples/proto \
	--go_opt paths=source_relative \
	--go-grpc_out ./examples/proto \
	--go-grpc_opt paths=source_relative \
	./examples/proto/*.proto