-- Emoji管理系统数据库表

-- 1. Emoji表情表
CREATE TABLE IF NOT EXISTS `emojis` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'emoji ID',
    `key` VARCHAR(50) NOT NULL UNIQUE COMMENT 'emoji键名，如e0,e1,e2',
    `filename` VARCHAR(255) NOT NULL COMMENT '原始文件名',
    `group_name` VARCHAR(100) NOT NULL COMMENT '组名，如系统基础表情',
    `sprite_group` INT NOT NULL COMMENT '雪碧图组号',
    `sprite_position_x` INT NOT NULL COMMENT '在雪碧图中的X位置',
    `sprite_position_y` INT NOT NULL COMMENT '在雪碧图中的Y位置',
    `file_size` INT DEFAULT 0 COMMENT '文件大小(字节)',
    `cdn_url` VARCHAR(500) DEFAULT '' COMMENT 'CDN地址',
    `upload_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1=active, 0=deleted(软删除)',
    `created_by` CHAR(36) COMMENT '上传者UUID',
    `deleted_at` TIMESTAMP NULL COMMENT '软删除时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX `idx_key` (`key`),
    INDEX `idx_group` (`group_name`),
    INDEX `idx_sprite_group` (`sprite_group`),
    INDEX `idx_status` (`status`),
    INDEX `idx_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Emoji表情表';

-- 2. Emoji组表
CREATE TABLE IF NOT EXISTS `emoji_groups` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '组ID',
    `group_name` VARCHAR(100) NOT NULL UNIQUE COMMENT '组名',
    `description` VARCHAR(500) DEFAULT '' COMMENT '组描述',
    `sort_order` INT DEFAULT 0 COMMENT '排序权重',
    `emoji_count` INT DEFAULT 0 COMMENT 'emoji数量',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1=active, 0=inactive',
    `created_by` CHAR(36) COMMENT '创建者UUID',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX `idx_group_name` (`group_name`),
    INDEX `idx_sort_order` (`sort_order`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Emoji组表';

-- 3. 雪碧图表
CREATE TABLE IF NOT EXISTS `emoji_sprites` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '雪碧图ID',
    `sprite_group` INT NOT NULL UNIQUE COMMENT '雪碧图组号',
    `filename` VARCHAR(255) NOT NULL COMMENT '雪碧图文件名',
    `cdn_url` VARCHAR(500) DEFAULT '' COMMENT 'CDN地址',
    `width` INT NOT NULL COMMENT '雪碧图宽度',
    `height` INT NOT NULL COMMENT '雪碧图高度',
    `emoji_count` INT NOT NULL COMMENT '包含的emoji数量',
    `file_size` INT DEFAULT 0 COMMENT '文件大小(字节)',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1=active, 0=inactive',
    `generated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '生成时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX `idx_sprite_group` (`sprite_group`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Emoji雪碧图表';

-- 4. 任务队列表（用于雪碧图生成等异步任务）
CREATE TABLE IF NOT EXISTS `emoji_tasks` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '任务ID',
    `task_type` VARCHAR(50) NOT NULL COMMENT '任务类型：regenerate_sprites, upload_emoji等',
    `status` VARCHAR(20) DEFAULT 'pending' COMMENT '任务状态：pending, running, completed, failed',
    `progress` INT DEFAULT 0 COMMENT '进度百分比 0-100',
    `message` VARCHAR(500) DEFAULT '' COMMENT '任务消息',
    `params` JSON COMMENT '任务参数',
    `result` JSON COMMENT '任务结果',
    `created_by` CHAR(36) COMMENT '创建者UUID',
    `started_at` TIMESTAMP NULL COMMENT '开始时间',
    `completed_at` TIMESTAMP NULL COMMENT '完成时间',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX `idx_task_type` (`task_type`),
    INDEX `idx_status` (`status`),
    INDEX `idx_created_by` (`created_by`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Emoji任务队列表';

-- 初始化默认emoji组
INSERT INTO `emoji_groups` (`group_name`, `description`, `sort_order`, `created_by`) VALUES 
('系统基础表情', '系统内置的基础表情包', 1, '37395c61-a2ec-464e-9567-ce6fa92630f7');

-- 查看表结构
-- SHOW CREATE TABLE emojis;
-- SHOW CREATE TABLE emoji_groups;
-- SHOW CREATE TABLE emoji_sprites;
-- SHOW CREATE TABLE emoji_tasks;
