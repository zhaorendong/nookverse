#!/bin/bash
# OpenAPI 文档验证脚本

echo "正在验证 OpenAPI 文档..."

# 检查文件是否存在
if [ ! -f "openapi.json" ]; then
    echo "❌ 错误: openapi.json 文件不存在"
    exit 1
fi

# 使用 swagger-cli 验证 (如果已安装)
if command -v swagger-cli &> /dev/null; then
    echo "🔍 使用 swagger-cli 验证..."
    swagger-cli validate openapi.json
    if [ $? -eq 0 ]; then
        echo "✅ OpenAPI 文档验证通过"
    else
        echo "❌ OpenAPI 文档验证失败"
        exit 1
    fi
else
    echo "⚠️  未找到 swagger-cli，跳过详细验证"
fi

# 基本的 JSON 格式验证
echo "🔍 验证 JSON 格式..."
python3 -m json.tool openapi.json > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ JSON 格式正确"
else
    echo "❌ JSON 格式错误"
    exit 1
fi

# 检查必需字段
echo "🔍 检查必需字段..."
REQUIRED_FIELDS=("openapi" "info" "paths" "components")
for field in "${REQUIRED_FIELDS[@]}"; do
    if python3 -c "import json; data = json.load(open('openapi.json')); exit(0 if '$field' in data else 1)" 2>/dev/null; then
        echo "✅ 包含必需字段: $field"
    else
        echo "❌ 缺少必需字段: $field"
        exit 1
    fi
done

echo "🎉 OpenAPI 文档验证完成！"
echo ""
echo "文档信息:"
python3 -c "
import json
with open('openapi.json') as f:
    data = json.load(f)
    print(f'标题: {data[\"info\"][\"title\"]}')
    print(f'版本: {data[\"info\"][\"version\"]}')
    print(f'路径数量: {len(data[\"paths\"])}')
    print(f'组件模式数量: {len(data[\"components\"][\"schemas\"])}')
"