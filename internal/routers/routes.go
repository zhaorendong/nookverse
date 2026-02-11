package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"nookverse/internal/services"
	"nookverse/pkg/api/v1/handlers"
)

// SetupRoutes 设置路由
func SetupRoutes(itemService services.ItemService, houseService services.HouseService) *gin.Engine {
	// 创建gin引擎
	r := gin.Default()

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Nookverse service is running",
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 物品管理路由
		itemHandler := handlers.NewItemHandler(itemService)
		items := v1.Group("/items")
		{
			items.POST("", itemHandler.CreateItem)
			items.GET("", itemHandler.ListItems)
			items.GET("/search", itemHandler.SearchItems)
			// 物品层级管理
			items.POST("/:itemId/move", itemHandler.MoveItem)
			items.POST("/:itemId/reminders", itemHandler.CreateReminder)
			
			// 单个物品操作
			items.GET("/:itemId", itemHandler.GetItem)
			items.PUT("/:itemId", itemHandler.UpdateItem)
			items.DELETE("/:itemId", itemHandler.DeleteItem)
			
			// 容器内容查询（特殊路由）
			items.GET("/container/:containerId/contents", itemHandler.GetContainerItems)
			
			// 提醒管理
			items.GET("/reminders/upcoming", itemHandler.GetUpcomingReminders)
			
			// 统计信息
			items.GET("/statistics", itemHandler.GetItemStatistics)
		}

		// 房间相关路由
		rooms := v1.Group("/rooms")
		{
			rooms.GET("/:roomId/items", itemHandler.GetItemsByRoom)
		}

		// 房屋管理路由
		houseHandler := handlers.NewHouseHandler(houseService)
		houses := v1.Group("/houses")
		{
			houses.POST("", houseHandler.CreateHouse)
			houses.GET("", houseHandler.ListHouses)
			houses.GET("/search", houseHandler.SearchHouses)
			
			// 单个房屋操作
			houses.GET("/:houseId", houseHandler.GetHouse)
			houses.PUT("/:houseId", houseHandler.UpdateHouse)
			houses.DELETE("/:houseId", houseHandler.DeleteHouse)
			
			// 房屋内房间管理
			houses.POST("/:houseId/rooms", houseHandler.CreateRoom)
			houses.GET("/:houseId/rooms", houseHandler.GetRoomsByHouse)
			
			// 统计信息
			houses.GET("/statistics", houseHandler.GetHouseStatistics)
		}

		// 房间独立操作路由
		independentRooms := v1.Group("/rooms")
		{
			independentRooms.GET("/:roomId", houseHandler.GetRoom)
			independentRooms.PUT("/:roomId", houseHandler.UpdateRoom)
			independentRooms.DELETE("/:roomId", houseHandler.DeleteRoom)
		}

		// TODO: 用户管理路由（待实现）
		// userHandler := handlers.NewUserHandler()
		// users := v1.Group("/users")
		// {
		// 	users.POST("/register", userHandler.Register)
		// 	users.POST("/login", userHandler.Login)
		// 	users.GET("/:id", userHandler.GetUserByID)
		// 	users.PUT("/:id", userHandler.UpdateUser)
		// 	users.DELETE("/:id", userHandler.DeleteUser)
		// }

		// TODO: 需要认证的路由组（待实现）
		// auth := v1.Group("/")
		// auth.Use(AuthMiddleware()) // 添加认证中间件
		// {
		// 	profile := auth.Group("/profile")
		// 	{
		// 		profile.GET("", userHandler.GetProfile)
		// 		profile.PUT("", userHandler.UpdateProfile)
		// 	}
		// }
	}

	return r
}

// AuthMiddleware 认证中间件示例
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里实现JWT认证逻辑
		// 从Authorization头中获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// 验证token逻辑...
		// 如果验证失败：
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// c.Abort()
		// return

		// 如果验证成功，设置用户信息到上下文
		// c.Set("user_id", userID)
		
		c.Next()
	}
}