# 项目结构说明

## 📋 概述

本项目采用**微服务架构**和**Go标准项目结构**，遵循Go生态系统的最佳实践。项目由4个服务组成：

- **server-blog**: 博客后端服务
- **server-auth-service**: SSO认证服务
- **web-blog**: 博客前端应用
- **web-auth-service**: SSO登录前端

项目结构清晰，职责分离明确，便于维护和扩展。

## 📁 完整目录结构

```
goBlog/
├── server-blog/                    # 博客后端服务
│   ├── internal/                   # 私有应用程序和库代码
│   │   ├── api/                   # API层 - HTTP处理器 (15个文件)
│   │   │   ├── ad.go             # 广告管理API
│   │   │   ├── ai_chat.go        # AI聊天API
│   │   │   ├── article.go        # 文章管理API
│   │   │   ├── category.go       # 分类管理API
│   │   │   ├── comment.go        # 评论管理API
│   │   │   ├── emoji.go          # Emoji表情API
│   │   │   ├── link.go           # 友情链接API
│   │   │   ├── message.go        # 留言板API
│   │   │   ├── photo.go          # 相册管理API
│   │   │   ├── tag.go            # 标签管理API
│   │   │   ├── talk.go           # 说说管理API
│   │   │   ├── upload.go         # 文件上传API
│   │   │   ├── user.go           # 用户管理API
│   │   │   └── website.go        # 网站配置API
│   │   ├── service/               # 业务逻辑层 (27个文件)
│   │   │   ├── ad_service.go
│   │   │   ├── ai_chat_service.go
│   │   │   ├── article_service.go
│   │   │   ├── category_service.go
│   │   │   ├── comment_service.go
│   │   │   ├── emoji_service.go
│   │   │   └── ...               # 其他业务逻辑
│   │   ├── model/                 # 数据模型 (47个文件)
│   │   │   ├── database/         # 数据库模型 (GORM)
│   │   │   ├── request/          # 请求参数模型
│   │   │   ├── response/         # 响应数据模型
│   │   │   ├── appTypes/         # 应用类型定义
│   │   │   ├── elasticsearch/    # ES相关模型
│   │   │   └── other/            # 其他模型
│   │   ├── middleware/            # 中间件 (5个文件)
│   │   │   ├── jwt.go            # JWT认证中间件
│   │   │   ├── logger.go         # 日志记录中间件
│   │   │   ├── cors.go           # 跨域处理中间件
│   │   │   └── ...
│   │   ├── router/                # 路由配置 (16个文件)
│   │   │   ├── router.go         # 主路由
│   │   │   ├── article.go        # 文章路由
│   │   │   ├── user.go           # 用户路由
│   │   │   └── ...
│   │   ├── initialize/            # 初始化模块 (8个文件)
│   │   │   ├── gorm.go           # 数据库初始化
│   │   │   ├── redis.go          # Redis初始化
│   │   │   ├── elasticsearch.go  # ES初始化
│   │   │   ├── cron.go           # 定时任务初始化
│   │   │   └── ...
│   │   └── task/                  # 定时任务 (5个文件)
│   │       ├── article_task.go   # 文章相关任务
│   │       ├── cache_task.go     # 缓存清理任务
│   │       └── ...
│   ├── pkg/                       # 可重用库代码
│   │   ├── config/               # 配置管理
│   │   │   └── config.go        # 配置结构定义
│   │   ├── utils/                # 工具函数
│   │   │   ├── jwt.go           # JWT工具
│   │   │   ├── md5.go           # MD5加密
│   │   │   ├── upload.go        # 上传工具
│   │   │   └── ...
│   │   ├── global/               # 全局变量
│   │   │   └── global.go        # 全局配置、DB、Redis等
│   │   └── core/                 # 核心功能
│   │       ├── server.go        # 服务器启动
│   │       ├── logger.go        # 日志配置
│   │       └── config.go        # 配置加载
│   ├── configs/                   # 配置文件
│   │   ├── config.yaml           # 开发环境配置
│   │   └── config.prod.yaml     # 生产环境配置
│   ├── scripts/                   # 脚本文件
│   │   └── flag/                 # 命令行工具
│   ├── keys/                      # RSA密钥文件
│   │   └── public.pem            # SSO公钥
│   ├── log/                       # 日志文件 (运行时生成)
│   ├── uploads/                   # 上传文件 (运行时生成)
│   ├── main.go                    # 程序入口
│   ├── go.mod                     # Go模块定义
│   ├── go.sum                     # 依赖版本锁定
│   ├── Dockerfile                 # Docker构建文件
│   └── .dockerignore             # Docker忽略文件
│
├── server-auth-service/           # SSO认证服务
│   ├── internal/                  # 私有应用程序和库代码
│   │   ├── handler/              # 请求处理器 (3个文件)
│   │   │   ├── auth_handler.go  # 认证处理
│   │   │   ├── user_handler.go  # 用户处理
│   │   │   └── oauth_handler.go # OAuth处理
│   │   ├── service/              # 业务逻辑层 (4个文件)
│   │   │   ├── auth_service.go
│   │   │   ├── user_service.go
│   │   │   ├── oauth_service.go
│   │   │   └── device_service.go
│   │   ├── model/                # 数据模型 (10个文件)
│   │   │   ├── entity/          # 数据库实体
│   │   │   │   ├── user.go
│   │   │   │   ├── oauth.go
│   │   │   │   ├── application.go
│   │   │   │   ├── device.go
│   │   │   │   └── login_log.go
│   │   │   ├── request/         # 请求模型
│   │   │   └── response/        # 响应模型
│   │   ├── router/               # 路由配置 (1个文件)
│   │   │   └── router.go
│   │   └── initialize/           # 初始化模块 (2个文件)
│   │       ├── gorm.go
│   │       └── redis.go
│   ├── pkg/                      # 公共库
│   │   ├── jwt/                 # JWT工具
│   │   │   ├── jwt.go
│   │   │   └── rsa.go          # RSA密钥加载
│   │   ├── crypto/              # 加密工具
│   │   │   └── password.go     # 密码加密
│   │   ├── config/              # 配置管理
│   │   ├── core/                # 核心功能
│   │   └── global/              # 全局变量
│   ├── configs/                  # 配置文件
│   │   ├── config.yaml
│   │   └── config.prod.yaml
│   ├── scripts/                  # 脚本文件
│   │   └── generate_keys.sh    # RSA密钥生成脚本
│   ├── keys/                     # RSA密钥文件
│   │   ├── private.pem          # 私钥
│   │   └── public.pem           # 公钥
│   ├── log/                      # 日志文件
│   ├── main.go                   # 程序入口
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── .dockerignore
│
├── web-blog/                      # 博客前端应用 (Vue3 + TypeScript)
│   ├── src/
│   │   ├── components/           # Vue组件 (37个)
│   │   │   ├── layout/          # 布局组件
│   │   │   │   ├── Header.vue
│   │   │   │   ├── Footer.vue
│   │   │   │   └── Sidebar.vue
│   │   │   ├── article/         # 文章组件
│   │   │   │   ├── ArticleCard.vue
│   │   │   │   ├── ArticleList.vue
│   │   │   │   └── ArticleDetail.vue
│   │   │   ├── comment/         # 评论组件
│   │   │   ├── editor/          # 编辑器组件
│   │   │   └── ...
│   │   ├── views/                # 页面视图 (40个)
│   │   │   ├── Home.vue         # 首页
│   │   │   ├── Article/         # 文章相关页面
│   │   │   ├── User/            # 用户相关页面
│   │   │   ├── Admin/           # 后台管理页面
│   │   │   └── ...
│   │   ├── api/                  # API接口 (14个)
│   │   │   ├── article.ts
│   │   │   ├── user.ts
│   │   │   ├── comment.ts
│   │   │   ├── ai-chat.ts
│   │   │   └── ...
│   │   ├── stores/               # Pinia状态管理 (5个)
│   │   │   ├── user.ts          # 用户状态
│   │   │   ├── article.ts       # 文章状态
│   │   │   ├── app.ts           # 应用状态
│   │   │   └── index.ts
│   │   ├── router/               # 路由配置 (1个)
│   │   │   └── index.ts
│   │   ├── utils/                # 工具函数 (6个)
│   │   │   ├── request.ts       # Axios封装
│   │   │   ├── auth.ts          # 认证工具
│   │   │   ├── date.ts          # 日期工具
│   │   │   └── ...
│   │   ├── assets/               # 静态资源 (2个)
│   │   │   ├── base.css
│   │   │   └── logo.svg
│   │   ├── App.vue               # 根组件
│   │   └── main.ts               # 入口文件
│   ├── public/                   # 公共资源
│   │   └── image_mapping.json   # 图片映射
│   ├── package.json              # 项目配置
│   ├── vite.config.ts            # Vite配置
│   ├── tsconfig.json             # TypeScript配置
│   ├── index.html                # HTML模板
│   ├── Dockerfile                # Docker构建文件
│   ├── nginx.conf                # Nginx配置
│   ├── .env                      # 环境变量
│   └── .dockerignore
│
├── web-auth-service/              # SSO登录前端 (Vue3)
│   ├── src/
│   │   ├── components/           # Vue组件
│   │   │   ├── LoginForm.vue
│   │   │   ├── RegisterForm.vue
│   │   │   └── OAuthButtons.vue
│   │   ├── views/                # 页面视图
│   │   │   ├── Login.vue
│   │   │   ├── Register.vue
│   │   │   └── Callback.vue
│   │   ├── api/                  # API接口
│   │   │   └── auth.js
│   │   ├── utils/                # 工具函数
│   │   ├── stores/               # 状态管理
│   │   ├── router/               # 路由配置
│   │   ├── App.vue
│   │   └── main.js
│   ├── public/
│   ├── package.json
│   ├── vite.config.js
│   ├── index.html
│   ├── Dockerfile
│   ├── nginx.conf
│   └── .dockerignore
│
├── nginx/                         # Nginx反向代理
│   ├── Dockerfile                # Nginx镜像构建
│   └── go_blog.conf              # Nginx配置文件
│       ├── www.hsk423.cn        # 博客域名配置
│       └── sso.hsk423.cn        # SSO域名配置
│
├── scripts/                       # 工具脚本
│   └── git-sync/                 # Git双仓库同步工具
│       ├── README.md             # 使用说明
│       └── sync_repos.sh         # 同步脚本
│
├── docs/                          # 项目文档
│   └── PROJECT_STRUCTURE.md      # 本文档
│
├── .workflow/                     # CI/CD工作流
│   ├── master-pipeline.yml       # 主分支流水线
│   ├── branch-pipeline.yml       # 分支流水线
│   └── pr-pipeline.yml           # PR流水线
│
├── docker-compose.base.yml        # 基础服务配置
│   ├── MySQL 8.0                 # 数据库服务
│   ├── Redis 6.2                 # 缓存服务
│   └── Elasticsearch 8.17        # 搜索引擎服务
│
├── docker-compose.yml             # 开发环境配置
│   ├── server-blog
│   ├── server-auth
│   ├── web-blog
│   ├── web-auth
│   └── nginx
│
├── docker-compose.prod.yml        # 生产环境配置
│
├── deploy.sh                      # 自动化部署脚本
├── README.md                      # 项目说明文档
├── README.en.md                   # 英文说明文档
├── .gitignore                     # Git忽略文件
└── dirMaker.bat                   # 目录创建脚本
```

