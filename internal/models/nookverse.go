package models

import (
	"time"
)

// House 房屋模型
type House struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Address     string         `json:"address" gorm:"type:text"`
	Description string         `json:"description" gorm:"type:text"`
	Area        float64        `json:"area" gorm:"type:decimal(10,2)"` // 面积（平方米）
	FloorCount  int            `json:"floor_count" gorm:"default:1"`   // 楼层数
	Metadata    map[string]any `json:"metadata" gorm:"type:jsonb"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	Rooms []Room `json:"rooms" gorm:"foreignKey:HouseID"`
}

// Room 房间模型
type Room struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	HouseID     string         `json:"house_id" gorm:"type:uuid;not null;index"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	RoomType    string         `json:"room_type" gorm:"size:50;not null"` // bedroom, living_room, kitchen等
	FloorNumber int            `json:"floor_number" gorm:"default:1"`     // 楼层号
	Area        float64        `json:"area" gorm:"type:decimal(8,2)"`     // 面积
	Description string         `json:"description" gorm:"type:text"`
	PositionData map[string]any `json:"position_data" gorm:"type:jsonb"` // 3D坐标和边界信息
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	Items []Item `json:"items" gorm:"foreignKey:RoomID"`
}

// Category 物品类别模型
type Category struct {
	ID       string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name     string    `json:"name" gorm:"size:100;not null"`
	ParentID *string   `json:"parent_id" gorm:"type:uuid;index"`
	Icon     string    `json:"icon" gorm:"size:50"`
	Color    string    `json:"color" gorm:"size:20;default:'#666666'"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	IsSystem bool      `json:"is_system" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`

	Children  []Category `json:"children" gorm:"foreignKey:ParentID"`
	Items     []Item     `json:"items" gorm:"foreignKey:CategoryID"`
}

// Item 物品模型（核心）
type Item struct {
	ID             string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name           string         `json:"name" gorm:"size:200;not null"`
	Description    string         `json:"description" gorm:"type:text"`
	CategoryID     *string        `json:"category_id" gorm:"type:uuid;index"`
	RoomID         *string        `json:"room_id" gorm:"type:uuid;index"`
	ContainerID    *string        `json:"container_id" gorm:"type:uuid;index"` // 容器关系
	Quantity       int            `json:"quantity" gorm:"default:1"`
	Status         string         `json:"status" gorm:"size:20;default:'active'"` // active, archived, discarded, borrowed

	// 重要属性
	ExpireDate     *time.Time     `json:"expire_date" gorm:"index"`
	PurchaseDate   *time.Time     `json:"purchase_date"`
	Price          *float64       `json:"price" gorm:"type:decimal(10,2)"`
	WarrantyPeriod *int           `json:"warranty_period"` // 保修期（月）
	Brand          *string        `json:"brand" gorm:"size:100"`
	Model          *string        `json:"model" gorm:"size:100"`

	// 位置详情
	Position       map[string]any `json:"position" gorm:"type:jsonb"` // 相对位置
	CustomPosition *string        `json:"custom_position" gorm:"type:text"`

	// 扩展属性
	Attributes     map[string]any `json:"attributes" gorm:"type:jsonb"`
	Labels         []string       `json:"labels" gorm:"type:text[]"`

	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`

	// 关联关系
	Category       *Category      `json:"category" gorm:"foreignKey:CategoryID"`
	Room           *Room          `json:"room" gorm:"foreignKey:RoomID"`
	Container      *Item          `json:"container" gorm:"foreignKey:ContainerID"`
	ContainedItems []Item         `json:"contained_items" gorm:"foreignKey:ContainerID"`
	MediaFiles     []MediaFile    `json:"media_files" gorm:"foreignKey:ItemID"`
	Reminders      []Reminder     `json:"reminders" gorm:"foreignKey:ItemID"`
}

// MediaFile 媒体文件模型
type MediaFile struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItemID      string    `json:"item_id" gorm:"type:uuid;not null;index"`
	FileURL     string    `json:"file_url" gorm:"type:text;not null"`
	ThumbnailURL string   `json:"thumbnail_url" gorm:"type:text"`
	FileType    string    `json:"file_type" gorm:"size:20;not null"` // image, video, document
	FileSize    *int64    `json:"file_size"`                          // 文件大小（字节）
	MimeType    *string   `json:"mime_type" gorm:"size:100"`
	AltText     *string   `json:"alt_text" gorm:"type:text"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`

	Item *Item `json:"item" gorm:"foreignKey:ItemID"`
}

