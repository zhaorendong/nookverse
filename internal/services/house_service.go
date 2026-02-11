package services

import (
	"context"
	"errors"

	"nookverse/internal/models"
	"gorm.io/gorm"
)

// HouseService 房屋服务接口
type HouseService interface {
	// 房屋基础CRUD操作
	CreateHouse(ctx context.Context, house *models.House) error
	GetHouseByID(ctx context.Context, id string) (*models.House, error)
	UpdateHouse(ctx context.Context, house *models.House) error
	DeleteHouse(ctx context.Context, id string) error
	ListHouses(ctx context.Context, filters HouseFilters) ([]models.House, int64, error)
	
	// 房间管理
	CreateRoom(ctx context.Context, room *models.Room) error
	GetRoomByID(ctx context.Context, id string) (*models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) error
	DeleteRoom(ctx context.Context, id string) error
	GetRoomsByHouse(ctx context.Context, houseID string) ([]models.Room, error)
	
	// 统计分析
	GetHouseStatistics(ctx context.Context) (*HouseStatistics, error)
	
	// 搜索功能
	SearchHouses(ctx context.Context, query string, filters HouseFilters) ([]models.House, int64, error)
}

// HouseFilters 房屋查询过滤条件
type HouseFilters struct {
	UserID     *string
	Name       *string
	Address    *string
	MinArea    *float64
	MaxArea    *float64
	MinFloors  *int
	MaxFloors  *int
	Page       int
	PageSize   int
	OrderBy    string
}

// HouseStatistics 房屋统计信息
type HouseStatistics struct {
	TotalHouses   int64            `json:"total_houses"`
	TotalRooms    int64            `json:"total_rooms"`
	AverageArea   float64          `json:"average_area"`
	ByFloorCount  map[int]int64    `json:"by_floor_count"`
	ByRoomType    map[string]int64 `json:"by_room_type"`
}

type houseService struct {
	db *gorm.DB
}

// NewHouseService 创建房屋服务实例
func NewHouseService(db *gorm.DB) HouseService {
	return &houseService{db: db}
}

// CreateHouse 创建房屋
func (s *houseService) CreateHouse(ctx context.Context, house *models.House) error {
	if house.Name == "" {
		return errors.New("房屋名称不能为空")
	}

	// 设置默认值
	if house.FloorCount <= 0 {
		house.FloorCount = 1
	}

	return s.db.WithContext(ctx).Create(house).Error
}

// GetHouseByID 根据ID获取房屋
func (s *houseService) GetHouseByID(ctx context.Context, id string) (*models.House, error) {
	var house models.House
	err := s.db.WithContext(ctx).
		Preload("Rooms").
		Preload("Rooms.Items").
		Preload("Rooms.Items.Category").
		First(&house, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("房屋不存在")
		}
		return nil, err
	}
	
	return &house, nil
}

// UpdateHouse 更新房屋
func (s *houseService) UpdateHouse(ctx context.Context, house *models.House) error {
	if house.ID == "" {
		return errors.New("房屋ID不能为空")
	}

	// 检查房屋是否存在
	var existing models.House
	if err := s.db.First(&existing, "id = ?", house.ID).Error; err != nil {
		return errors.New("房屋不存在")
	}

	return s.db.WithContext(ctx).Save(house).Error
}

// DeleteHouse 删除房屋
func (s *houseService) DeleteHouse(ctx context.Context, id string) error {
	// 检查是否有房间
	var count int64
	s.db.Model(&models.Room{}).Where("house_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该房屋包含房间，不能直接删除")
	}

	result := s.db.WithContext(ctx).Delete(&models.House{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("房屋不存在")
	}
	
	return nil
}

// ListHouses 列出房屋
func (s *houseService) ListHouses(ctx context.Context, filters HouseFilters) ([]models.House, int64, error) {
	var houses []models.House
	var total int64

	query := s.db.WithContext(ctx).Model(&models.House{})

	// 应用过滤条件
	if filters.UserID != nil {
		// TODO: 根据用户权限过滤
	}
	
	if filters.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filters.Name+"%")
	}
	
	if filters.Address != nil {
		query = query.Where("address ILIKE ?", "%"+*filters.Address+"%")
	}
	
	if filters.MinArea != nil {
		query = query.Where("area >= ?", *filters.MinArea)
	}
	
	if filters.MaxArea != nil {
		query = query.Where("area <= ?", *filters.MaxArea)
	}
	
	if filters.MinFloors != nil {
		query = query.Where("floor_count >= ?", *filters.MinFloors)
	}
	
	if filters.MaxFloors != nil {
		query = query.Where("floor_count <= ?", *filters.MaxFloors)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用分页和排序
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}
	
	offset := (filters.Page - 1) * filters.PageSize
	query = query.Offset(offset).Limit(filters.PageSize)
	
	orderBy := "created_at DESC"
	if filters.OrderBy != "" {
		orderBy = filters.OrderBy
	}
	query = query.Order(orderBy)

	// 预加载房间数据
	err := query.Preload("Rooms").Find(&houses).Error
	
	return houses, total, err
}

