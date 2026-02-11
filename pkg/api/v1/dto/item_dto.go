package dto

import (
	"time"

	"nookverse/internal/models"
)

// CreateItemRequest 创建物品请求
type CreateItemRequest struct {
	Name           string            `json:"name" binding:"required"`
	Description    *string           `json:"description,omitempty"`
	CategoryID     *string           `json:"category_id,omitempty"`
	RoomID         *string           `json:"room_id,omitempty"`
	ContainerID    *string           `json:"container_id,omitempty"`
	Quantity       int               `json:"quantity" binding:"required"`
	Status         string            `json:"status,omitempty"`
	ExpireDate     *time.Time        `json:"expire_date,omitempty"`
	PurchaseDate   *time.Time        `json:"purchase_date,omitempty"`
	Price          *float64          `json:"price,omitempty"`
	WarrantyPeriod *int              `json:"warranty_period,omitempty"`
	Brand          *string           `json:"brand,omitempty"`
	Model          *string           `json:"model,omitempty"`
	Position       map[string]any    `json:"position,omitempty"`
	CustomPosition *string           `json:"custom_position,omitempty"`
	Attributes     map[string]any    `json:"attributes,omitempty"`
	Labels         []string          `json:"labels,omitempty"`
}

// UpdateItemRequest 更新物品请求
type UpdateItemRequest struct {
	Name           *string           `json:"name,omitempty"`
	Description    *string           `json:"description,omitempty"`
	CategoryID     *string           `json:"category_id,omitempty"`
	RoomID         *string           `json:"room_id,omitempty"`
	ContainerID    *string           `json:"container_id,omitempty"`
	Quantity       *int              `json:"quantity,omitempty"`
	Status         *string           `json:"status,omitempty"`
	ExpireDate     *time.Time        `json:"expire_date,omitempty"`
	PurchaseDate   *time.Time        `json:"purchase_date,omitempty"`
	Price          *float64          `json:"price,omitempty"`
	WarrantyPeriod *int              `json:"warranty_period,omitempty"`
	Brand          *string           `json:"brand,omitempty"`
	Model          *string           `json:"model,omitempty"`
	Position       *map[string]any   `json:"position,omitempty"`
	CustomPosition *string           `json:"custom_position,omitempty"`
	Attributes     *map[string]any   `json:"attributes,omitempty"`
	Labels         *[]string         `json:"labels,omitempty"`
}

// ItemResponse 物品响应
type ItemResponse struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    *string           `json:"description,omitempty"`
	Category       *CategoryResponse `json:"category,omitempty"`
	Room           *RoomResponse     `json:"room,omitempty"`
	Container      *ItemResponse     `json:"container,omitempty"`
	Quantity       int               `json:"quantity"`
	Status         string            `json:"status"`
	ExpireDate     *time.Time        `json:"expire_date,omitempty"`
	PurchaseDate   *time.Time        `json:"purchase_date,omitempty"`
	Price          *float64          `json:"price,omitempty"`
	WarrantyPeriod *int              `json:"warranty_period,omitempty"`
	Brand          *string           `json:"brand,omitempty"`
	Model          *string           `json:"model,omitempty"`
	Position       map[string]any    `json:"position,omitempty"`
	CustomPosition *string           `json:"custom_position,omitempty"`
	Attributes     map[string]any    `json:"attributes,omitempty"`
	Labels         []string          `json:"labels,omitempty"`
	MediaFiles     []MediaFileResponse `json:"media_files,omitempty"`
	Reminders      []ReminderResponse  `json:"reminders,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// CategoryResponse 分类响应
type CategoryResponse struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	ParentID  *string             `json:"parent_id,omitempty"`
	Icon      string              `json:"icon"`
	Color     string              `json:"color"`
	SortOrder int                 `json:"sort_order"`
	Children  []CategoryResponse  `json:"children,omitempty"`
}

// RoomResponse 房间响应
type RoomResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	RoomType    string    `json:"room_type"`
	FloorNumber int       `json:"floor_number"`
	Area        *float64  `json:"area,omitempty"`
	Description *string   `json:"description,omitempty"`
}

// MediaFileResponse 媒体文件响应
type MediaFileResponse struct {
	ID          string    `json:"id"`
	FileURL     string    `json:"file_url"`
	ThumbnailURL *string  `json:"thumbnail_url,omitempty"`
	FileType    string    `json:"file_type"`
	FileSize    *int64    `json:"file_size,omitempty"`
	MimeType    *string   `json:"mime_type,omitempty"`
	AltText     *string   `json:"alt_text,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

