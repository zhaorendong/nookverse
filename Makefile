# NookVerse Makefile

.PHONY: help build run test clean migrate

# 默认目标
help:
	@echo "NookVerse 项目 Makefile"
	@echo ""
	@echo "可用命令:"
	@echo "  build     - 构建项目"
	@echo "  run       - 运行项目"
	@echo "  test      - 运行测试"
	@echo "  clean     - 清理构建文件"
	@echo "  migrate   - 执行数据库迁移（自动创建数据库并初始化）"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  docker-run   - 运行Docker容器"

# 构建项目
build:
	go build -o bin/nookverse cmd/server/main.go

# 运行项目
run: build
	./bin/nookverse

# 运行测试
test:
	go test -v ./...

# 清理构建文件
clean:
	rm -rf bin/
	go clean

# 数据库迁移（自动处理数据库创建和初始化）
migrate:
	go run db/migrate.go

# 构建Docker镜像
docker-build:
	docker build -t nookverse .

# 运行Docker容器
docker-run:
	docker run -p 8080:8080 --name nookverse-app nookverse

# 启动Docker Compose服务
docker-compose-up:
	docker-compose up -d

# 停止Docker Compose服务
docker-compose-down:
	docker-compose down

# 开发环境运行
dev:
	go run cmd/server/main.go

# 安装依赖
deps:
	go mod tidy
	go mod download

# 代码格式化
fmt:
	go fmt ./...

# 代码检查
vet:
	go vet ./...