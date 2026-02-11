-- NookVerse 数据库初始化脚本
-- PostgreSQL 数据库结构

-- 扩展安装
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- 用于模糊搜索

-- 1. 房屋表
CREATE TABLE IF NOT EXISTS houses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    address TEXT,
    description TEXT,
    area DECIMAL(10,2), -- 面积（平方米）
    floor_count INTEGER DEFAULT 1, -- 楼层数
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- 2. 房间表
CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    house_id UUID REFERENCES houses(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    room_type VARCHAR(50) NOT NULL, -- bedroom, living_room, kitchen, bathroom, study, etc.
    floor_number INTEGER DEFAULT 1, -- 楼层号
    area DECIMAL(8,2), -- 面积
    description TEXT,
    position_data JSONB DEFAULT '{}', -- 3D坐标和边界信息
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 3. 物品类别表
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL, -- 支持分类层级
    icon VARCHAR(50),
    color VARCHAR(20) DEFAULT '#666666',
    sort_order INTEGER DEFAULT 0, -- 排序
    is_system BOOLEAN DEFAULT FALSE, -- 是否系统分类
    created_at TIMESTAMP DEFAULT NOW()
);

-- 4. 物品表（核心）
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    room_id UUID REFERENCES rooms(id) ON DELETE SET NULL,
    container_id UUID REFERENCES items(id) ON DELETE SET NULL, -- 容器关系（物品可以包含其他物品）
    quantity INTEGER DEFAULT 1,
    status VARCHAR(20) DEFAULT 'active', -- active, archived, discarded, borrowed
    
    -- 重要属性
    expire_date DATE, -- 过期日期
    purchase_date DATE, -- 购买日期
    price DECIMAL(10,2), -- 价格
    warranty_period INTEGER, -- 保修期（月）
    brand VARCHAR(100), -- 品牌
    model VARCHAR(100), -- 型号
    
    -- 位置详情
    position JSONB DEFAULT '{}', -- 相对位置（x,y,z坐标或描述性位置）
    custom_position TEXT, -- 用户自定义位置描述
    
    -- 扩展属性
    attributes JSONB DEFAULT '{}', -- 存储品牌、型号、颜色等
    labels TEXT[], -- 标签数组
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 5. 媒体文件表
CREATE TABLE IF NOT EXISTS media_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID REFERENCES items(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    thumbnail_url TEXT,
    file_type VARCHAR(20) NOT NULL, -- image, video, document
    file_size BIGINT, -- 文件大小（字节）
    mime_type VARCHAR(100),
    alt_text TEXT, -- 替代文本
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 6. 提醒表
CREATE TABLE IF NOT EXISTS reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID REFERENCES items(id) ON DELETE CASCADE,
    reminder_type VARCHAR(20) NOT NULL, -- expire, maintenance, warranty, custom
    trigger_time TIMESTAMP NOT NULL,
    message TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, sent, completed, cancelled
    notify_channels TEXT[], -- notification channels: app, email, sms, voice
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 7. 用户表
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    nickname VARCHAR(100),
    avatar_url TEXT,
    phone VARCHAR(20),
    status INTEGER DEFAULT 1, -- 1:正常 2:禁用
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 8. 家庭表
CREATE TABLE IF NOT EXISTS families (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    invite_code VARCHAR(20) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 9. 家庭成员表
CREATE TABLE IF NOT EXISTS family_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id UUID REFERENCES families(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member', -- owner, admin, member, viewer
    joined_at TIMESTAMP DEFAULT NOW(),
    status INTEGER DEFAULT 1, -- 1:正常 2:禁用
    
    UNIQUE(family_id, user_id)
);

-- 10. 物品权限表
CREATE TABLE IF NOT EXISTS item_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID REFERENCES items(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    permission_level VARCHAR(20) DEFAULT 'view', -- owner, edit, view
    granted_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(item_id, user_id)
);

-- 11. 操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    item_id UUID REFERENCES items(id) ON DELETE SET NULL,
    operation_type VARCHAR(50) NOT NULL, -- create, update, delete, move, borrow
    description TEXT,
    old_value JSONB,
    new_value JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 12. 物品层级关系表（闭包表）
CREATE TABLE IF NOT EXISTS item_hierarchy (
    ancestor_id UUID REFERENCES items(id) ON DELETE CASCADE,
    descendant_id UUID REFERENCES items(id) ON DELETE CASCADE,
    depth INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    
    PRIMARY KEY (ancestor_id, descendant_id)
);

-- 13. 搜索索引表（用于全文搜索）
CREATE TABLE IF NOT EXISTS search_index (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID REFERENCES items(id) ON DELETE CASCADE,
    searchable_content TSVECTOR,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_rooms_house ON rooms(house_id);
CREATE INDEX IF NOT EXISTS idx_rooms_type ON rooms(room_type);
CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_items_room ON items(room_id);
CREATE INDEX IF NOT EXISTS idx_items_container ON items(container_id);
CREATE INDEX IF NOT EXISTS idx_items_category ON items(category_id);
CREATE INDEX IF NOT EXISTS idx_items_expire ON items(expire_date);
CREATE INDEX IF NOT EXISTS idx_items_status ON items(status);
CREATE INDEX IF NOT EXISTS idx_items_labels ON items USING GIN(labels);
CREATE INDEX IF NOT EXISTS idx_media_item ON media_files(item_id);
CREATE INDEX IF NOT EXISTS idx_media_type ON media_files(file_type);
CREATE INDEX IF NOT EXISTS idx_reminders_item ON reminders(item_id);
CREATE INDEX IF NOT EXISTS idx_reminders_type ON reminders(reminder_type);
CREATE INDEX IF NOT EXISTS idx_reminders_status ON reminders(status);
CREATE INDEX IF NOT EXISTS idx_reminders_trigger ON reminders(trigger_time);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_families_owner ON families(owner_id);
CREATE INDEX IF NOT EXISTS idx_family_members_family ON family_members(family_id);
CREATE INDEX IF NOT EXISTS idx_family_members_user ON family_members(user_id);
CREATE INDEX IF NOT EXISTS idx_item_permissions_item ON item_permissions(item_id);
CREATE INDEX IF NOT EXISTS idx_item_permissions_user ON item_permissions(user_id);
CREATE INDEX IF NOT EXISTS idx_logs_user ON operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_logs_item ON operation_logs(item_id);
CREATE INDEX IF NOT EXISTS idx_logs_operation ON operation_logs(operation_type);
CREATE INDEX IF NOT EXISTS idx_logs_created ON operation_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_hierarchy_ancestor ON item_hierarchy(ancestor_id);
CREATE INDEX IF NOT EXISTS idx_hierarchy_descendant ON item_hierarchy(descendant_id);
CREATE INDEX IF NOT EXISTS idx_search_item ON search_index(item_id);
CREATE INDEX IF NOT EXISTS idx_search_content ON search_index USING GIN(searchable_content);

-- 创建触发器函数
-- 更新时间戳
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建更新时间戳触发器
DROP TRIGGER IF EXISTS update_houses_updated_at ON houses;
CREATE TRIGGER update_houses_updated_at 
    BEFORE UPDATE ON houses 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_rooms_updated_at ON rooms;
CREATE TRIGGER update_rooms_updated_at 
    BEFORE UPDATE ON rooms 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_items_updated_at ON items;
CREATE TRIGGER update_items_updated_at 
    BEFORE UPDATE ON items 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_reminders_updated_at ON reminders;
CREATE TRIGGER update_reminders_updated_at 
    BEFORE UPDATE ON reminders 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_search_index_updated_at ON search_index;
CREATE TRIGGER update_search_index_updated_at 
    BEFORE UPDATE ON search_index 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 创建全文搜索函数
CREATE OR REPLACE FUNCTION update_search_index()
RETURNS TRIGGER AS $$
BEGIN
    -- 先删除现有的搜索索引
    DELETE FROM search_index WHERE item_id = NEW.id;
    
    -- 插入新的搜索索引
    INSERT INTO search_index (item_id, searchable_content)
    VALUES (
        NEW.id,
        setweight(to_tsvector('simple', COALESCE(NEW.name, '')), 'A') ||
        setweight(to_tsvector('simple', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('simple', array_to_string(COALESCE(NEW.labels, '{}'), ' ')), 'C')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建搜索索引触发器
DROP TRIGGER IF EXISTS trigger_update_search_index ON items;
CREATE TRIGGER trigger_update_search_index
    AFTER INSERT OR UPDATE ON items
    FOR EACH ROW EXECUTE FUNCTION update_search_index();

-- 清空现有数据以确保一致性（按依赖顺序删除）
DELETE FROM reminders;
DELETE FROM item_permissions;
DELETE FROM media_files;
DELETE FROM items;
DELETE FROM item_hierarchy;
DELETE FROM rooms;
DELETE FROM houses;
DELETE FROM family_members;
DELETE FROM families;
DELETE FROM categories WHERE is_system = TRUE;
DELETE FROM users WHERE username = 'admin';

-- 插入初始数据
-- 系统分类
INSERT INTO categories (name, icon, color, is_system, sort_order) VALUES
('电子设备', '📱', '#FF6B6B', TRUE, 1),
('家具', '🪑', '#4ECDC4', TRUE, 2),
('服装', '👕', '#45B7D1', TRUE, 3),
('书籍', '📚', '#96CEB4', TRUE, 4),
('食品', '🍎', '#FFEAA7', TRUE, 5),
('药品', '💊', '#DDA0DD', TRUE, 6),
('工具', '🔧', '#D9B573', TRUE, 7),
('运动用品', '⚽', '#FF8A80', TRUE, 8),
('化妆品', '💄', '#FFB6C1', TRUE, 9),
('其他', '📦', '#9E9E9E', TRUE, 10);

-- 系统用户（管理员）
INSERT INTO users (username, email, password_hash, nickname, status) VALUES
('admin', 'admin@nookverse.com', '$2a$10$abcdefghijklmnopqrstuvABCDEFGHIJKLMNO', '系统管理员', 1);

-- 示例家庭
INSERT INTO families (name, description, owner_id) 
SELECT '默认家庭', '系统默认家庭', id FROM users WHERE username = 'admin';

-- 示例房屋
INSERT INTO houses (name, address, description) VALUES
('示例住宅', '北京市朝阳区示例街道123号', '这是一个示例住宅');

-- 示例房间
INSERT INTO rooms (house_id, name, room_type, floor_number, description) 
SELECT h.id, r.name, r.room_type, r.floor_number, r.description
FROM houses h
CROSS JOIN (
    VALUES 
        ('主卧室', 'bedroom', 1, '主人的卧室'),
        ('客厅', 'living_room', 1, '家庭聚会的主要场所'),
        ('厨房', 'kitchen', 1, '烹饪美食的地方')
) AS r(name, room_type, floor_number, description)
WHERE h.name = '示例住宅';

-- 示例物品
INSERT INTO items (name, description, category_id, room_id, quantity, price, purchase_date, brand) 
SELECT 'iPhone 15', '苹果最新款手机', 
       (SELECT id FROM categories WHERE name = '电子设备'),
       (SELECT r.id FROM rooms r JOIN houses h ON r.house_id = h.id WHERE h.name = '示例住宅' AND r.name = '主卧室'),
       1, 6999.00, '2024-01-15', 'Apple';

INSERT INTO items (name, description, category_id, room_id, quantity, price, purchase_date, brand) 
SELECT '真皮沙发', '客厅主沙发',
       (SELECT id FROM categories WHERE name = '家具'),
       (SELECT r.id FROM rooms r JOIN houses h ON r.house_id = h.id WHERE h.name = '示例住宅' AND r.name = '客厅'),
       1, 8999.00, '2023-12-01', '顾家家居';

INSERT INTO items (name, description, category_id, room_id, quantity, price, expire_date, brand) 
SELECT '牛奶', '纯牛奶',
       (SELECT id FROM categories WHERE name = '食品'),
       (SELECT r.id FROM rooms r JOIN houses h ON r.house_id = h.id WHERE h.name = '示例住宅' AND r.name = '厨房'),
       2, 12.50, '2024-03-01', '伊利';

-- 示例提醒
INSERT INTO reminders (item_id, reminder_type, trigger_time, message, notify_channels)
SELECT i.id, 'expire', '2024-02-28 09:00:00', '牛奶即将过期，请及时处理', ARRAY['app', 'email']
FROM items i WHERE i.name = '牛奶';