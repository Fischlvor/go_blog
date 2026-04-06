-- =============================================
-- go_blog PostgreSQL 数据库初始化脚本
-- 最终状态（已合并所有迁移）
-- =============================================

-- 启用 uuid-ossp 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==================== 用户表 ====================
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    nickname VARCHAR(50),
    avatar VARCHAR(500),
    email VARCHAR(100),
    role_id INT DEFAULT 1,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_uuid ON users(uuid);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- ==================== 文章分类表 ====================
CREATE TABLE IF NOT EXISTS article_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) UNIQUE,
    article_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_article_categories_deleted_at ON article_categories(deleted_at);

-- ==================== 文章标签表 ====================
CREATE TABLE IF NOT EXISTS article_tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) UNIQUE,
    article_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_article_tags_deleted_at ON article_tags(deleted_at);

-- ==================== 文章表 ====================
CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL UNIQUE,
    excerpt TEXT,
    content TEXT NOT NULL,
    featured_image VARCHAR(500),
    author_uuid UUID,
    category_id BIGINT NOT NULL REFERENCES article_categories(id),
    status VARCHAR(20) DEFAULT 'draft',
    visibility VARCHAR(20) NOT NULL DEFAULT 'public', -- 文章可见性: public(公开) | private(私有)
    tag_ids BIGINT[] DEFAULT '{}', -- 标签 ID 数组
    read_time VARCHAR(20),
    views INT DEFAULT 0,
    likes INT DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    meta_title VARCHAR(200),
    meta_description VARCHAR(500),
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_articles_author_uuid ON articles(author_uuid);
CREATE INDEX IF NOT EXISTS idx_articles_category_id ON articles(category_id);
CREATE INDEX IF NOT EXISTS idx_articles_status ON articles(status);
CREATE INDEX IF NOT EXISTS idx_articles_slug ON articles(slug);
CREATE INDEX IF NOT EXISTS idx_articles_deleted_at ON articles(deleted_at);
CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles(published_at);
CREATE INDEX IF NOT EXISTS idx_articles_tag_ids ON articles USING GIN (tag_ids); -- GIN 索引加速标签查询

-- ==================== 文章点赞表 ====================
CREATE TABLE IF NOT EXISTS article_likes (
    id BIGSERIAL PRIMARY KEY,
    article_slug VARCHAR(200) NOT NULL,
    user_uuid UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE (article_slug, user_uuid)
);

CREATE INDEX IF NOT EXISTS idx_article_likes_article_slug ON article_likes(article_slug);
CREATE INDEX IF NOT EXISTS idx_article_likes_user_uuid ON article_likes(user_uuid);
CREATE INDEX IF NOT EXISTS idx_article_likes_deleted_at ON article_likes(deleted_at);