// ReminderResponse 提醒响应
type ReminderResponse struct {
	ID             string    `json:"id"`
	ReminderType   string    `json:"reminder_type"`
	TriggerTime    time.Time `json:"trigger_time"`
	Message        string    `json:"message"`
	Status         string    `json:"status"`
	NotifyChannels []string  `json:"notify_channels"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// MoveItemRequest 移动物品请求
type MoveItemRequest struct {
	ContainerID string `json:"container_id" binding:"required"`
}

// CreateReminderRequest 创建提醒请求
type CreateReminderRequest struct {
	ReminderType   string    `json:"reminder_type" binding:"required"`
	TriggerTime    time.Time `json:"trigger_time" binding:"required"`
	Message        string    `json:"message" binding:"required"`
	NotifyChannels []string  `json:"notify_channels"`
}

// ToItemResponse 转换物品模型为响应格式
func ToItemResponse(item *models.Item) ItemResponse {
	resp := ItemResponse{
		ID:             item.ID,
		Name:           item.Name,
		Description:    &item.Description,
		Quantity:       item.Quantity,
		Status:         item.Status,
		ExpireDate:     item.ExpireDate,
		PurchaseDate:   item.PurchaseDate,
		Price:          item.Price,
		WarrantyPeriod: item.WarrantyPeriod,
		Brand:          item.Brand,
		Model:          item.Model,
		Position:       item.Position,
		CustomPosition: item.CustomPosition,
		Attributes:     item.Attributes,
		Labels:         item.Labels,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}

	// 转换分类信息
	if item.Category != nil {
		resp.Category = &CategoryResponse{
			ID:        item.Category.ID,
			Name:      item.Category.Name,
			ParentID:  item.Category.ParentID,
			Icon:      item.Category.Icon,
			Color:     item.Category.Color,
			SortOrder: item.Category.SortOrder,
		}
	}

	if item.Room != nil {
		resp.Room = &RoomResponse{
			ID:          item.Room.ID,
			Name:        item.Room.Name,
			RoomType:    item.Room.RoomType,
			FloorNumber: item.Room.FloorNumber,
			Area:        &item.Room.Area,
			Description: &item.Room.Description,
		}
	}

	if item.Container != nil {
		resp.Container = &ItemResponse{
			ID:   item.Container.ID,
			Name: item.Container.Name,
		}
	}

	// 转换媒体文件
	for _, media := range item.MediaFiles {
		resp.MediaFiles = append(resp.MediaFiles, MediaFileResponse{
			ID:          media.ID,
			FileURL:     media.FileURL,
			ThumbnailURL: &media.ThumbnailURL,
			FileType:    media.FileType,
			FileSize:    media.FileSize,
			MimeType:    media.MimeType,
			AltText:     media.AltText,
			SortOrder:   media.SortOrder,
			CreatedAt:   media.CreatedAt,
		})
	}

	// 转换提醒
	for _, reminder := range item.Reminders {
		resp.Reminders = append(resp.Reminders, ReminderResponse{
			ID:             reminder.ID,
			ReminderType:   reminder.ReminderType,
			TriggerTime:    reminder.TriggerTime,
			Message:        reminder.Message,
			Status:         reminder.Status,
			NotifyChannels: reminder.NotifyChannels,
			CreatedAt:      reminder.CreatedAt,
			UpdatedAt:      reminder.UpdatedAt,
		})
	}

	return resp
}

// ToReminderResponse 转换提醒模型为响应格式
func ToReminderResponse(reminder *models.Reminder) ReminderResponse {
	return ReminderResponse{
		ID:             reminder.ID,
		ReminderType:   reminder.ReminderType,
		TriggerTime:    reminder.TriggerTime,
		Message:        reminder.Message,
		Status:         reminder.Status,
		NotifyChannels: reminder.NotifyChannels,
		CreatedAt:      reminder.CreatedAt,
		UpdatedAt:      reminder.UpdatedAt,
	}
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Data       []ItemResponse `json:"data"`
	Query      string         `json:"query"`
	Pagination Pagination     `json:"pagination"`
}

// Pagination 分页信息
type Pagination struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// StatisticsResponse 统计响应
type StatisticsResponse struct {
	TotalItems     int64            `json:"total_items"`
	ByStatus       map[string]int64 `json:"by_status"`
	ByCategory     map[string]int64 `json:"by_category"`
	TotalValue     float64          `json:"total_value"`
	ExpiringSoon   int64            `json:"expiring_soon"`
	LowStockItems  int64            `json:"low_stock_items"`
}