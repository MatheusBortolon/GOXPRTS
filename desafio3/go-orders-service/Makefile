GO ?= go
BINARY_NAME = go-orders-service

.PHONY: all build run test clean migrate

all: build

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/rest/main.go
	$(GO) build -o $(BINARY_NAME)-grpc ./cmd/grpc/main.go
	$(GO) build -o $(BINARY_NAME)-graphql ./cmd/graphql/main.go

run: build
	./$(BINARY_NAME)

test:
	$(GO) test ./...

clean:
	$(GO) clean
	rm -f $(BINARY_NAME) $(BINARY_NAME)-grpc $(BINARY_NAME)-graphql

migrate:
	./scripts/migrate.sh