// Reminder 提醒模型
type Reminder struct {
	ID             string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItemID         string    `json:"item_id" gorm:"type:uuid;not null;index"`
	ReminderType   string    `json:"reminder_type" gorm:"size:20;not null"` // expire, maintenance, warranty, custom
	TriggerTime    time.Time `json:"trigger_time" gorm:"not null;index"`
	Message        string    `json:"message" gorm:"type:text;not null"`
	Status         string    `json:"status" gorm:"size:20;default:'pending';index"` // pending, sent, completed, cancelled
	NotifyChannels []string  `json:"notify_channels" gorm:"type:text[]"`            // app, email, sms, voice
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Item *Item `json:"item" gorm:"foreignKey:ItemID"`
}

// User 用户模型
type User struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string    `json:"username" gorm:"size:50;uniqueIndex;not null"`
	Email        string    `json:"email" gorm:"size:100;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"` // 不序列化密码
	Nickname     *string   `json:"nickname" gorm:"size:100"`
	AvatarURL    *string   `json:"avatar_url" gorm:"type:text"`
	Phone        *string   `json:"phone" gorm:"size:20"`
	Status       int       `json:"status" gorm:"default:1"` // 1:正常 2:禁用
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Families []FamilyMember `json:"families" gorm:"foreignKey:UserID"`
	Items    []ItemPermission `json:"items" gorm:"foreignKey:UserID"`
}

// Family 家庭模型
type Family struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description *string   `json:"description" gorm:"type:text"`
	OwnerID     string    `json:"owner_id" gorm:"type:uuid;not null;index"`
	InviteCode  *string   `json:"invite_code" gorm:"size:20;uniqueIndex"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Owner    User           `json:"owner" gorm:"foreignKey:OwnerID"`
	Members  []FamilyMember `json:"members" gorm:"foreignKey:FamilyID"`
	Houses   []House        `json:"houses" gorm:"many2many:family_houses"`
}

// FamilyMember 家庭成员模型
type FamilyMember struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FamilyID  string    `json:"family_id" gorm:"type:uuid;not null;index"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
	Role      string    `json:"role" gorm:"size:20;default:'member'"` // owner, admin, member, viewer
	JoinedAt  time.Time `json:"joined_at"`
	Status    int       `json:"status" gorm:"default:1"` // 1:正常 2:禁用

	Family Family `json:"family" gorm:"foreignKey:FamilyID"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`
}

// ItemPermission 物品权限模型
type ItemPermission struct {
	ID             string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItemID         string    `json:"item_id" gorm:"type:uuid;not null;index"`
	UserID         string    `json:"user_id" gorm:"type:uuid;not null;index"`
	PermissionLevel string   `json:"permission_level" gorm:"size:20;default:'view'"` // owner, edit, view
	GrantedBy      *string   `json:"granted_by" gorm:"type:uuid"`
	CreatedAt      time.Time `json:"created_at"`

	Item User `json:"item" gorm:"foreignKey:ItemID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// OperationLog 操作日志模型
type OperationLog struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      *string        `json:"user_id" gorm:"type:uuid;index"`
	ItemID      *string        `json:"item_id" gorm:"type:uuid;index"`
	OperationType string       `json:"operation_type" gorm:"size:50;not null;index"`
	Description *string        `json:"description" gorm:"type:text"`
	OldValue    map[string]any `json:"old_value" gorm:"type:jsonb"`
	NewValue    map[string]any `json:"new_value" gorm:"type:jsonb"`
	IPAddress   *string        `json:"ip_address" gorm:"type:inet"`
	UserAgent   *string        `json:"user_agent" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"index"`

	User *User `json:"user" gorm:"foreignKey:UserID"`
	Item *Item `json:"item" gorm:"foreignKey:ItemID"`
}

// ItemHierarchy 物品层级关系模型（闭包表）
type ItemHierarchy struct {
	AncestorID  string    `json:"ancestor_id" gorm:"type:uuid;not null;primaryKey"`
	DescendantID string   `json:"descendant_id" gorm:"type:uuid;not null;primaryKey"`
	Depth       int       `json:"depth" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`

	Ancestor  Item `json:"ancestor" gorm:"foreignKey:AncestorID"`
	Descendant Item `json:"descendant" gorm:"foreignKey:DescendantID"`
}

// TableName 指定表名
func (House) TableName() string { return "houses" }
func (Room) TableName() string { return "rooms" }
func (Category) TableName() string { return "categories" }
func (Item) TableName() string { return "items" }
func (MediaFile) TableName() string { return "media_files" }
func (Reminder) TableName() string { return "reminders" }
func (User) TableName() string { return "users" }
func (Family) TableName() string { return "families" }
func (FamilyMember) TableName() string { return "family_members" }
func (ItemPermission) TableName() string { return "item_permissions" }
func (OperationLog) TableName() string { return "operation_logs" }
func (ItemHierarchy) TableName() string { return "item_hierarchy" }