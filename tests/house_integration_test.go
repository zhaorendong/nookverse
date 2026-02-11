package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"nookverse/internal/database"
	"nookverse/internal/models"
	"nookverse/internal/routers"
	"nookverse/internal/services"
	"nookverse/pkg/api/v1/dto"
)

func TestHouseIntegration(t *testing.T) {
	t.Skip("跳过需要数据库的集成测试")
	
	// 连接测试数据库
	db, err := database.NewConnection(database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "nookverse_test",
		SSLMode:  "disable",
		TimeZone: "Asia/Shanghai",
	})
	require.NoError(t, err, "数据库连接失败")

	// 自动迁移表结构
	err = db.AutoMigrate(
		&models.House{},
		&models.Room{},
		&models.Category{},
		&models.Item{},
		&models.MediaFile{},
		&models.Reminder{},
	)
	require.NoError(t, err, "数据库迁移失败")

	// 清理测试数据
	cleanupTestData(db)

	// 初始化服务
	houseService := services.NewHouseService(db)

	// 设置路由
	router := routers.SetupRoutes(nil, houseService)

	t.Run("创建房屋", func(t *testing.T) {
		houseReq := dto.CreateHouseRequest{
			Name:        "测试房屋",
			Address:     stringPtr("测试地址123号"),
			Description: stringPtr("这是一个测试房屋"),
			Area:        float64Ptr(120.5),
			FloorCount:  intPtr(2),
			Metadata: map[string]any{
				"year_built": 2020,
				"has_garden": true,
			},
		}

		jsonData, _ := json.Marshal(houseReq)
		req, _ := http.NewRequest("POST", "/api/v1/houses", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data, ok := response["data"].(map[string]any)
		require.True(t, ok, "响应中应该包含data字段")
		assert.Equal(t, "测试房屋", data["name"])
		assert.Equal(t, "测试地址123号", data["address"])
		assert.Equal(t, 2, int(data["floor_count"].(float64)))
	})

	t.Run("获取房屋列表", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/houses", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data, ok := response["data"].([]any)
		require.True(t, ok, "响应中应该包含data数组")
		assert.GreaterOrEqual(t, len(data), 1, "应该至少有一个房屋")
	})

	t.Run("创建房间", func(t *testing.T) {
		// 先获取一个房屋ID
		var house models.House
		err := db.First(&house).Error
		require.NoError(t, err)

		roomReq := dto.CreateRoomRequest{
			Name:        "测试房间",
			RoomType:    "bedroom",
			FloorNumber: intPtr(1),
			Area:        float64Ptr(25.0),
			Description: stringPtr("主卧室"),
			PositionData: map[string]any{
				"x": 10,
				"y": 20,
				"z": 0,
			},
		}

		jsonData, _ := json.Marshal(roomReq)
		url := fmt.Sprintf("/api/v1/houses/%s/rooms", house.ID)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data, ok := response["data"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "测试房间", data["name"])
		assert.Equal(t, "bedroom", data["room_type"])
	})

	t.Run("获取房屋详情", func(t *testing.T) {
		var house models.House
		err := db.Preload("Rooms").First(&house).Error
		require.NoError(t, err)

		url := fmt.Sprintf("/api/v1/houses/%s", house.ID)
		req, _ := http.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data, ok := response["data"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, house.Name, data["name"])

		// 检查房间是否被预加载
		rooms, ok := data["rooms"].([]any)
		require.True(t, ok)
		assert.GreaterOrEqual(t, len(rooms), 1)
	})

	t.Run("更新房屋", func(t *testing.T) {
		var house models.House
		err := db.First(&house).Error
		require.NoError(t, err)

		updateReq := dto.UpdateHouseRequest{
			Name:        stringPtr("更新后的房屋名称"),
			Description: stringPtr("更新后的描述"),
		}

		jsonData, _ := json.Marshal(updateReq)
		url := fmt.Sprintf("/api/v1/houses/%s", house.ID)
		req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证更新是否成功
		var updatedHouse models.House
		err = db.First(&updatedHouse, "id = ?", house.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "更新后的房屋名称", updatedHouse.Name)
	})

	t.Run("搜索房屋", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/houses/search?q=测试", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data, ok := response["data"].([]any)
		require.True(t, ok)
		assert.GreaterOrEqual(t, len(data), 1)
	})

	t.Run("获取房屋统计", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/houses/statistics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data, ok := response["data"].(map[string]any)
		require.True(t, ok)

		totalHouses, ok := data["total_houses"].(float64)
		require.True(t, ok)
		assert.GreaterOrEqual(t, int(totalHouses), 1)

		totalRooms, ok := data["total_rooms"].(float64)
		require.True(t, ok)
		assert.GreaterOrEqual(t, int(totalRooms), 1)
	})

	t.Run("删除房间", func(t *testing.T) {
		var room models.Room
		err := db.First(&room).Error
		require.NoError(t, err)

		url := fmt.Sprintf("/api/v1/rooms/%s", room.ID)
		req, _ := http.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证房间已被删除
		var count int64
		db.Model(&models.Room{}).Where("id = ?", room.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("删除房屋", func(t *testing.T) {
		// 先删除所有房间
		db.Delete(&models.Room{}, "1=1")

		var house models.House
		err := db.First(&house).Error
		require.NoError(t, err)

		url := fmt.Sprintf("/api/v1/houses/%s", house.ID)
		req, _ := http.NewRequest("DELETE", url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证房屋已被删除
		var count int64
		db.Model(&models.House{}).Where("id = ?", house.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

// 辅助函数
func cleanupTestData(db *gorm.DB) {
	// 按正确的依赖顺序删除数据
	db.Exec("DELETE FROM reminders")
	db.Exec("DELETE FROM media_files")
	db.Exec("DELETE FROM items")
	db.Exec("DELETE FROM rooms")
	db.Exec("DELETE FROM houses")
	db.Exec("DELETE FROM categories WHERE is_system = false")
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}