// CreateRoom 创建房间
func (s *houseService) CreateRoom(ctx context.Context, room *models.Room) error {
	if room.Name == "" {
		return errors.New("房间名称不能为空")
	}

	if room.HouseID == "" {
		return errors.New("房屋ID不能为空")
	}

	// 验证房屋存在
	var house models.House
	if err := s.db.First(&house, "id = ?", room.HouseID).Error; err != nil {
		return errors.New("指定的房屋不存在")
	}

	// 设置默认值
	if room.FloorNumber <= 0 {
		room.FloorNumber = 1
	}

	if room.RoomType == "" {
		room.RoomType = "other"
	}

	return s.db.WithContext(ctx).Create(room).Error
}

// GetRoomByID 根据ID获取房间
func (s *houseService) GetRoomByID(ctx context.Context, id string) (*models.Room, error) {
	var room models.Room
	err := s.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Category").
		Preload("Items.MediaFiles").
		First(&room, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("房间不存在")
		}
		return nil, err
	}
	
	return &room, nil
}

// UpdateRoom 更新房间
func (s *houseService) UpdateRoom(ctx context.Context, room *models.Room) error {
	if room.ID == "" {
		return errors.New("房间ID不能为空")
	}

	// 检查房间是否存在
	var existing models.Room
	if err := s.db.First(&existing, "id = ?", room.ID).Error; err != nil {
		return errors.New("房间不存在")
	}

	return s.db.WithContext(ctx).Save(room).Error
}

// DeleteRoom 删除房间
func (s *houseService) DeleteRoom(ctx context.Context, id string) error {
	// 检查是否有物品
	var count int64
	s.db.Model(&models.Item{}).Where("room_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该房间包含物品，不能直接删除")
	}

	result := s.db.WithContext(ctx).Delete(&models.Room{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("房间不存在")
	}
	
	return nil
}

// GetRoomsByHouse 获取房屋内房间
func (s *houseService) GetRoomsByHouse(ctx context.Context, houseID string) ([]models.Room, error) {
	var rooms []models.Room
	err := s.db.WithContext(ctx).
		Where("house_id = ?", houseID).
		Preload("Items").
		Preload("Items.Category").
		Order("floor_number, name").
		Find(&rooms).Error
	
	return rooms, err
}

// GetHouseStatistics 获取房屋统计信息
func (s *houseService) GetHouseStatistics(ctx context.Context) (*HouseStatistics, error) {
	stats := &HouseStatistics{
		ByFloorCount: make(map[int]int64),
		ByRoomType:   make(map[string]int64),
	}

	// 获取房屋总数
	s.db.WithContext(ctx).Model(&models.House{}).Count(&stats.TotalHouses)

	// 获取房间总数
	s.db.WithContext(ctx).Model(&models.Room{}).Count(&stats.TotalRooms)

	// 计算平均面积
	var avgArea float64
	s.db.WithContext(ctx).
		Model(&models.House{}).
		Select("COALESCE(AVG(area), 0)").
		Scan(&avgArea)
	stats.AverageArea = avgArea

	// 按楼层数统计
	s.db.WithContext(ctx).
		Model(&models.House{}).
		Select("floor_count, count(*)").
		Group("floor_count").
		Scan(&stats.ByFloorCount)

	// 按房间类型统计
	s.db.WithContext(ctx).
		Model(&models.Room{}).
		Select("room_type, count(*)").
		Group("room_type").
		Scan(&stats.ByRoomType)

	return stats, nil
}

// SearchHouses 搜索房屋
func (s *houseService) SearchHouses(ctx context.Context, query string, filters HouseFilters) ([]models.House, int64, error) {
	var houses []models.House
	var total int64

	dbQuery := s.db.WithContext(ctx).Model(&models.House{})

	// 构建搜索条件
	searchCondition := s.db.Where(
		s.db.Where("name ILIKE ?", "%"+query+"%").
			Or("address ILIKE ?", "%"+query+"%").
			Or("description ILIKE ?", "%"+query+"%"),
	)

	// 应用过滤条件
	if filters.MinArea != nil {
		searchCondition = searchCondition.Where("area >= ?", *filters.MinArea)
	}
	
	if filters.MaxArea != nil {
		searchCondition = searchCondition.Where("area <= ?", *filters.MaxArea)
	}
	
	if filters.MinFloors != nil {
		searchCondition = searchCondition.Where("floor_count >= ?", *filters.MinFloors)
	}
	
	if filters.MaxFloors != nil {
		searchCondition = searchCondition.Where("floor_count <= ?", *filters.MaxFloors)
	}

	// 获取总数
	if err := dbQuery.Where(searchCondition).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用分页
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}
	
	offset := (filters.Page - 1) * filters.PageSize
	searchCondition = searchCondition.Offset(offset).Limit(filters.PageSize)

	// 执行查询
	err := searchCondition.Preload("Rooms").Find(&houses).Error
	
	return houses, total, err
}