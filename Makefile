# Allow user to have custom version of go.
GOBINARY ?= go
GOOS ?= linux

LDFLAGS += -w -s

# Install build dependencies.
tools-update:
	$(GOBINARY) get -u github.com/golang/dep/cmd/dep/...
	$(GOBINARY) get -u google.golang.org/grpc
	$(GOBINARY) get -u github.com/golang/protobuf/proto
	$(GOBINARY) get -u github.com/golang/protobuf/protoc-gen-go
	$(GOBINARY) get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	$(GOBINARY) get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
.PHONY: tools-update

deps-update:
	$(ENVVAR) GOOS=$(GOOS) dep ensure
.PHONY: deps-update

gen:
	# Generate Go code for gRPC server and client
	protoc \
		-Iapi/protobuf-spec \
		--go_out=plugins=grpc:pkg/colorer \
		api/protobuf-spec/colorer.proto

	protoc \
		-Iapi/protobuf-spec \
		--go_out=plugins=grpc:pkg/aggregator \
		api/protobuf-spec/aggregator.proto

	# Add the following options to generate
	# - Go code for JSON gRPC gateway
	# - Generate Swagger spec
	#-I$(GOPATH)/src
	#-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
	#--grpc-gateway_out=logtostderr=true:pkg
	#--swagger_out=logtostderr=true:pkg
.PHONY: gen

build:
	$(ENVVAR) GOOS=$(GOOS) $(GOBINARY) build -i -v -ldflags '$(LDFLAGS)' -o target/server cmd/server.go
.PHONY: build

all: tools-update deps-update gen build
.PHONY: all
