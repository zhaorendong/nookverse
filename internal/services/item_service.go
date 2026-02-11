package services

import (
	"context"
	"errors"
	"time"

	"nookverse/internal/models"
	"gorm.io/gorm"
)

// ItemService 物品服务接口
type ItemService interface {
	// 基础CRUD操作
	CreateItem(ctx context.Context, item *models.Item) error
	GetItemByID(ctx context.Context, id string) (*models.Item, error)
	UpdateItem(ctx context.Context, item *models.Item) error
	DeleteItem(ctx context.Context, id string) error
	
	// 查询操作
	ListItems(ctx context.Context, filters ItemFilters) ([]models.Item, int64, error)
	SearchItems(ctx context.Context, query string, filters ItemFilters) ([]models.Item, int64, error)
	GetItemsByRoom(ctx context.Context, roomID string) ([]models.Item, error)
	GetItemsByCategory(ctx context.Context, categoryID string) ([]models.Item, error)
	
	// 层级管理
	GetItemHierarchy(ctx context.Context, itemID string) ([]models.Item, error)
	MoveItemToContainer(ctx context.Context, itemID, containerID string) error
	GetContainerItems(ctx context.Context, containerID string) ([]models.Item, error)
	
	// 提醒管理
	CreateReminder(ctx context.Context, reminder *models.Reminder) error
	GetUpcomingReminders(ctx context.Context, days int) ([]models.Reminder, error)
	
	// 统计分析
	GetItemStatistics(ctx context.Context, userID string) (*ItemStatistics, error)
}

// ItemFilters 物品查询过滤条件
type ItemFilters struct {
	UserID      *string
	RoomID      *string
	CategoryID  *string
	Status      *string
	Labels      []string
	ExpireDate  *time.Time
	MinPrice    *float64
	MaxPrice    *float64
	Page        int
	PageSize    int
	OrderBy     string
}

// ItemStatistics 物品统计信息
type ItemStatistics struct {
	TotalItems     int64            `json:"total_items"`
	ByStatus       map[string]int64 `json:"by_status"`
	ByCategory     map[string]int64 `json:"by_category"`
	TotalValue     float64          `json:"total_value"`
	ExpiringSoon   int64            `json:"expiring_soon"` // 30天内过期
	LowStockItems  int64            `json:"low_stock_items"` // 数量小于等于1的物品
}

type itemService struct {
	db *gorm.DB
}

// NewItemService 创建物品服务实例
func NewItemService(db *gorm.DB) ItemService {
	return &itemService{db: db}
}

// CreateItem 创建物品
func (s *itemService) CreateItem(ctx context.Context, item *models.Item) error {
	if item.Name == "" {
		return errors.New("物品名称不能为空")
	}

	// 验证关联关系
	if item.RoomID != nil {
		var room models.Room
		if err := s.db.First(&room, "id = ?", *item.RoomID).Error; err != nil {
			return errors.New("指定的房间不存在")
		}
	}

	if item.CategoryID != nil {
		var category models.Category
		if err := s.db.First(&category, "id = ?", *item.CategoryID).Error; err != nil {
			return errors.New("指定的分类不存在")
		}
	}

	// 设置默认值
	if item.Quantity <= 0 {
		item.Quantity = 1
	}
	if item.Status == "" {
		item.Status = "active"
	}

	return s.db.WithContext(ctx).Create(item).Error
}

// GetItemByID 根据ID获取物品
func (s *itemService) GetItemByID(ctx context.Context, id string) (*models.Item, error) {
	var item models.Item
	err := s.db.WithContext(ctx).
		Preload("Category").
		Preload("Room").
		Preload("Container").
		Preload("MediaFiles").
		Preload("Reminders").
		First(&item, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("物品不存在")
		}
		return nil, err
	}
	
	return &item, nil
}