## 📖 详细说明

### 一、后端服务 (Go)

#### 1. server-blog (博客后端服务)

**服务职责**: 提供博客核心功能的API服务

**主要模块**:

##### internal/api/ - API处理器层
- **职责**: 处理HTTP请求和响应
- **功能**: 
  - 参数验证和绑定
  - 调用service层处理业务逻辑
  - 返回标准JSON响应
  - 错误处理和异常捕获
- **主要API**:
  - `article.go`: 文章增删改查、点赞、收藏
  - `ai_chat.go`: AI聊天对话、流式响应
  - `comment.go`: 评论管理、回复功能
  - `user.go`: 用户信息管理
  - `upload.go`: 文件上传（七牛云）
  - `emoji.go`: Emoji表情管理

##### internal/service/ - 业务逻辑层
- **职责**: 实现具体的业务逻辑
- **功能**:
  - 业务规则验证
  - 数据库操作（通过GORM）
  - Redis缓存管理
  - Elasticsearch搜索
  - 第三方服务调用
- **设计模式**: 面向接口编程，便于单元测试

##### internal/model/ - 数据模型层
- **database/**: GORM数据库模型，对应数据库表结构
- **request/**: API请求参数结构体，用于参数绑定和验证
- **response/**: API响应数据结构体，统一响应格式
- **appTypes/**: 应用类型定义（枚举、常量等）
- **elasticsearch/**: ES文档模型

##### internal/middleware/ - 中间件
- **jwt.go**: JWT令牌验证，支持SSO统一认证
- **logger.go**: 请求日志记录（请求路径、耗时、状态码）
- **cors.go**: 跨域处理
- **rate_limit.go**: 访问频率限制
- **login_record.go**: 登录日志记录

##### internal/router/ - 路由配置
- **职责**: 定义API路由和中间件链
- **特点**: 
  - 路由分组管理
  - 中间件灵活组合
  - 支持路由版本控制

##### internal/initialize/ - 初始化模块
- **gorm.go**: MySQL数据库连接池初始化
- **redis.go**: Redis连接初始化
- **elasticsearch.go**: ES客户端初始化
- **cron.go**: 定时任务调度器初始化
- **router.go**: 路由初始化

##### internal/task/ - 定时任务
- **article_task.go**: 文章相关定时任务（浏览量同步）
- **cache_task.go**: 缓存预热和清理
- **backup_task.go**: 数据备份任务

##### pkg/ - 公共库
- **config/**: 配置文件解析和管理
- **utils/**: 工具函数（JWT、MD5、时间处理等）
- **global/**: 全局变量（DB、Redis、ES、Config等）
- **core/**: 核心功能（服务器启动、日志初始化）

#### 2. server-auth-service (SSO认证服务)

**服务职责**: 提供统一的身份认证和授权服务

**主要模块**:

##### internal/handler/ - 请求处理器
- **auth_handler.go**: 登录、注册、登出、令牌刷新
- **user_handler.go**: 用户信息管理
- **oauth_handler.go**: 第三方OAuth登录（GitHub、Google等）

##### internal/service/ - 业务逻辑
- **auth_service.go**: 认证逻辑（密码验证、JWT生成）
- **user_service.go**: 用户管理
- **oauth_service.go**: OAuth流程处理
- **device_service.go**: 设备管理（多设备登录）

##### internal/model/entity/ - 数据库实体
- **user.go**: SSO用户表
- **oauth.go**: OAuth绑定表
- **application.go**: 应用注册表
- **device.go**: 设备信息表
- **login_log.go**: 登录日志表

##### pkg/jwt/ - JWT工具
- **jwt.go**: JWT生成和验证
- **rsa.go**: RSA密钥加载（非对称加密）

##### pkg/crypto/ - 加密工具
- **password.go**: 密码加密（bcrypt）

**安全特性**:
- RSA非对称加密（私钥签发，公钥验证）
- 密码bcrypt加密存储
- 令牌过期和刷新机制
- 设备指纹识别
- 登录日志审计

### 二、前端应用 (Vue3)

#### 1. web-blog (博客前端)

**技术栈**: Vue3 + TypeScript + Element Plus + Vite

**主要模块**:

##### src/components/ - Vue组件
- **layout/**: 布局组件（Header、Footer、Sidebar）
- **article/**: 文章组件（列表、卡片、详情）
- **comment/**: 评论组件（列表、表单、回复）
- **editor/**: Markdown编辑器组件
- **common/**: 通用组件（分页、加载、空状态）

##### src/views/ - 页面视图
- **Home.vue**: 首页（文章列表、推荐文章）
- **Article/**: 文章相关页面（详情、编辑、发布）
- **User/**: 用户相关页面（个人中心、设置）
- **Admin/**: 后台管理页面（文章管理、评论管理）
- **AIChat/**: AI聊天页面

##### src/api/ - API接口
- 使用Axios封装的HTTP请求
- 统一的请求拦截器（添加Token）
- 统一的响应拦截器（错误处理）
- 接口按模块分文件管理

##### src/stores/ - Pinia状态管理
- **user.ts**: 用户状态（登录信息、权限）
- **article.ts**: 文章状态（列表、详情）
- **app.ts**: 应用状态（主题、语言）

##### src/router/ - 路由配置
- 路由懒加载
- 路由守卫（权限验证）
- 动态路由（根据权限生成）

##### src/utils/ - 工具函数
- **request.ts**: Axios封装
- **auth.ts**: 认证工具（Token存储、SSO跳转）
- **date.ts**: 日期格式化
- **validate.ts**: 表单验证

#### 2. web-auth-service (SSO登录前端)

**技术栈**: Vue3 + Vite

**主要功能**:
- 统一登录页面
- 注册页面
- 第三方OAuth登录
- 登录后回调处理
- Token传递给业务应用

### 三、基础设施

#### 1. Nginx反向代理

**配置文件**: `nginx/go_blog.conf`

**功能**:
- **域名路由**: 
  - `www.hsk423.cn` → 博客服务
  - `sso.hsk423.cn` → SSO服务
- **反向代理**: 
  - `/api/` → 后端服务
  - `/` → 前端静态文件
- **SSL/TLS**: HTTPS支持
- **Gzip压缩**: 静态资源压缩
- **流式响应**: AI聊天流式数据支持
- **安全头**: HSTS、X-Frame-Options等

#### 2. Docker容器化

**基础服务** (`docker-compose.base.yml`):
- **MySQL 8.0**: 数据库服务，数据持久化
- **Redis 6.2**: 缓存服务，配置持久化
- **Elasticsearch 8.17**: 搜索引擎服务，数据持久化

**业务服务** (`docker-compose.yml` / `docker-compose.prod.yml`):
- **server-blog**: 博客后端容器
- **server-auth**: SSO认证容器
- **web-blog**: 博客前端容器（Nginx）
- **web-auth**: SSO前端容器（Nginx）
- **nginx**: 反向代理容器

**网络**: 所有服务在同一个Docker网络 `blog-network`，容器间通过服务名通信

#### 3. 自动化部署

**deploy.sh脚本功能**:
1. **构建镜像**: 本地构建所有服务的Docker镜像
2. **保存上传**: 将镜像保存为tar文件并上传到服务器
3. **远程部署**: 在服务器加载镜像并启动容器

**支持的部署模式**:
- 全量部署：`./deploy.sh all`
- 分步部署：`./deploy.sh build/upload/deploy`
- 单服务部署：`./deploy.sh single server-blog all`

## 🏗️ 架构设计原则

### 1. 分层架构

```
┌─────────────────────────────────┐
│     Presentation Layer          │  API层（HTTP处理）
├─────────────────────────────────┤
│     Business Logic Layer        │  Service层（业务逻辑）
├─────────────────────────────────┤
│     Data Access Layer           │  Model层（数据访问）
├─────────────────────────────────┤
│     Infrastructure Layer        │  基础设施（DB、Redis、ES）
└─────────────────────────────────┘
```

### 2. 微服务架构

- **服务拆分**: 博客服务和认证服务独立部署
- **服务通信**: HTTP/HTTPS RESTful API
- **数据隔离**: 各服务独立数据库（逻辑隔离）
- **统一网关**: Nginx作为API网关

### 3. 设计原则

1. **单一职责原则 (SRP)**: 每个模块只负责一个功能
2. **开闭原则 (OCP)**: 对扩展开放，对修改关闭
3. **依赖倒置原则 (DIP)**: 依赖抽象而非具体实现
4. **接口隔离原则 (ISP)**: 使用专门的接口
5. **关注点分离**: API、Service、Model职责清晰

### 4. Go项目结构规范

- **internal/**: 私有代码，Go编译器强制隔离
- **pkg/**: 公共库，可被外部导入
- **cmd/**: 应用程序入口（本项目main.go在根目录）
- **configs/**: 配置文件
- **scripts/**: 工具脚本

## 🔧 技术实现细节

### 1. SSO单点登录流程

```
1. 用户访问博客 → 检测未登录 → 跳转SSO登录页
2. 用户在SSO登录 → SSO验证成功 → 生成JWT（RSA私钥签名）
3. SSO回调博客 → 传递JWT → 博客验证JWT（RSA公钥验证）
4. 验证成功 → 博客生成自己的Session → 用户登录成功
```

### 2. JWT认证机制

- **SSO服务**: 使用RSA私钥签发JWT
- **博客服务**: 使用SSO公钥验证JWT
- **优势**: 
  - 无需共享密钥
  - 安全性高
  - 支持跨域认证

### 3. 缓存策略

**Redis缓存层级**:
1. **热点数据**: 文章详情、用户信息（TTL: 1小时）
2. **列表数据**: 文章列表、分类列表（TTL: 10分钟）
3. **统计数据**: 浏览量、点赞数（实时更新，定时同步到DB）

**缓存更新策略**:
- **Cache Aside**: 先更新DB，再删除缓存
- **定时同步**: 统计数据定时批量写入DB

### 4. 全文搜索

**Elasticsearch索引设计**:
- **文章索引**: 标题、内容、标签、分类
- **分词器**: IK中文分词
- **搜索功能**: 
  - 关键词搜索
  - 高亮显示
  - 相关度排序

### 5. AI聊天实现

**技术方案**:
- **API**: DeepSeek AI API
- **流式响应**: Server-Sent Events (SSE)
- **上下文管理**: Redis存储会话历史
- **Nginx配置**: 禁用缓冲，支持流式传输

## 📝 开发规范

### 1. 代码规范

- **Go**: 遵循 Go Code Review Comments
- **Vue**: 遵循 Vue.js Style Guide
- **命名**: 
  - Go: 驼峰命名（CamelCase）
  - Vue: 短横线命名（kebab-case）
  - 数据库: 下划线命名（snake_case）

### 2. Git规范

**分支管理**:
- `master`: 生产环境分支
- `develop`: 开发分支
- `feature/*`: 功能分支
- `hotfix/*`: 紧急修复分支

**提交规范**:
```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式调整
refactor: 重构
test: 测试相关
chore: 构建/工具相关
```

### 3. API规范

**RESTful设计**:
- `GET /api/articles`: 获取文章列表
- `GET /api/articles/:id`: 获取文章详情
- `POST /api/articles`: 创建文章
- `PUT /api/articles/:id`: 更新文章
- `DELETE /api/articles/:id`: 删除文章

**响应格式**:
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

### 4. 错误处理

- **统一错误码**: 定义在 `appTypes/` 中
- **错误日志**: 记录到日志文件
- **用户提示**: 友好的错误信息

## 🚀 扩展指南

### 添加新功能的步骤

1. **确定功能模块**: 属于哪个服务（博客/认证）
2. **设计数据模型**: 在 `model/database/` 添加数据库模型
3. **实现业务逻辑**: 在 `service/` 添加业务逻辑
4. **添加API接口**: 在 `api/` 添加HTTP处理器
5. **配置路由**: 在 `router/` 添加路由
6. **前端开发**: 
   - 在 `api/` 添加接口调用
   - 在 `views/` 添加页面
   - 在 `components/` 添加组件
7. **测试**: 编写单元测试和集成测试
8. **文档**: 更新API文档和使用文档

### 性能优化建议

1. **数据库优化**: 
   - 添加索引
   - 查询优化
   - 读写分离
2. **缓存优化**: 
   - 增加缓存层级
   - 优化缓存策略
3. **前端优化**: 
   - 路由懒加载
   - 组件按需加载
   - 图片懒加载
4. **服务器优化**: 
   - 负载均衡
   - CDN加速
   - Gzip压缩

## 📚 参考资料

- [Go标准项目布局](https://github.com/golang-standards/project-layout)
- [Vue.js官方文档](https://vuejs.org/)
- [Element Plus文档](https://element-plus.org/)
- [Docker官方文档](https://docs.docker.com/)
- [Nginx官方文档](https://nginx.org/en/docs/) 