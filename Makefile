PROJECT_PATH = /home/kadr/go/src/github.com/cost_control
PROTO_FILES_PATH = $(PROJECT_PATH)/internal/handlers/rpc/proto
SRC_FILES_PATH = $(PROJECT_PATH)/internal/handlers/rpc/src
gen:
	protoc -I=$(PROTO_FILES_PATH) --go_out=$(SRC_FILES_PATH) --experimental_allow_proto3_optional --proto_path=/usr/include $(PROTO_FILES_PATH)/*.proto
gen-full:
	protoc -I=$(PROTO_FILES_PATH) --go_out=$(SRC_FILES_PATH) --go-grpc_out=$(SRC_FILES_PATH) --experimental_allow_proto3_optional --proto_path=/usr/include $(PROTO_FILES_PATH)/*.proto