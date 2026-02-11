# Nookverse API 文档

## 概述

Nookverse API 是一个完整的家庭物品管理系统API，提供了物品管理、房间管理、提醒管理和统计分析等功能。

## API 规范

API 遵循 OpenAPI 3.0.3 规范，完整文档位于 [`openapi.json`](openapi.json) 文件中。

## 主要功能模块

### 1. 物品管理 (Items)
- **创建物品**: `POST /api/v1/items`
- **获取物品列表**: `GET /api/v1/items`
- **搜索物品**: `GET /api/v1/items/search`
- **获取物品详情**: `GET /api/v1/items/{itemId}`
- **更新物品**: `PUT /api/v1/items/{itemId}`
- **删除物品**: `DELETE /api/v1/items/{itemId}`

### 2. 物品层级管理
- **移动物品**: `POST /api/v1/items/{itemId}/move`
- **获取容器内容**: `GET /api/v1/items/container/{containerId}/contents`

### 3. 提醒管理 (Reminders)
- **创建提醒**: `POST /api/v1/items/{itemId}/reminders`
- **获取即将到来的提醒**: `GET /api/v1/items/reminders/upcoming`

### 4. 房间管理 (Rooms)
- **获取房间内物品**: `GET /api/v1/rooms/{roomId}/items`

### 5. 统计分析 (Statistics)
- **获取物品统计信息**: `GET /api/v1/items/statistics`

## 认证机制

部分接口需要 JWT Token 认证，在请求头中添加：
```
Authorization: Bearer <your-jwt-token>
```

## 错误响应格式

所有错误响应遵循统一格式：
```json
{
  "error": "错误描述信息"
}
```

常见HTTP状态码：
- `200`: 请求成功
- `201`: 创建成功
- `400`: 请求参数错误
- `401`: 未授权访问
- `404`: 资源不存在
- `500`: 服务器内部错误

## 数据模型

### 物品模型 (Item)
```json
{
  "id": "物品唯一标识",
  "name": "物品名称",
  "description": "物品描述",
  "category": "分类信息",
  "room": "房间信息",
  "quantity": "数量",
  "status": "状态(active/archived/discarded/borrowed)",
  "expire_date": "过期时间",
  "purchase_date": "购买时间",
  "price": "价格",
  "brand": "品牌",
  "model": "型号",
  "position": "位置信息",
  "labels": ["标签列表"],
  "media_files": ["媒体文件列表"],
  "reminders": ["提醒列表"]
}
```

### 提醒模型 (Reminder)
```json
{
  "id": "提醒ID",
  "reminder_type": "提醒类型(expire/maintenance/warranty/custom)",
  "trigger_time": "触发时间",
  "message": "提醒消息",
  "status": "状态(pending/sent/completed/cancelled)",
  "notify_channels": ["通知渠道(app/email/sms/voice)"]
}
```

## 使用示例

### 1. 创建物品
```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "description": "工作用笔记本电脑",
    "category_id": "category-uuid",
    "room_id": "room-uuid",
    "quantity": 1,
    "status": "active",
    "price": 12999.00,
    "brand": "Apple",
    "model": "MacBook Pro 14-inch",
    "labels": ["electronics", "work"]
  }'
```

### 2. 搜索物品
```bash
curl -X GET "http://localhost:8080/api/v1/items/search?q=MacBook&page=1&page_size=10"
```

### 3. 创建提醒
```bash
curl -X POST http://localhost:8080/api/v1/items/item-uuid/reminders \
  -H "Content-Type: application/json" \
  -d '{
    "reminder_type": "warranty",
    "trigger_time": "2025-02-10T10:00:00Z",
    "message": "MacBook保修期即将到期",
    "notify_channels": ["app", "email"]
  }'
```

## 开发工具

推荐使用以下工具来查看和测试API：

1. **Swagger UI**: 可以直接在浏览器中查看和测试API
2. **Postman**: 导入openapi.json文件进行测试
3. **Insomnia**: 支持OpenAPI规范的API测试工具

## 注意事项

1. 所有时间字段使用ISO 8601格式 (`YYYY-MM-DDTHH:mm:ssZ`)
2. ID字段使用UUID格式
3. 分页查询默认每页10条记录
4. 搜索功能支持模糊匹配
5. 部分接口需要认证，请确保携带有效的JWT Token

## 版本历史

- **v1.0.0**: 初始版本，包含基本的物品管理功能

---
*本文档基于OpenAPI 3.0.3规范生成*