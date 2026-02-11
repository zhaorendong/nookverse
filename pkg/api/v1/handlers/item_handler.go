package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"nookverse/internal/models"
	"nookverse/internal/services"
	"nookverse/pkg/api/v1/dto"
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

// stringPtr 返回字符串指针
func stringPtr(s string) *string {
	return &s
}

// intPtr 返回整数指针
func intPtr(i int) *int {
	return &i
}

// float64Ptr 返回浮点数指针
func float64Ptr(f float64) *float64 {
	return &f
}

// ItemHandler 物品处理器
type ItemHandler struct {
	itemService services.ItemService
}

// NewItemHandler 创建物品处理器实例
func NewItemHandler(itemService services.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

// CreateItem 创建物品
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败: " + err.Error(),
		})
		return
	}

	// 转换DTO到模型
	item := &models.Item{
		Name:           req.Name,
		Description:    getValueOrDefault(req.Description, ""),
		CategoryID:     req.CategoryID,
		RoomID:         req.RoomID,
		ContainerID:    req.ContainerID,
		Quantity:       req.Quantity,
		Status:         req.Status,
		ExpireDate:     req.ExpireDate,
		PurchaseDate:   req.PurchaseDate,
		Price:          req.Price,
		WarrantyPeriod: req.WarrantyPeriod,
		Brand:          req.Brand,
		Model:          req.Model,
		Position:       req.Position,
		CustomPosition: req.CustomPosition,
		Attributes:     req.Attributes,
		Labels:         req.Labels,
	}

	// 创建物品
	if err := h.itemService.CreateItem(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建物品失败: " + err.Error(),
		})
		return
	}

	// 返回响应
	c.JSON(http.StatusCreated, gin.H{
		"message": "物品创建成功",
		"data":    dto.ToItemResponse(item),
	})
}

// GetItem 获取物品详情
func (h *ItemHandler) GetItem(c *gin.Context) {
	id := c.Param("itemId")
	
	// 验证参数
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "物品ID不能为空",
		})
		return
	}
	
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "物品ID格式不正确",
		})
		return
	}
	
	item, err := h.itemService.GetItemByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "物品不存在: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.ToItemResponse(item),
	})
}

// UpdateItem 更新物品
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	id := c.Param("itemId")
	
	// 获取现有物品
	existingItem, err := h.itemService.GetItemByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "物品不存在",
		})
		return
	}

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败: " + err.Error(),
		})
		return
	}

	// 更新字段
	if req.Name != nil {
		existingItem.Name = *req.Name
	}
	if req.Description != nil {
		existingItem.Description = *req.Description
	}
	if req.CategoryID != nil {
		existingItem.CategoryID = req.CategoryID
	}
	if req.RoomID != nil {
		existingItem.RoomID = req.RoomID
	}
	if req.ContainerID != nil {
		existingItem.ContainerID = req.ContainerID
	}
	if req.Quantity != nil {
		existingItem.Quantity = *req.Quantity
	}
	if req.Status != nil {
		existingItem.Status = *req.Status
	}
	if req.ExpireDate != nil {
		existingItem.ExpireDate = req.ExpireDate
	}
	if req.PurchaseDate != nil {
		existingItem.PurchaseDate = req.PurchaseDate
	}
	if req.Price != nil {
		existingItem.Price = req.Price
	}
	if req.WarrantyPeriod != nil {
		existingItem.WarrantyPeriod = req.WarrantyPeriod
	}
	if req.Brand != nil {
		existingItem.Brand = req.Brand
	}
	if req.Model != nil {
		existingItem.Model = req.Model
	}
	if req.Position != nil {
		existingItem.Position = *req.Position
	}
	if req.CustomPosition != nil {
		existingItem.CustomPosition = req.CustomPosition
	}
	if req.Attributes != nil {
		existingItem.Attributes = *req.Attributes
	}
	if req.Labels != nil {
		existingItem.Labels = *req.Labels
	}

	// 保存更新
	if err := h.itemService.UpdateItem(c.Request.Context(), existingItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新物品失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "物品更新成功",
		"data":    dto.ToItemResponse(existingItem),
	})
}

// DeleteItem 删除物品
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	id := c.Param("itemId")
	
	if err := h.itemService.DeleteItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除物品失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "物品删除成功",
	})
}