-- ==================== 文章浏览记录表 ====================
CREATE TABLE IF NOT EXISTS article_views (
    id BIGSERIAL PRIMARY KEY,
    article_slug VARCHAR(200) NOT NULL,
    ip_address VARCHAR(50) NOT NULL,
    user_agent TEXT,
    referer TEXT,
    viewed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_article_views_article_slug ON article_views(article_slug);
CREATE INDEX IF NOT EXISTS idx_article_views_viewed_at ON article_views(viewed_at);

-- ==================== 评论表 ====================
CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    article_slug VARCHAR(200) NOT NULL,
    parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    user_uuid UUID NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'approved',
    likes INT DEFAULT 0,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_comments_article_slug ON comments(article_slug);
CREATE INDEX IF NOT EXISTS idx_comments_user_uuid ON comments(user_uuid);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_deleted_at ON comments(deleted_at);

-- ==================== 评论点赞表 ====================
CREATE TABLE IF NOT EXISTS comment_likes (
    id BIGSERIAL PRIMARY KEY,
    comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    user_uuid UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (comment_id, user_uuid)
);

CREATE INDEX IF NOT EXISTS idx_comment_likes_comment_id ON comment_likes(comment_id);
CREATE INDEX IF NOT EXISTS idx_comment_likes_user_uuid ON comment_likes(user_uuid);

-- ==================== AI 模型配置表 ====================
CREATE TABLE IF NOT EXISTS ai_models (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    endpoint VARCHAR(255),
    api_key VARCHAR(255),
    max_tokens INTEGER DEFAULT 4096,
    temperature DOUBLE PRECISION DEFAULT 0.7,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_ai_models_provider ON ai_models(provider);
CREATE INDEX IF NOT EXISTS idx_ai_models_is_active ON ai_models(is_active);
CREATE INDEX IF NOT EXISTS idx_ai_models_deleted_at ON ai_models(deleted_at);

-- ==================== AI 聊天会话表 ====================
CREATE TABLE IF NOT EXISTS ai_chat_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_uuid UUID NOT NULL,
    title VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_ai_chat_sessions_user_uuid ON ai_chat_sessions(user_uuid);
CREATE INDEX IF NOT EXISTS idx_ai_chat_sessions_deleted_at ON ai_chat_sessions(deleted_at);

-- ==================== AI 聊天消息表 ====================
CREATE TABLE IF NOT EXISTS ai_chat_messages (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL REFERENCES ai_chat_sessions(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    tokens INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_ai_chat_messages_session_id ON ai_chat_messages(session_id);
CREATE INDEX IF NOT EXISTS idx_ai_chat_messages_deleted_at ON ai_chat_messages(deleted_at);

-- ==================== 反馈表 ====================
CREATE TABLE IF NOT EXISTS feedbacks (
    id BIGSERIAL PRIMARY KEY,
    user_uuid UUID,
    type VARCHAR(20),
    content TEXT NOT NULL,
    contact VARCHAR(100),
    reply TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_feedbacks_user_uuid ON feedbacks(user_uuid);
CREATE INDEX IF NOT EXISTS idx_feedbacks_deleted_at ON feedbacks(deleted_at);

-- ==================== 友链表 ====================
CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(500) NOT NULL,
    logo VARCHAR(500),
    description VARCHAR(200),
    sort INT DEFAULT 0,
    is_visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_links_deleted_at ON links(deleted_at);

-- ==================== 文件表 ====================
CREATE TABLE IF NOT EXISTS files (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(500) NOT NULL UNIQUE,
    filename VARCHAR(255),
    file_hash VARCHAR(64),
    size BIGINT,
    mime_type VARCHAR(100),
    usage VARCHAR(50),
    resource_id BIGINT,
    user_uuid UUID,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_files_key ON files(key);
CREATE INDEX IF NOT EXISTS idx_files_file_hash ON files(file_hash);
CREATE INDEX IF NOT EXISTS idx_files_user_uuid ON files(user_uuid);
CREATE INDEX IF NOT EXISTS idx_files_deleted_at ON files(deleted_at);

-- ==================== 广告表 ====================
CREATE TABLE IF NOT EXISTS advertisements (
    id BIGSERIAL PRIMARY KEY,
    ad_name VARCHAR(100) NOT NULL,
    ad_image VARCHAR(500),
    ad_link VARCHAR(500),
    ad_type INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ==================== 页脚链接表 ====================
CREATE TABLE IF NOT EXISTS footer_links (
    title VARCHAR(100) PRIMARY KEY,
    link VARCHAR(500) NOT NULL
);

-- ==================== 表情组表 ====================
CREATE TABLE IF NOT EXISTS emoji_groups (
    id BIGSERIAL PRIMARY KEY,
    group_name VARCHAR(100) NOT NULL UNIQUE,
    group_key VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(500) DEFAULT '',
    sort_order INT DEFAULT 0,
    emoji_count INT DEFAULT 0,
    status SMALLINT DEFAULT 1,
    sprite_conf_url VARCHAR(500) DEFAULT '',
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_emoji_groups_status ON emoji_groups(status);
CREATE INDEX IF NOT EXISTS idx_emoji_groups_sort_order ON emoji_groups(sort_order);

-- ==================== 表情表 ====================
CREATE TABLE IF NOT EXISTS emojis (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(50) NOT NULL UNIQUE,
    filename VARCHAR(255) NOT NULL,
    group_key VARCHAR(50) NOT NULL,
    sprite_group INT NOT NULL DEFAULT 0,
    sprite_position_x INT NOT NULL DEFAULT 0,
    sprite_position_y INT NOT NULL DEFAULT 0,
    file_size INT DEFAULT 0,
    cdn_url VARCHAR(500) DEFAULT '',
    upload_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    status SMALLINT DEFAULT 1,
    created_by UUID,
    deleted_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_emojis_group_key ON emojis(group_key);
CREATE INDEX IF NOT EXISTS idx_emojis_status ON emojis(status);
CREATE INDEX IF NOT EXISTS idx_emojis_sprite_group ON emojis(sprite_group);

-- ==================== 雪碧图表 ====================
CREATE TABLE IF NOT EXISTS emoji_sprites (
    id BIGSERIAL PRIMARY KEY,
    sprite_group INT NOT NULL UNIQUE,
    filename VARCHAR(255) NOT NULL,
    cdn_url VARCHAR(500) DEFAULT '',
    width INT NOT NULL,
    height INT NOT NULL,
    emoji_count INT NOT NULL,
    file_size INT DEFAULT 0,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_emoji_sprites_status ON emoji_sprites(status);

-- ==================== 表情任务表 ====================
CREATE TABLE IF NOT EXISTS emoji_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    progress INT DEFAULT 0,
    message VARCHAR(500) DEFAULT '',
    params JSONB,
    result JSONB,
    created_by UUID,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_emoji_tasks_status ON emoji_tasks(status);
CREATE INDEX IF NOT EXISTS idx_emoji_tasks_task_type ON emoji_tasks(task_type);

-- ==================== 资源表 ====================
CREATE TABLE IF NOT EXISTS resources (
    id BIGSERIAL PRIMARY KEY,
    file_key VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_hash VARCHAR(64),
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    user_uuid UUID,
    transcode_status SMALLINT DEFAULT 0,
    transcode_key VARCHAR(500),
    thumbnail_key VARCHAR(500),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_resources_file_hash ON resources(file_hash);
CREATE INDEX IF NOT EXISTS idx_resources_user_uuid ON resources(user_uuid);
CREATE INDEX IF NOT EXISTS idx_resources_deleted_at ON resources(deleted_at);

-- ==================== 资源上传任务表 ====================
CREATE TABLE IF NOT EXISTS resource_upload_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_id VARCHAR(64) NOT NULL UNIQUE,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_hash VARCHAR(64),
    mime_type VARCHAR(100) NOT NULL,
    chunk_size INT DEFAULT 4194304,
    total_chunks INT NOT NULL,
    status SMALLINT DEFAULT 0,
    user_uuid UUID,
    expires_at TIMESTAMPTZ,
    qiniu_contexts JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_resource_upload_tasks_task_id ON resource_upload_tasks(task_id);
CREATE INDEX IF NOT EXISTS idx_resource_upload_tasks_user_uuid ON resource_upload_tasks(user_uuid);
CREATE INDEX IF NOT EXISTS idx_resource_upload_tasks_expires_at ON resource_upload_tasks(expires_at);

-- ==================== 登录记录表 ====================
CREATE TABLE IF NOT EXISTS logins (
    id BIGSERIAL PRIMARY KEY,
    user_uuid UUID,
    login_method VARCHAR(50),
    ip VARCHAR(50),
    address VARCHAR(200),
    os VARCHAR(100),
    device_info VARCHAR(200),
    browser_info VARCHAR(200),
    status INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_logins_user_uuid ON logins(user_uuid);
CREATE INDEX IF NOT EXISTS idx_logins_created_at ON logins(created_at);
CREATE INDEX IF NOT EXISTS idx_logins_deleted_at ON logins(deleted_at);

-- ==================== 站点配置表 ====================
CREATE TABLE IF NOT EXISTS site_settings (
    id BIGSERIAL PRIMARY KEY,
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT NOT NULL,
    setting_type VARCHAR(20) NOT NULL DEFAULT 'string',
    description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_site_settings_is_public ON site_settings(is_public);

INSERT INTO site_settings (setting_key, setting_value, setting_type, description, is_public)
VALUES
    ('profile.avatar', '', 'string', '头像', TRUE),
    ('website.title', '', 'string', '网站标题', TRUE),
    ('website.description', '', 'string', '网站描述', TRUE),
    ('profile.intro', '', 'string', '个人介绍', TRUE),
    ('website.version', '', 'string', '网站版本', TRUE),
    ('website.created_at', '', 'string', '网站创建日期', TRUE),
    ('website.icp_filing', '', 'string', 'ICP备案号', TRUE),
    ('profile.bilibili_url', '', 'string', 'Bilibili 链接', TRUE),
    ('profile.github_url', '', 'string', 'GitHub 链接', TRUE),
    ('website.steam_url', '', 'string', 'Steam 链接', TRUE),
    ('profile.name', '', 'string', '名称', TRUE),
    ('profile.job', '', 'string', '职业', TRUE),
    ('profile.address', '', 'string', '地址', TRUE),
    ('profile.email', '', 'string', '联系邮箱', TRUE)
ON CONFLICT (setting_key) DO NOTHING;