// UpdateItem 更新物品
func (s *itemService) UpdateItem(ctx context.Context, item *models.Item) error {
	if item.ID == "" {
		return errors.New("物品ID不能为空")
	}

	// 检查物品是否存在
	var existing models.Item
	if err := s.db.First(&existing, "id = ?", item.ID).Error; err != nil {
		return errors.New("物品不存在")
	}

	return s.db.WithContext(ctx).Save(item).Error
}

// DeleteItem 删除物品
func (s *itemService) DeleteItem(ctx context.Context, id string) error {
	// 检查是否有子物品
	var count int64
	s.db.Model(&models.Item{}).Where("container_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该物品包含其他物品，不能直接删除")
	}

	result := s.db.WithContext(ctx).Delete(&models.Item{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("物品不存在")
	}
	
	return nil
}

// ListItems 列出物品
func (s *itemService) ListItems(ctx context.Context, filters ItemFilters) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	query := s.db.WithContext(ctx).Model(&models.Item{})

	// 应用过滤条件
	if filters.UserID != nil {
		// TODO: 根据用户权限过滤
	}
	
	if filters.RoomID != nil {
		query = query.Where("room_id = ?", *filters.RoomID)
	}
	
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	
	if len(filters.Labels) > 0 {
		for _, label := range filters.Labels {
			query = query.Where("? = ANY(labels)", label)
		}
	}
	
	if filters.ExpireDate != nil {
		query = query.Where("expire_date <= ?", *filters.ExpireDate)
	}
	
	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}
	
	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
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

	// 预加载关联数据
	err := query.
		Preload("Category").
		Preload("Room").
		Preload("Container").
		Find(&items).Error
	
	return items, total, err
}

// SearchItems 搜索物品
func (s *itemService) SearchItems(ctx context.Context, query string, filters ItemFilters) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	dbQuery := s.db.WithContext(ctx).Model(&models.Item{})

	// 构建搜索条件
	searchCondition := s.db.Where(
		s.db.Where("name ILIKE ?", "%"+query+"%").
			Or("description ILIKE ?", "%"+query+"%").
			Or("brand ILIKE ?", "%"+query+"%").
			Or("model ILIKE ?", "%"+query+"%"),
	)

	// 应用过滤条件
	if filters.RoomID != nil {
		searchCondition = searchCondition.Where("room_id = ?", *filters.RoomID)
	}
	
	if filters.CategoryID != nil {
		searchCondition = searchCondition.Where("category_id = ?", *filters.CategoryID)
	}
	
	if filters.Status != nil {
		searchCondition = searchCondition.Where("status = ?", *filters.Status)
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
	err := searchCondition.
		Preload("Category").
		Preload("Room").
		Preload("Container").
		Find(&items).Error
	
	return items, total, err
}

// GetItemsByRoom 获取房间内物品
func (s *itemService) GetItemsByRoom(ctx context.Context, roomID string) ([]models.Item, error) {
	var items []models.Item
	err := s.db.WithContext(ctx).
		Where("room_id = ? AND container_id IS NULL", roomID).
		Preload("Category").
		Preload("MediaFiles").
		Order("name").
		Find(&items).Error
	
	return items, err
}

// GetItemsByCategory 获取分类下物品
func (s *itemService) GetItemsByCategory(ctx context.Context, categoryID string) ([]models.Item, error) {
	var items []models.Item
	err := s.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Preload("Room").
		Preload("Category").
		Order("name").
		Find(&items).Error
	
	return items, err
}

// GetItemHierarchy 获取物品层级关系
func (s *itemService) GetItemHierarchy(ctx context.Context, itemID string) ([]models.Item, error) {
	var ancestors []models.Item
	
	// 查询祖先节点
	err := s.db.WithContext(ctx).
		Raw(`
			SELECT i.* FROM items i
			INNER JOIN item_hierarchy ih ON i.id = ih.ancestor_id
			WHERE ih.descendant_id = ?
			ORDER BY ih.depth DESC
		`, itemID).
		Scan(&ancestors).Error
	
	return ancestors, err
}

