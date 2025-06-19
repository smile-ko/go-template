APP_NAME = server


.PHONY: proto
proto-user:
	protoc --proto_path=docs/proto/user/v1 \
	       --go_out=docs/proto/user/v1/gen --go_opt=paths=source_relative \
	       --go-grpc_out=docs/proto/user/v1/gen --go-grpc_opt=paths=source_relative \
	      	docs/proto/user/v1/*.proto


run:
	go run ./cmd/$(APP_NAME)/
