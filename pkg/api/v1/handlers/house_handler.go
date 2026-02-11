package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"nookverse/internal/models"
	"nookverse/internal/services"
	"nookverse/pkg/api/v1/dto"
)

// HouseHandler 房屋处理器
type HouseHandler struct {
	houseService services.HouseService
}

// NewHouseHandler 创建房屋处理器实例
func NewHouseHandler(houseService services.HouseService) *HouseHandler {
	return &HouseHandler{
		houseService: houseService,
	}
}

// CreateHouse 创建房屋
func (h *HouseHandler) CreateHouse(c *gin.Context) {
	var req dto.CreateHouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败: " + err.Error(),
		})
		return
	}

	// 转换DTO到模型
	house := &models.House{
		Name:        req.Name,
		Address:     getValueOrDefault(req.Address, ""),
		Description: getValueOrDefault(req.Description, ""),
		Area:        getValueOrDefault(req.Area, 0),
		FloorCount:  getValueOrDefault(req.FloorCount, 1),
		Metadata:    req.Metadata,
	}

	// 创建房屋
	if err := h.houseService.CreateHouse(c.Request.Context(), house); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建房屋失败: " + err.Error(),
		})
		return
	}

	// 返回响应
	c.JSON(http.StatusCreated, gin.H{
		"message": "房屋创建成功",
		"data":    dto.ToHouseResponse(house),
	})
}

// GetHouse 获取房屋详情
func (h *HouseHandler) GetHouse(c *gin.Context) {
	id := c.Param("houseId")
	
	// 验证参数
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "房屋ID不能为空",
		})
		return
	}
	
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "房屋ID格式不正确",
		})
		return
	}
	
	house, err := h.houseService.GetHouseByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "房屋不存在: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.ToHouseResponse(house),
	})
}

// UpdateHouse 更新房屋
func (h *HouseHandler) UpdateHouse(c *gin.Context) {
	id := c.Param("houseId")
	
	// 获取现有房屋
	existingHouse, err := h.houseService.GetHouseByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "房屋不存在",
		})
		return
	}

	var req dto.UpdateHouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败: " + err.Error(),
		})
		return
	}

	// 更新字段
	if req.Name != nil {
		existingHouse.Name = *req.Name
	}
	if req.Address != nil {
		existingHouse.Address = *req.Address
	}
	if req.Description != nil {
		existingHouse.Description = *req.Description
	}
	if req.Area != nil {
		existingHouse.Area = *req.Area
	}
	if req.FloorCount != nil {
		existingHouse.FloorCount = *req.FloorCount
	}
	if req.Metadata != nil {
		existingHouse.Metadata = *req.Metadata
	}

	// 保存更新
	if err := h.houseService.UpdateHouse(c.Request.Context(), existingHouse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新房屋失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "房屋更新成功",
		"data":    dto.ToHouseResponse(existingHouse),
	})
}

// DeleteHouse 删除房屋
func (h *HouseHandler) DeleteHouse(c *gin.Context) {
	id := c.Param("houseId")
	
	if err := h.houseService.DeleteHouse(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除房屋失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "房屋删除成功",
	})
}

