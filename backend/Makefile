.PHONY: help build run test clean migrate swagger docker-build docker-run

# 变量定义
BINARY_NAME=learn-hub
GO=go
GOFLAGS=-v
SWAG=swag

help:
	@echo "Learn Hub Backend - Available Commands"
	@echo "========================================"
	@echo "make build          - 构建二进制文件"
	@echo "make run            - 运行应用"
	@echo "make dev            - 开发模式运行（支持热重载）"
	@echo "make test           - 运行单元测试"
	@echo "make test-coverage  - 运行测试并生成覆盖率报告"
	@echo "make migrate        - 执行数据库迁移"
	@echo "make swagger        - 生成 Swagger 文档"
	@echo "make clean          - 清理构建文件"
	@echo "make docker-build   - 构建 Docker 镜像"
	@echo "make docker-run     - 运行 Docker 容器"
	@echo "make lint           - 代码检查"
	@echo "make fmt            - 代码格式化"

# 构建二进制文件
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o bin/$(BINARY_NAME) ./cmd/main.go
	@echo "Build complete: bin/$(BINARY_NAME)"

# 运行应用
run: build
	@echo "Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME)

# 开发模式运行（需要 air 工具）
dev:
	@echo "Running in development mode..."
	@command -v air >/dev/null 2>&1 || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

# 运行单元测试
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 数据库迁移
migrate:
	@echo "Running database migrations..."
	$(GO) run ./cmd/migrate/main.go
	@echo "Migration complete"

# 生成 Swagger 文档
swagger:
	@echo "Generating Swagger documentation..."
	@command -v swag >/dev/null 2>&1 || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/main.go -o docs
	@echo "Swagger documentation generated in docs/"
	@echo "Access Swagger UI at: http://localhost:8080/swagger/index.html"

# 清理构建文件
clean:
	@echo "Cleaning up..."
	$(GO) clean
	rm -f bin/$(BINARY_NAME)
	rm -rf docs/
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# 代码检查
lint:
	@echo "Running linter..."
	@command -v golangci-lint >/dev/null 2>&1 || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# 代码格式化
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	goimports -w .

# 构建 Docker 镜像
docker-build:
	@echo "Building Docker image..."
	docker build -t learn-hub-backend:latest .
	@echo "Docker image built: learn-hub-backend:latest"

# 运行 Docker 容器
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env learn-hub-backend:latest
	@echo "Container started"

# 安装依赖
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies installed"

# 生成 API 文档和构建
all: deps swagger build
	@echo "All tasks completed"
