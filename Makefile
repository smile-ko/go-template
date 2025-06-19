APP_NAME = server

PROTO_DIR = proto
OUT_DIR = proto/pb

.PHONY: proto
proto:
	protoc --proto_path=$(PROTO_DIR) \
	       --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(shell find $(PROTO_DIR) -name '*.proto')

proto-clean:
	rm -rf $(OUT_DIR)/*.pb.go $(OUT_DIR)/*.grpc.pb.go

run:
	go run ./cmd/$(APP_NAME)/
