PROTO_SRC := ./proto
PROTO_OUT := ./pkg/proto
GO_MAIN := ./cmd/UnitTracker/main.go
GO_OUT := ./bin


# Tools
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc

# Ensuring tools are installed
.PHONY: tools
tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Compiling protobuf files
.PHONY: proto
proto: tools
	mkdir -p $(PROTO_OUT)
	$(PROTOC) --proto_path=$(PROTO_SRC) --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative $(PROTO_SRC)/*.proto

# Building the application
.PHONY: build
build: proto
	mkdir -p $(GO_OUT)
	go build -o $(GO_OUT) $(GO_MAIN)

# Running the application
.PHONY: run
run: build
	$(GO_OUT)/main


# Cleaning up generated files and build artifacts
.PHONY: clean
clean:
	rm -rf $(PROTO_OUT)/*
	rm -rf $(GO_OUT)

