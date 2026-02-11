# Docker部署说明

## 快速开始

### 1. 准备配置文件
```bash
# 复制配置文件模板
cp config.example.json config.json

# 编辑配置文件，设置数据库密码等参数
vim config.json
```

### 2. 启动所有服务
```bash
# 使用Docker Compose启动所有服务
make docker-compose-up

# 或者直接使用docker-compose命令
docker-compose up -d
```

### 3. 验证服务状态
```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f app
docker-compose logs -f db
docker-compose logs -f redis
```

## 服务访问地址

- **应用服务**: http://localhost:8080
- **健康检查**: http://localhost:8080/health
- **Redis管理界面**: http://localhost:8081
- **PostgreSQL**: localhost:5432

## 常用命令

```bash
# 停止所有服务
make docker-compose-down
# 或
docker-compose down

# 重新构建并启动
docker-compose up -d --build

# 查看容器日志
docker-compose logs -f

# 进入应用容器
docker-compose exec app sh

# 进入数据库容器
docker-compose exec db psql -U postgres -d nookverse

# 进入Redis容器
docker-compose exec redis redis-cli
```

## 环境变量

可以在 `.env` 文件中设置环境变量：

```bash
# 数据库密码
POSTGRES_PASSWORD=your_secure_password

# Redis密码（如果启用）
REDIS_PASSWORD=your_redis_password
```

## 数据持久化

数据卷会自动创建并持久化：

- `postgres_data`: PostgreSQL数据
- `redis_data`: Redis数据

## 故障排除

### 1. 端口冲突
如果端口被占用，可以修改docker-compose.yml中的端口映射：

```yaml
ports:
  - "8080:8080"  # 改为你想要的端口
```

### 2. 权限问题
确保配置文件和数据目录有正确的权限：

```bash
chmod 644 config.json
chmod 755 uploads/
```

### 3. 数据库连接失败
检查数据库配置是否正确，确认PostgreSQL服务正常运行。

### 4. Redis连接失败
检查Redis配置和服务状态，确认端口6379未被占用。