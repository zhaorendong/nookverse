# House API 使用指南

## 概述

House API 提供了完整的房屋和房间管理功能，允许用户创建、查询、更新和删除房屋及房间信息。

## 主要功能

### 1. 房屋管理
- 创建房屋
- 查询房屋列表（支持分页和过滤）
- 搜索房屋
- 获取房屋详情
- 更新房屋信息
- 删除房屋

### 2. 房间管理
- 在房屋内创建房间
- 查询房屋内房间列表
- 获取房间详情
- 更新房间信息
- 删除房间

### 3. 统计功能
- 获取房屋和房间统计信息

## API 端点详解

### 房屋相关接口

#### 创建房屋
```
POST /api/v1/houses
```

**请求体示例：**
```json
{
  "name": "我的家",
  "address": "北京市朝阳区某某街道123号",
  "description": "这是一套温馨的三居室",
  "area": 120.5,
  "floor_count": 2,
  "metadata": {
    "year_built": 2020,
    "has_garden": true,
    "parking_spaces": 2
  }
}
```

#### 获取房屋列表
```
GET /api/v1/houses?page=1&page_size=20&name=我家&min_area=100
```

**支持的查询参数：**
- `name`: 房屋名称模糊搜索
- `address`: 地址模糊搜索
- `min_area`: 最小面积
- `max_area`: 最大面积
- `min_floors`: 最小楼层数
- `max_floors`: 最大楼层数
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）
- `order_by`: 排序字段（created_at, name, area, floor_count）

#### 搜索房屋
```
GET /api/v1/houses/search?q=北京&min_area=80
```

#### 获取房屋详情
```
GET /api/v1/houses/{houseId}
```

#### 更新房屋
```
PUT /api/v1/houses/{houseId}
```

#### 删除房屋
```
DELETE /api/v1/houses/{houseId}
```
> 注意：只有当房屋内没有房间时才能删除

### 房间相关接口

#### 创建房间
```
POST /api/v1/houses/{houseId}/rooms
```

**请求体示例：**
```json
{
  "name": "主卧室",
  "room_type": "bedroom",
  "floor_number": 1,
  "area": 25.0,
  "description": "朝南的主卧室",
  "position_data": {
    "x": 10,
    "y": 20,
    "z": 0,
    "width": 5,
    "length": 4
  }
}
```

**支持的房间类型：**
- `bedroom`: 卧室
- `living_room`: 客厅
- `kitchen`: 厨房
- `bathroom`: 浴室
- `study`: 书房
- `storage`: 储藏室
- `other`: 其他

#### 获取房屋内房间
```
GET /api/v1/houses/{houseId}/rooms
```

#### 获取房间详情
```
GET /api/v1/rooms/{roomId}
```

#### 更新房间
```
PUT /api/v1/rooms/{roomId}
```

#### 删除房间
```
DELETE /api/v1/rooms/{roomId}
```
> 注意：只有当房间内没有物品时才能删除

### 统计接口

#### 获取房屋统计
```
GET /api/v1/houses/statistics
```

**响应示例：**
```json
{
  "data": {
    "total_houses": 5,
    "total_rooms": 20,
    "average_area": 115.5,
    "by_floor_count": {
      "1": 2,
      "2": 2,
      "3": 1
    },
    "by_room_type": {
      "bedroom": 8,
      "living_room": 5,
      "kitchen": 4,
      "bathroom": 3
    }
  }
}
```

## 错误处理

API 使用标准的 HTTP 状态码：

- `200`: 请求成功
- `201`: 创建成功
- `400`: 请求参数错误
- `404`: 资源不存在
- `500`: 服务器内部错误

**错误响应格式：**
```json
{
  "error": "错误描述信息"
}
```

## 使用示例

### JavaScript 客户端示例

```javascript
// 创建房屋
async function createHouse(houseData) {
  const response = await fetch('/api/v1/houses', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(houseData)
  });
  
  if (!response.ok) {
    throw new Error(await response.text());
  }
  
  return response.json();
}

// 获取房屋列表
async function getHouses(filters = {}) {
  const params = new URLSearchParams(filters);
  const response = await fetch(`/api/v1/houses?${params}`);
  
  if (!response.ok) {
    throw new Error(await response.text());
  }
  
  return response.json();
}

// 创建房间
async function createRoom(houseId, roomData) {
  const response = await fetch(`/api/v1/houses/${houseId}/rooms`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(roomData)
  });
  
  if (!response.ok) {
    throw new Error(await response.text());
  }
  
  return response.json();
}
```

### Python 客户端示例

```python
import requests

class HouseClient:
    def __init__(self, base_url='http://localhost:8080'):
        self.base_url = base_url
    
    def create_house(self, house_data):
        response = requests.post(
            f'{self.base_url}/api/v1/houses',
            json=house_data
        )
        response.raise_for_status()
        return response.json()
    
    def get_houses(self, **filters):
        response = requests.get(
            f'{self.base_url}/api/v1/houses',
            params=filters
        )
        response.raise_for_status()
        return response.json()
    
    def create_room(self, house_id, room_data):
        response = requests.post(
            f'{self.base_url}/api/v1/houses/{house_id}/rooms',
            json=room_data
        )
        response.raise_for_status()
        return response.json()
```

## 注意事项

1. **数据完整性**：删除房屋前必须先删除其所有房间
2. **数据验证**：所有必填字段都需要提供有效值
3. **UUID格式**：ID字段使用标准UUID格式
4. **分页机制**：大列表查询建议使用分页参数
5. **搜索功能**：支持模糊搜索房屋名称、地址和描述

## 性能优化建议

1. 合理使用分页参数避免一次性加载大量数据
2. 利用过滤参数减少不必要的数据传输
3. 对于频繁查询的数据考虑客户端缓存
4. 批量操作时注意数据库事务的一致性

这个API设计遵循RESTful原则，提供了完整的房屋和房间管理功能，可以满足大多数家庭物品管理系统的房屋结构需求。