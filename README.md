# NookVerse - 智能家居物品管理系统

NookVerse是一个将房子里所有物品，大到家具电器，小到柜子里的每一个物品，全部数字化存储并进行管理的系统，让每个家庭都能实现"一屋之内，自有乾坤"的智慧生活，Nookverse 旨在通过数字化+智能化手段，彻底解决"家中物品找不到、记不清、管不好"的痛点，让每个家庭都能实现"一屋之内，自有乾坤"的智慧生活。

## 🌟 核心功能

### 📦 物品全生命周期管理
- 物品登记、编辑、删除
- 分类管理、标签系统
- 多媒体文件关联
- 状态追踪（在用、闲置、丢弃、出借）

### 🏠 空间层级管理
- 多房屋/地址管理
- 房间划分和3D空间映射
- 无限层级嵌套（房间→家具→抽屉→盒子→物品）

### 🔍 智能检索
- 快速搜索和条件筛选
- 语义搜索和位置检索
- 高级组合查询

### ⏰ 智能提醒
- 过期提醒（食品、药品等）
- 保修期提醒
- 维护提醒
- 自定义时间节点提醒

### 📱 多端协同
- Web端管理后台
- 移动APP（iOS/Android）
- 小程序支持
- 智能音箱语音交互

## 🏗️ 项目架构

```
nookverse/
├── cmd/
│   └── server/
│       └── main.go              # 应用程序入口
├── internal/
│   ├── config/                  # 配置管理
│   ├── database/                # 数据库连接
│   ├── models/                  # 数据模型
│   ├── routers/                 # 路由定义
│   ├── services/                # 业务逻辑层
│   └── utils/                   # 工具函数
├── pkg/
│   └── api/
│       └── v1/
│           ├── handlers/        # HTTP处理器
│           └── dto/             # 数据传输对象
├── db/                          # 数据库脚本
│   ├── init.sql                 # 初始化SQL
│   └── migrate.go               # 迁移工具
├── Dockerfile                   # Docker配置
├── Makefile                     # 构建脚本
├── go.mod                       # 依赖管理
└── README.md                    # 项目文档
```

## 🚀 快速开始

### 环境要求
- Go 1.25+
- PostgreSQL 12+
- Docker (可选)

### 1. 克隆项目
```bash
git clone <repository-url>
cd nookverse
```

### 2. 配置环境
```bash
# 复制配置文件模板
cp config.example.json config.json

# 编辑配置文件，设置数据库连接等参数
vim config.json
```

### 3. 初始化数据库
```bash
# 创建数据库
createdb nookverse

# 执行初始化脚本
psql -d nookverse -f db/init.sql
```

或者使用Makefile：
```bash
make migrate
```

### 4. 运行应用
```bash
# 安装依赖
go mod tidy

# 运行开发环境
make dev

# 或者直接运行
go run cmd/server/main.go
```

### 5. 访问应用
打开浏览器访问：http://localhost:8080

健康检查端点：http://localhost:8080/health

## 📊 API 接口文档

### 物品管理
```
POST   /api/v1/items              # 创建物品
GET    /api/v1/items              # 获取物品列表
GET    /api/v1/items/search       # 搜索物品
GET    /api/v1/items/{id}         # 获取物品详情
PUT    /api/v1/items/{id}         # 更新物品
DELETE /api/v1/items/{id}         # 删除物品

POST   /api/v1/items/{id}/move    # 移动物品到容器
GET    /api/v1/items/{containerId}/contents  # 获取容器内物品

POST   /api/v1/items/{id}/reminders  # 创建提醒
GET    /api/v1/items/reminders/upcoming  # 获取即将到来的提醒

GET    /api/v1/items/statistics   # 获取物品统计信息
```

### 房间管理
```
GET    /api/v1/rooms/{roomId}/items  # 获取房间内物品
```

### 用户管理
```
POST   /api/v1/users/register     # 用户注册
POST   /api/v1/users/login        # 用户登录
GET    /api/v1/users/{id}         # 获取用户信息
PUT    /api/v1/users/{id}         # 更新用户信息
DELETE /api/v1/users/{id}         # 删除用户
```

## 🐳 Docker 部署

### 构建镜像
```bash
make docker-build
# 或
docker build -t nookverse .
```

### 运行容器
```bash
make docker-run
# 或
docker run -p 8080:8080 --name nookverse-app nookverse
```

### Docker Compose (推荐)
```yaml
version: '3.8'

services:
  # 应用服务
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - CONFIG_PATH=/app/config.json
    volumes:
      - ./config.json:/app/config.json:ro
      - ./uploads:/app/uploads
    depends_on:
      - db
      - redis
    restart: unless-stopped
    networks:
      - nookverse-network

  # PostgreSQL数据库
  db:
    image: postgres:13-alpine
    environment:
      POSTGRES_DB: nookverse
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-password}
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    ports:
      - "5432:5432"
    restart: unless-stopped
    networks:
      - nookverse-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis缓存服务
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
      - ./redis.conf:/usr/local/etc/redis/redis.conf:ro
    command: redis-server /usr/local/etc/redis/redis.conf
    restart: unless-stopped
    networks:
      - nookverse-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3

  # Redis Commander (可选的Web管理界面)
  redis-commander:
    image: rediscommander/redis-commander:latest
    environment:
      - REDIS_HOSTS=local:redis:6379
    ports:
      - "8081:8081"
    depends_on:
      - redis
    restart: unless-stopped
    networks:
      - nookverse-network

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local

networks:
  nookverse-network:
    driver: bridge
```

## 🛠️ 开发工具

### Makefile 命令
```bash
make help      # 查看所有可用命令
make build     # 构建项目
make run       # 运行项目
make test      # 运行测试
make clean     # 清理构建文件
make migrate   # 执行数据库迁移
make init-db   # 初始化数据库
make fmt       # 代码格式化
make vet       # 代码检查
```

## 🔧 环境变量配置

项目支持通过环境变量来配置关键参数。可以使用 `.env` 文件来管理环境变量。

### 环境变量文件
```bash
# 复制环境变量模板
cp .env.example .env

# 编辑环境变量
code .env  # 或使用 vim .env
```

### 主要环境变量说明
| 变量名 | 说明 | 默认值 | 重要性 |
|--------|------|--------|--------|
| `POSTGRES_PASSWORD` | PostgreSQL数据库密码 | password | ⭐⭐⭐ |
| `CONFIG_PATH` | 配置文件路径 | /app/config.json | ⭐⭐ |
| `JWT_SECRET` | JWT密钥 | your-jwt-secret-key | ⭐⭐⭐ |
| `REDIS_PASSWORD` | Redis密码 | 空 | ⭐⭐ |
| `SERVER_PORT` | 服务器端口 | 8080 | ⭐ |

> ⚠️ **安全提醒**：生产环境务必修改默认密码和密钥！

## 🔧 配置说明

### config.json 配置项
```json
{
  "server": {
    "port": 8080
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "password",
    "name": "nookverse"
  },
  "jwt": {
    "secret": "your-jwt-secret-key",
    "expire": 24
  },
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "upload": {
    "path": "./uploads",
    "max_size": 10485760,
    "allowed_types": ["image/jpeg", "image/png", "image/gif"]
  }
}
```

## 🧪 测试

```bash
# 运行所有测试
make test

# 运行特定包测试
go test ./internal/services/...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 代码规范
- 遵循 Go 官方编码规范
- 使用 `make fmt` 格式化代码
- 添加必要的单元测试
- 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

*Made with ❤️ for smarter home management*