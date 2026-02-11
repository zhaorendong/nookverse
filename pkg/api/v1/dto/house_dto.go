package dto

import (
	"time"

	"nookverse/internal/models"
)

// CreateHouseRequest 创建房屋请求
type CreateHouseRequest struct {
	Name        string         `json:"name" binding:"required"`
	Address     *string        `json:"address,omitempty"`
	Description *string        `json:"description,omitempty"`
	Area        *float64       `json:"area,omitempty"`
	FloorCount  *int           `json:"floor_count,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// UpdateHouseRequest 更新房屋请求
type UpdateHouseRequest struct {
	Name        *string        `json:"name,omitempty"`
	Address     *string        `json:"address,omitempty"`
	Description *string        `json:"description,omitempty"`
	Area        *float64       `json:"area,omitempty"`
	FloorCount  *int           `json:"floor_count,omitempty"`
	Metadata    *map[string]any `json:"metadata,omitempty"`
}

// HouseResponse 房屋响应
type HouseResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Address     *string            `json:"address,omitempty"`
	Description *string            `json:"description,omitempty"`
	Area        *float64           `json:"area,omitempty"`
	FloorCount  int                `json:"floor_count"`
	Metadata    map[string]any     `json:"metadata,omitempty"`
	Rooms       []HouseRoomResponse `json:"rooms,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	Name        string         `json:"name" binding:"required"`
	RoomType    string         `json:"room_type" binding:"required"`
	FloorNumber *int           `json:"floor_number,omitempty"`
	Area        *float64       `json:"area,omitempty"`
	Description *string        `json:"description,omitempty"`
	PositionData map[string]any `json:"position_data,omitempty"`
}

// UpdateRoomRequest 更新房间请求
type UpdateRoomRequest struct {
	Name        *string        `json:"name,omitempty"`
	RoomType    *string        `json:"room_type,omitempty"`
	FloorNumber *int           `json:"floor_number,omitempty"`
	Area        *float64       `json:"area,omitempty"`
	Description *string        `json:"description,omitempty"`
	PositionData *map[string]any `json:"position_data,omitempty"`
}

// HouseRoomResponse 房间响应
type HouseRoomResponse struct {
	ID          string         `json:"id"`
	HouseID     string         `json:"house_id"`
	Name        string         `json:"name"`
	RoomType    string         `json:"room_type"`
	FloorNumber int            `json:"floor_number"`
	Area        *float64       `json:"area,omitempty"`
	Description *string        `json:"description,omitempty"`
	PositionData map[string]any `json:"position_data,omitempty"`
	Items       []ItemResponse `json:"items,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// HouseStatisticsResponse 房屋统计响应
type HouseStatisticsResponse struct {
	TotalHouses   int64            `json:"total_houses"`
	TotalRooms    int64            `json:"total_rooms"`
	AverageArea   float64          `json:"average_area"`
	ByFloorCount  map[int]int64    `json:"by_floor_count"`
	ByRoomType    map[string]int64 `json:"by_room_type"`
}

// ToHouseResponse 转换房屋模型为响应格式
func ToHouseResponse(house *models.House) HouseResponse {
	resp := HouseResponse{
		ID:          house.ID,
		Name:        house.Name,
		Address:     &house.Address,
		Description: &house.Description,
		Area:        &house.Area,
		FloorCount:  house.FloorCount,
		Metadata:    house.Metadata,
		CreatedAt:   house.CreatedAt,
		UpdatedAt:   house.UpdatedAt,
	}

	// 转换房间信息
	for _, room := range house.Rooms {
		resp.Rooms = append(resp.Rooms, ToHouseRoomResponse(&room))
	}

	return resp
}

// ToHouseRoomResponse 转换房间模型为响应格式
func ToHouseRoomResponse(room *models.Room) HouseRoomResponse {
	resp := HouseRoomResponse{
		ID:          room.ID,
		HouseID:     room.HouseID,
		Name:        room.Name,
		RoomType:    room.RoomType,
		FloorNumber: room.FloorNumber,
		Area:        &room.Area,
		Description: &room.Description,
		PositionData: room.PositionData,
		CreatedAt:   room.CreatedAt,
		UpdatedAt:   room.UpdatedAt,
	}

	// 转换物品信息
	for _, item := range room.Items {
		resp.Items = append(resp.Items, ToItemResponse(&item))
	}

	return resp
}