// MoveItemToContainer 移动物品到容器
func (s *itemService) MoveItemToContainer(ctx context.Context, itemID, containerID string) error {
	// 验证容器存在且不是自己
	if itemID == containerID {
		return errors.New("不能将物品移动到自身")
	}

	var container models.Item
	if err := s.db.First(&container, "id = ?", containerID).Error; err != nil {
		return errors.New("目标容器不存在")
	}

	// 检查是否会形成循环引用
	var count int64
	s.db.Model(&models.ItemHierarchy{}).
		Where("ancestor_id = ? AND descendant_id = ?", itemID, containerID).
		Count(&count)
	
	if count > 0 {
		return errors.New("不能形成循环引用")
	}

	// 更新物品的容器ID
	return s.db.WithContext(ctx).
		Model(&models.Item{}).
		Where("id = ?", itemID).
		Update("container_id", containerID).Error
}

// GetContainerItems 获取容器内物品
func (s *itemService) GetContainerItems(ctx context.Context, containerID string) ([]models.Item, error) {
	var items []models.Item
	err := s.db.WithContext(ctx).
		Where("container_id = ?", containerID).
		Preload("Category").
		Preload("MediaFiles").
		Order("name").
		Find(&items).Error
	
	return items, err
}

// CreateReminder 创建提醒
func (s *itemService) CreateReminder(ctx context.Context, reminder *models.Reminder) error {
	if reminder.ItemID == "" {
		return errors.New("物品ID不能为空")
	}

	// 验证物品存在
	var item models.Item
	if err := s.db.First(&item, "id = ?", reminder.ItemID).Error; err != nil {
		return errors.New("物品不存在")
	}

	// 验证提醒时间
	if reminder.TriggerTime.Before(time.Now()) {
		return errors.New("提醒时间不能早于当前时间")
	}

	if reminder.Status == "" {
		reminder.Status = "pending"
	}

	return s.db.WithContext(ctx).Create(reminder).Error
}

// GetUpcomingReminders 获取即将到来的提醒
func (s *itemService) GetUpcomingReminders(ctx context.Context, days int) ([]models.Reminder, error) {
	var reminders []models.Reminder
	
	endTime := time.Now().AddDate(0, 0, days)
	
	err := s.db.WithContext(ctx).
		Where("status = ? AND trigger_time BETWEEN ? AND ?", 
			"pending", time.Now(), endTime).
		Preload("Item").
		Order("trigger_time").
		Find(&reminders).Error
	
	return reminders, err
}

// GetItemStatistics 获取物品统计信息
func (s *itemService) GetItemStatistics(ctx context.Context, userID string) (*ItemStatistics, error) {
	stats := &ItemStatistics{
		ByStatus:   make(map[string]int64),
		ByCategory: make(map[string]int64),
	}

	// 获取总数
	s.db.WithContext(ctx).Model(&models.Item{}).Count(&stats.TotalItems)

	// 按状态统计
	s.db.WithContext(ctx).
		Model(&models.Item{}).
		Select("status, count(*)").
		Group("status").
		Scan(&stats.ByStatus)

	// 按分类统计
	s.db.WithContext(ctx).
		Model(&models.Item{}).
		Select("c.name, count(*)").
		Joins("LEFT JOIN categories c ON items.category_id = c.id").
		Group("c.name").
		Scan(&stats.ByCategory)

	// 计算总价值
	var totalValue float64
	s.db.WithContext(ctx).
		Model(&models.Item{}).
		Select("COALESCE(SUM(price * quantity), 0)").
		Scan(&totalValue)
	stats.TotalValue = totalValue

	// 统计即将过期的物品（30天内）
	expireTime := time.Now().AddDate(0, 0, 30)
	s.db.WithContext(ctx).
		Model(&models.Item{}).
		Where("expire_date IS NOT NULL AND expire_date <= ?", expireTime).
		Count(&stats.ExpiringSoon)

	// 统计低库存物品
	s.db.WithContext(ctx).
		Model(&models.Item{}).
		Where("quantity <= 1").
		Count(&stats.LowStockItems)

	return stats, nil
}