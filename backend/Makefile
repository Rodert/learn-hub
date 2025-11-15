.PHONY: run build test clean deps

# 安装依赖
deps:
	go mod download
	go mod tidy

# 运行服务
run:
	go run main.go

# 构建
build:
	go build -o bin/server main.go

# 测试
test:
	go test ./...

# 清理
clean:
	rm -rf bin/

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

