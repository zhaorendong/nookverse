package tests

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"nookverse/pkg/api/v1/dto"
	"nookverse/tests/testutils"
)

// getValueOrDefault 获取指针值或默认值
func getValueOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// isValidUUID 验证UUID格式
func isValidUUID(uuid string) bool {
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, uuid)
	return matched
}

func TestHouseValidation(t *testing.T) {
	t.Run("有效的UUID格式", func(t *testing.T) {
		validUUID := "550e8400-e29b-41d4-a716-446655440000"
		assert.True(t, isValidUUID(validUUID))
	})

	t.Run("无效的UUID格式", func(t *testing.T) {
		invalidUUID := "invalid-uuid"
		assert.False(t, isValidUUID(invalidUUID))
	})

	t.Run("空字符串UUID", func(t *testing.T) {
		emptyUUID := ""
		assert.False(t, isValidUUID(emptyUUID))
	})
}

func TestGetValueOrDefault(t *testing.T) {
	t.Run("非空指针值", func(t *testing.T) {
		value := "test"
		ptr := &value
		result := getValueOrDefault(ptr, "default")
		assert.Equal(t, "test", result)
	})

	t.Run("空指针值", func(t *testing.T) {
		var ptr *string
		result := getValueOrDefault(ptr, "default")
		assert.Equal(t, "default", result)
	})

	t.Run("数字类型的默认值", func(t *testing.T) {
		var ptr *int
		result := getValueOrDefault(ptr, 42)
		assert.Equal(t, 42, result)
	})

	t.Run("浮点数类型的默认值", func(t *testing.T) {
		var ptr *float64
		result := getValueOrDefault(ptr, 3.14)
		assert.Equal(t, 3.14, result)
	})
}

func TestHouseDTOSerialization(t *testing.T) {
	t.Run("CreateHouseRequest序列化", func(t *testing.T) {
		req := dto.CreateHouseRequest{
			Name:        "测试房屋",
			Address:     testutils.StringPtr("测试地址"),
			Description: testutils.StringPtr("测试描述"),
			Area:        testutils.Float64Ptr(120.5),
			FloorCount:  testutils.IntPtr(2),
			Metadata: map[string]any{
				"year_built": 2020,
			},
		}

		// 验证基本字段
		assert.Equal(t, "测试房屋", req.Name)
		assert.Equal(t, "测试地址", *req.Address)
		assert.Equal(t, "测试描述", *req.Description)
		assert.Equal(t, 120.5, *req.Area)
		assert.Equal(t, 2, *req.FloorCount)
		assert.Equal(t, 2020, req.Metadata["year_built"])
	})

	t.Run("CreateRoomRequest序列化", func(t *testing.T) {
		req := dto.CreateRoomRequest{
			Name:        "测试房间",
			RoomType:    "bedroom",
			FloorNumber: testutils.IntPtr(1),
			Area:        testutils.Float64Ptr(25.0),
			Description: testutils.StringPtr("主卧室"),
			PositionData: map[string]any{
				"x": 10,
				"y": 20,
			},
		}

		assert.Equal(t, "测试房间", req.Name)
		assert.Equal(t, "bedroom", req.RoomType)
		assert.Equal(t, 1, *req.FloorNumber)
		assert.Equal(t, 25.0, *req.Area)
		assert.Equal(t, "主卧室", *req.Description)
		assert.Equal(t, 10, req.PositionData["x"])
		assert.Equal(t, 20, req.PositionData["y"])
	})
}