// ListItems 分页查询物品
func (h *ItemHandler) ListItems(c *gin.Context) {
	// 构建查询过滤条件
	var filters services.ItemFilters
	
	if roomID := c.Query("room_id"); roomID != "" {
		filters.RoomID = &roomID
	}
	
	if categoryID := c.Query("category_id"); categoryID != "" {
		filters.CategoryID = &categoryID
	}
	
	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}
	
	if labels := c.Query("labels"); labels != "" {
		// 将逗号分隔的标签转换为数组
		filters.Labels = []string{labels}
	}
	
	if page := c.Query("page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil {
			filters.Page = pageNum
		}
	}
	
	if pageSize := c.Query("page_size"); pageSize != "" {
		if size, err := strconv.Atoi(pageSize); err == nil {
			filters.PageSize = size
		}
	}
	
	if orderBy := c.Query("order_by"); orderBy != "" {
		filters.OrderBy = orderBy
	}

	// 执行查询
	items, total, err := h.itemService.ListItems(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询物品列表失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var responses []dto.ItemResponse
	for _, item := range items {
		responses = append(responses, dto.ToItemResponse(&item))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"total":     total,
			"page":      filters.Page,
			"page_size": filters.PageSize,
		},
	})
}

// SearchItems 搜索物品
func (h *ItemHandler) SearchItems(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "搜索关键词不能为空",
		})
		return
	}

	// 构建查询过滤条件
	var filters services.ItemFilters
	
	if page := c.Query("page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil {
			filters.Page = pageNum
		}
	}
	
	if pageSize := c.Query("page_size"); pageSize != "" {
		if size, err := strconv.Atoi(pageSize); err == nil {
			filters.PageSize = size
		}
	}

	// 执行搜索
	items, total, err := h.itemService.SearchItems(c.Request.Context(), query, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "搜索物品失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var responses []dto.ItemResponse
	for _, item := range items {
		responses = append(responses, dto.ToItemResponse(&item))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"total":     total,
			"page":      filters.Page,
			"page_size": filters.PageSize,
		},
		"query": query,
	})
}

// GetItemsByRoom 获取房间内物品
func (h *ItemHandler) GetItemsByRoom(c *gin.Context) {
	roomID := c.Param("roomId")
	
	items, err := h.itemService.GetItemsByRoom(c.Request.Context(), roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取房间物品失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var responses []dto.ItemResponse
	for _, item := range items {
		responses = append(responses, dto.ToItemResponse(&item))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
	})
}

// MoveItem 移动物品到容器
func (h *ItemHandler) MoveItem(c *gin.Context) {
	itemID := c.Param("itemId")
	
	// 验证参数
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "物品ID不能为空",
		})
		return
	}
	
	if !isValidUUID(itemID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "物品ID格式不正确",
		})
		return
	}
	
	var req dto.MoveItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败: " + err.Error(),
		})
		return
	}

	if err := h.itemService.MoveItemToContainer(c.Request.Context(), itemID, req.ContainerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "移动物品失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "物品移动成功",
	})
}

// GetContainerItems 获取容器内物品
func (h *ItemHandler) GetContainerItems(c *gin.Context) {
	containerID := c.Param("containerId")
	
	items, err := h.itemService.GetContainerItems(c.Request.Context(), containerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取容器物品失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var responses []dto.ItemResponse
	for _, item := range items {
		responses = append(responses, dto.ToItemResponse(&item))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
	})
}

// CreateReminder 创建提醒
func (h *ItemHandler) CreateReminder(c *gin.Context) {
	itemID := c.Param("itemId")
	
	// 验证参数
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "物品ID不能为空",
		})
		return
	}
	
	if !isValidUUID(itemID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "物品ID格式不正确",
		})
		return
	}
	
	var req dto.CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败: " + err.Error(),
		})
		return
	}

	reminder := &models.Reminder{
		ItemID:         itemID,
		ReminderType:   req.ReminderType,
		TriggerTime:    req.TriggerTime,
		Message:        req.Message,
		NotifyChannels: req.NotifyChannels,
	}

	if err := h.itemService.CreateReminder(c.Request.Context(), reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建提醒失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "提醒创建成功",
		"data":    dto.ToReminderResponse(reminder),
	})
}

// GetUpcomingReminders 获取即将到来的提醒
func (h *ItemHandler) GetUpcomingReminders(c *gin.Context) {
	days := 7 // 默认7天
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
			days = d
		}
	}

	reminders, err := h.itemService.GetUpcomingReminders(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取提醒列表失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var responses []dto.ReminderResponse
	for _, reminder := range reminders {
		responses = append(responses, dto.ToReminderResponse(&reminder))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
	})
}

// GetItemStatistics 获取物品统计数据
func (h *ItemHandler) GetItemStatistics(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授权访问",
		})
		return
	}

	stats, err := h.itemService.GetItemStatistics(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取统计数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}