#!/bin/bash

# House API 演示脚本
# 展示如何使用House CRUD接口

echo "🏠 NookVerse House API 演示"
echo "============================="

# 服务器地址
BASE_URL="http://localhost:8080"

echo "1. 创建房屋"
echo "-----------"
curl -X POST "${BASE_URL}/api/v1/houses" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "我的温馨小家",
    "address": "北京市朝阳区幸福小区1号楼101室",
    "description": "一套两居室的温馨住宅",
    "area": 85.5,
    "floor_count": 2,
    "metadata": {
      "year_built": 2018,
      "has_elevator": true,
      "orientation": "south"
    }
  }' | jq '.'

echo -e "\n2. 获取房屋列表"
echo "---------------"
curl -X GET "${BASE_URL}/api/v1/houses?page=1&page_size=10" | jq '.'

echo -e "\n3. 搜索房屋"
echo "-----------"
curl -X GET "${BASE_URL}/api/v1/houses/search?q=北京" | jq '.'

echo -e "\n4. 创建房间（需要替换HOUSE_ID）"
echo "-----------------------------"
# 这里需要先获取房屋ID
HOUSE_ID=$(curl -s -X GET "${BASE_URL}/api/v1/houses" | jq -r '.data[0].id')

if [ "$HOUSE_ID" != "null" ]; then
  echo "使用房屋ID: $HOUSE_ID"
  
  curl -X POST "${BASE_URL}/api/v1/houses/${HOUSE_ID}/rooms" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "主卧室",
      "room_type": "bedroom",
      "floor_number": 1,
      "area": 20.0,
      "description": "朝南的主卧室，采光良好",
      "position_data": {
        "x": 0,
        "y": 0,
        "z": 0,
        "width": 4,
        "length": 5
      }
    }' | jq '.'
  
  echo -e "\n5. 创建客厅"
  curl -X POST "${BASE_URL}/api/v1/houses/${HOUSE_ID}/rooms" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "客厅",
      "room_type": "living_room",
      "floor_number": 1,
      "area": 25.0,
      "description": "宽敞明亮的客厅"
    }' | jq '.'
  
  echo -e "\n6. 获取房屋详情（包含房间）"
  curl -X GET "${BASE_URL}/api/v1/houses/${HOUSE_ID}" | jq '.'
  
  echo -e "\n7. 获取房屋内房间列表"
  curl -X GET "${BASE_URL}/api/v1/houses/${HOUSE_ID}/rooms" | jq '.'
  
  echo -e "\n8. 获取房屋统计信息"
  curl -X GET "${BASE_URL}/api/v1/houses/statistics" | jq '.'
else
  echo "未找到房屋，跳过房间创建步骤"
fi

echo -e "\n演示完成！ 🎉"