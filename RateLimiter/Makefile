# Makefile para Rate Limiter

.PHONY: help build run test clean docker-up docker-down docker-logs

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Compila a aplicação
	go build -o bin/server ./cmd/server

run: ## Executa a aplicação localmente
	go run ./cmd/server/main.go

test: ## Executa todos os testes
	go test ./... -v

test-cover: ## Executa testes com cobertura
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Cobertura gerada em coverage.html"

clean: ## Limpa arquivos temporários e binários
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -testcache

docker-up: ## Inicia os containers Docker
	docker-compose up -d

docker-build: ## Builda e inicia os containers
	docker-compose up -d --build

docker-down: ## Para os containers Docker
	docker-compose down

docker-logs: ## Mostra os logs dos containers
	docker-compose logs -f

docker-clean: ## Para containers e remove volumes
	docker-compose down -v

lint: ## Executa linter
	golangci-lint run ./...

deps: ## Baixa as dependências
	go mod download
	go mod tidy

integration-test: ## Executa testes de integração (Linux/Mac)
	bash test-integration.sh

integration-test-windows: ## Executa testes de integração (Windows)
	powershell -ExecutionPolicy Bypass -File test-integration.ps1