// ListHouses 分页查询房屋
func (h *HouseHandler) ListHouses(c *gin.Context) {
	// 构建查询过滤条件
	var filters services.HouseFilters
	
	if name := c.Query("name"); name != "" {
		filters.Name = &name
	}
	
	if address := c.Query("address"); address != "" {
		filters.Address = &address
	}
	
	if minArea := c.Query("min_area"); minArea != "" {
		if area, err := strconv.ParseFloat(minArea, 64); err == nil {
			filters.MinArea = &area
		}
	}
	
	if maxArea := c.Query("max_area"); maxArea != "" {
		if area, err := strconv.ParseFloat(maxArea, 64); err == nil {
			filters.MaxArea = &area
		}
	}
	
	if minFloors := c.Query("min_floors"); minFloors != "" {
		if floors, err := strconv.Atoi(minFloors); err == nil {
			filters.MinFloors = &floors
		}
	}
	
	if maxFloors := c.Query("max_floors"); maxFloors != "" {
		if floors, err := strconv.Atoi(maxFloors); err == nil {
			filters.MaxFloors = &floors
		}
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
	houses, total, err := h.houseService.ListHouses(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询房屋列表失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var responses []dto.HouseResponse
	for _, house := range houses {
		responses = append(responses, dto.ToHouseResponse(&house))
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

// SearchHouses 搜索房屋
func (h *HouseHandler) SearchHouses(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "搜索关键词不能为空",
		})
		return
	}

	// 构建查询过滤条件
	var filters services.HouseFilters
	
	if minArea := c.Query("min_area"); minArea != "" {
		if area, err := strconv.ParseFloat(minArea, 64); err == nil {
			filters.MinArea = &area
		}
	}
	
	if maxArea := c.Query("max_area"); maxArea != "" {
		if area, err := strconv.ParseFloat(maxArea, 64); err == nil {
			filters.MaxArea = &area
		}
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

	// 执行搜索
	houses, total, err := h.houseService.SearchHouses(c.Request.Context(), query, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "搜索房屋失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var responses []dto.HouseResponse
	for _, house := range houses {
		responses = append(responses, dto.ToHouseResponse(&house))
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

// CreateRoom 创建房间
func (h *HouseHandler) CreateRoom(c *gin.Context) {
	houseID := c.Param("houseId")
	
	// 验证参数
	if houseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "房屋ID不能为空",
		})
		return
	}
	
	if !isValidUUID(houseID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "房屋ID格式不正确",
		})
		return
	}
	
	var req dto.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败: " + err.Error(),
		})
		return
	}

	room := &models.Room{
		HouseID:      houseID,
		Name:         req.Name,
		RoomType:     req.RoomType,
		FloorNumber:  getValueOrDefault(req.FloorNumber, 1),
		Area:         getValueOrDefault(req.Area, 0),
		Description:  getValueOrDefault(req.Description, ""),
		PositionData: req.PositionData,
	}

	if err := h.houseService.CreateRoom(c.Request.Context(), room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建房间失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "房间创建成功",
		"data":    dto.ToHouseRoomResponse(room),
	})
}

// GetRoom 获取房间详情
func (h *HouseHandler) GetRoom(c *gin.Context) {
	id := c.Param("roomId")
	
	// 验证参数
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "房间ID不能为空",
		})
		return
	}
	
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "房间ID格式不正确",
		})
		return
	}
	
	room, err := h.houseService.GetRoomByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "房间不存在: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.ToHouseRoomResponse(room),
	})
}

// UpdateRoom 更新房间
func (h *HouseHandler) UpdateRoom(c *gin.Context) {
	id := c.Param("roomId")
	
	// 获取现有房间
	existingRoom, err := h.houseService.GetRoomByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "房间不存在",
		})
		return
	}

	var req dto.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数验证失败: " + err.Error(),
		})
		return
	}

	// 更新字段
	if req.Name != nil {
		existingRoom.Name = *req.Name
	}
	if req.RoomType != nil {
		existingRoom.RoomType = *req.RoomType
	}
	if req.FloorNumber != nil {
		existingRoom.FloorNumber = *req.FloorNumber
	}
	if req.Area != nil {
		existingRoom.Area = *req.Area
	}
	if req.Description != nil {
		existingRoom.Description = *req.Description
	}
	if req.PositionData != nil {
		existingRoom.PositionData = *req.PositionData
	}

	// 保存更新
	if err := h.houseService.UpdateRoom(c.Request.Context(), existingRoom); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新房间失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "房间更新成功",
		"data":    dto.ToHouseRoomResponse(existingRoom),
	})
}

// DeleteRoom 删除房间
func (h *HouseHandler) DeleteRoom(c *gin.Context) {
	id := c.Param("roomId")
	
	if err := h.houseService.DeleteRoom(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除房间失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "房间删除成功",
	})
}

// GetRoomsByHouse 获取房屋内房间
func (h *HouseHandler) GetRoomsByHouse(c *gin.Context) {
	houseID := c.Param("houseId")
	
	rooms, err := h.houseService.GetRoomsByHouse(c.Request.Context(), houseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取房屋房间失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var responses []dto.HouseRoomResponse
	for _, room := range rooms {
		responses = append(responses, dto.ToHouseRoomResponse(&room))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
	})
}

// GetHouseStatistics 获取房屋统计数据
func (h *HouseHandler) GetHouseStatistics(c *gin.Context) {
	stats, err := h.houseService.GetHouseStatistics(c.Request.Context())
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