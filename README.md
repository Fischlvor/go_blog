# goBlog

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.23-blue)
![Vue Version](https://img.shields.io/badge/Vue-3.5-green)
![License](https://img.shields.io/badge/license-MIT-orange)

一个基于 Go + Vue3 的现代化全栈博客系统，集成 SSO 单点登录和 AI 聊天助手

[在线演示](https://www.hsk423.cn) | [项目文档](./docs/PROJECT_STRUCTURE.md) | [更新日志](#)

</div>

---

## 📖 项目介绍

goBlog 是一个功能完善的现代化博客系统，采用**前后端分离架构**和**微服务设计**，提供完整的博客管理功能。项目遵循 Go 标准项目结构，代码组织清晰，易于维护和扩展。

### ✨ 核心特性

- 🔐 **SSO 单点登录** - 独立的认证服务，支持多应用统一登录
- 🤖 **AI 聊天助手** - 集成 DeepSeek AI，提供智能对话功能
- 🚀 **高性能架构** - Redis 缓存 + Elasticsearch 全文搜索
- 📝 **Markdown 编辑器** - 支持实时预览和富文本编辑
- 🎨 **现代化 UI** - Vue3 + TypeScript + Element Plus
- 🐳 **Docker 部署** - 完整的容器化部署方案
- 🔄 **CI/CD 支持** - 自动化构建和部署脚本

## 🛠️ 技术栈

### 后端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.23 | 主要开发语言 |
| Gin | 1.10 | Web 框架 |
| GORM | 1.25 | ORM 框架 |
| MySQL | 8.0 | 主数据库 |
| Redis | 6.2 | 缓存数据库 |
| Elasticsearch | 8.17 | 全文搜索引擎 |
| JWT | - | 身份认证 |
| Zap | 1.27 | 日志框架 |
| go-ratelimiter | 0.3.0 | 限流中间件 |

### 前端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.5 | 前端框架 |
| TypeScript | 5.8 | 类型系统 |
| Element Plus | 2.10 | UI 组件库 |
| Vite | 6.2 | 构建工具 |
| Pinia | 3.0 | 状态管理 |
| Vue Router | 4.5 | 路由管理 |
| Axios | 1.11 | HTTP 客户端 |
| md-editor-v3 | 5.8 | Markdown 编辑器 |

### 基础设施

- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **SSL**: HTTPS 支持
- **限流保护**: 多层限流架构
- **部署**: 自动化部署脚本

## 📁 项目结构

```
goBlog/
├── server-blog/              # 博客后端服务
│   ├── internal/             # 私有代码
│   │   ├── api/             # API 处理器
│   │   ├── service/         # 业务逻辑
│   │   ├── model/           # 数据模型
│   │   ├── middleware/      # 中间件
│   │   ├── router/          # 路由配置
│   │   ├── initialize/      # 初始化模块
│   │   └── task/            # 定时任务
│   ├── pkg/                 # 公共库
│   │   ├── config/          # 配置管理
│   │   ├── utils/           # 工具函数
│   │   ├── global/          # 全局变量
│   │   └── core/            # 核心功能
│   ├── configs/             # 配置文件
│   ├── keys/                # RSA 密钥文件
│   ├── log/                 # 日志文件
│   ├── scripts/             # 工具脚本
│   ├── sql/                 # SQL 脚本
│   ├── uploads/             # 上传文件
│   └── main.go              # 程序入口
├── server-auth-service/      # SSO 认证服务
│   ├── internal/             # 私有代码
│   │   ├── api/             # API 处理器
│   │   │   ├── auth.go      # 认证 API
│   │   │   ├── captcha.go   # 验证码 API
│   │   │   ├── device.go    # 设备 API
│   │   │   ├── manage.go    # 管理 API
│   │   │   ├── oauth.go     # OAuth API
│   │   │   └── enter.go     # API 聚合
│   │   ├── service/         # 业务逻辑
│   │   │   ├── auth_service.go    # 认证服务
│   │   │   ├── code_service.go    # 授权码服务
│   │   │   ├── device_service.go  # 设备服务
│   │   │   ├── manage_service.go  # 管理服务
│   │   │   └── qq_service.go      # QQ 登录服务
│   │   ├── model/           # 数据模型
│   │   │   ├── appTypes/    # 应用类型枚举
│   │   │   ├── database/    # 数据库模型
│   │   │   ├── request/     # 请求体
│   │   │   ├── response/    # 响应体
│   │   │   └── other/       # 其他模型
│   │   ├── middleware/      # 中间件
│   │   │   ├── auth.go      # JWT 认证中间件
│   │   │   ├── client_auth.go # 客户端认证中间件
│   │   │   └── cors.go      # CORS 中间件
│   │   ├── router/          # 路由配置
│   │   │   ├── enter.go     # 路由聚合
│   │   │   ├── auth.go      # 认证路由
│   │   │   ├── base.go      # 基础路由
│   │   │   ├── device.go    # 设备路由
│   │   │   ├── manage.go    # 管理路由
│   │   │   ├── oauth.go     # OAuth 路由
│   │   │   ├── user.go      # 用户路由
│   │   │   └── router.go    # 路由设置
│   │   └── initialize/      # 初始化模块
│   ├── pkg/                 # 公共库
│   │   ├── jwt/             # JWT 工具
│   │   ├── crypto/          # 加密工具
│   │   ├── config/          # 配置管理
│   │   ├── utils/           # 工具函数
│   │   ├── global/          # 全局变量
│   │   └── core/            # 核心功能
│   ├── scripts/             # 工具脚本
│   │   └── flag/            # 命令行标志处理
│   ├── configs/             # 配置文件
│   ├── keys/                # RSA 密钥文件
│   ├── logs/                # 日志文件
│   ├── go.mod               # Go 模块文件
│   ├── go.sum               # 依赖锁定文件
│   └── main.go              # 程序入口
├── web-blog/                 # 博客前端 (Vue3 + TypeScript)
│   ├── src/
│   │   ├── components/      # Vue 组件
│   │   ├── views/           # 页面视图
│   │   ├── api/             # API 接口
│   │   ├── stores/          # Pinia 状态管理
│   │   ├── router/          # 路由配置
│   │   └── utils/           # 工具函数
│   └── package.json
├── web-auth-service/         # SSO 登录前端 (Vue3)
│   ├── src/
│   │   ├── components/      # Vue 组件
│   │   ├── views/           # 页面视图
│   │   └── api/             # API 接口
│   └── package.json
├── nginx/                    # Nginx 配置
│   ├── Dockerfile
│   └── go_blog.conf         # 反向代理配置
├── scripts/                  # 工具脚本
│   └── git-sync/            # Git 双仓库同步工具
├── docs/                     # 项目文档
│   └── PROJECT_STRUCTURE.md # 详细架构说明
├── docker-compose.base.yml   # 基础服务配置
├── docker-compose.yml        # 开发环境
├── docker-compose.prod.yml   # 生产环境
└── deploy.sh                 # 自动化部署脚本
```

详细架构说明请查看 [项目结构文档](./docs/PROJECT_STRUCTURE.md)

## 🚀 快速开始

### 环境要求

- Go 1.23+
- Node.js 18+
- MySQL 8.0+
- Redis 6.2+
- Elasticsearch 8.17+ (可选)
- Docker & Docker Compose (推荐)

### 开发环境部署

#### 1. 克隆项目

```bash
git clone https://gitee.com/qiyana423/go_blog.git
cd go_blog
```

#### 2. 启动基础服务

```bash
# 启动 MySQL、Redis、Elasticsearch
docker-compose -f docker-compose.base.yml up -d
```

#### 3. 配置后端服务

**博客服务配置:**

```bash
cd server-blog

# 复制配置文件模板
cp configs/config.example.yaml configs/config.yaml

# 编辑配置文件，修改数据库、Redis等配置
vim configs/config.yaml

# 生成 RSA 密钥对 (用于验证 SSO JWT)
mkdir -p keys
# 从 SSO 服务复制公钥到 keys/public.pem

# 安装依赖
go mod tidy

# 运行服务
go run main.go
```

**SSO 认证服务配置:**

```bash
cd server-auth-service

# 复制配置文件模板
cp configs/config.example.yaml configs/config.yaml

# 编辑配置文件，修改数据库、Redis、邮件等配置
vim configs/config.yaml

# 生成 RSA 密钥对
bash scripts/generate_keys.sh

# 安装依赖
go mod tidy

# 运行服务
go run main.go
```

#### 4. 配置前端服务

**博客前端:**

```bash
cd web-blog

# 安装依赖
npm install

# 开发模式运行
npm run dev
```

**SSO 登录前端:**

```bash
cd web-auth-service

# 安装依赖
npm install

# 开发模式运行
npm run dev
```

### 生产环境部署

#### 使用 Docker Compose 部署

```bash
# 1. 启动基础服务
docker-compose -f docker-compose.base.yml up -d

# 2. 构建并启动业务服务
docker-compose -f docker-compose.yml up -d
```

#### 使用自动化部署脚本

项目提供了完整的自动化部署脚本，支持分步执行：

```bash
# 完整部署流程
./deploy.sh all

# 或分步执行
./deploy.sh build    # 1. 构建镜像
./deploy.sh upload   # 2. 上传到服务器
./deploy.sh deploy   # 3. 远程部署

# 单服务部署
./deploy.sh single server-blog all    # 部署博客服务
./deploy.sh single web-blog build     # 仅构建前端
```

详细部署说明请查看 `deploy.sh` 脚本注释。

## 📚 主要功能

### 用户功能

- ✅ 用户注册、登录 (SSO 单点登录)
- ✅ 个人信息管理
- ✅ 密码修改、找回
- ✅ 第三方登录 (OAuth)
- ✅ 多设备管理

### 文章功能

- ✅ Markdown 编辑器
- ✅ 文章发布、编辑、删除
- ✅ 文章分类、标签管理
- ✅ 文章搜索 (Elasticsearch)
- ✅ 文章点赞、收藏
- ✅ 浏览量统计
- ✅ 文章置顶、推荐

### 评论功能

- ✅ 文章评论
- ✅ 评论回复
- ✅ 评论点赞
- ✅ 评论管理

### AI 功能

- ✅ AI 聊天助手
- ✅ 流式响应
- ✅ 上下文记忆
- ✅ 会话管理

### 内容管理

- ✅ 广告位管理
- ✅ 友情链接
- ✅ 网站配置
- ✅ 图片上传 (七牛云)
- ✅ Emoji 表情管理

### 系统功能

- ✅ 定时任务
- ✅ 日志记录
- ✅ 访问统计
- ✅ 健康检查
- ✅ 数据备份
- ✅ 多层限流保护
  - Nginx 全局限流（10000 QPS）
  - 应用层用户限流（1000 QPS/用户）
  - 业务规则限流（登录5次/分钟）
  - 自动拉黑机制

## 🔧 配置说明

### 后端配置

**博客服务配置**:
```bash
cd server-blog/configs
cp config.example.yaml config.yaml
# 编辑 config.yaml，修改数据库、Redis、七牛云等配置
```

**SSO认证服务配置**:
```bash
cd server-auth-service/configs
cp config.example.yaml config.yaml
# 编辑 config.yaml，修改数据库、Redis、邮件等配置
```

配置文件包含详细的注释说明，请根据实际情况修改。主要配置项包括：
- 数据库连接（MySQL）
- 缓存配置（Redis）
- 搜索引擎（Elasticsearch）
- 对象存储（七牛云）
- JWT密钥配置
- 邮件服务配置
- 第三方登录（QQ等）

### 前端环境变量

在前端项目根目录创建 `.env` 文件：

```bash
# 博客前端 (web-blog/.env)
VITE_API_BASE_URL=http://localhost:8080/api
VITE_SSO_URL=http://localhost:8000/api
VITE_CDN_URL=https://cdn.example.com

# SSO前端 (web-auth-service/.env)
VITE_API_BASE_URL=http://localhost:8000/api
```

## 📊 系统架构

### 服务架构

```
┌─────────────┐      ┌─────────────┐
│  web-blog   │      │  web-auth   │
│  (Vue3)     │      │  (Vue3)     │
└──────┬──────┘      └──────┬──────┘
       │                    │
       └────────┬───────────┘
                │
         ┌──────▼──────┐
         │    Nginx    │
         │  (反向代理)  │
         └──────┬──────┘
                │
       ┌────────┴────────┐
       │                 │
┌──────▼──────┐   ┌─────▼──────┐
│ server-blog │   │server-auth │
│   (Go)      │   │   (Go)     │
└──────┬──────┘   └─────┬──────┘
       │                │
       └────────┬───────┘
                │
    ┌───────────┼───────────┐
    │           │           │
┌───▼───┐  ┌───▼───┐  ┌───▼───┐
│ MySQL │  │ Redis │  │  ES   │
└───────┘  └───────┘  └───────┘
```

### 技术亮点

1. **微服务架构** - 博客服务和认证服务分离，职责清晰
2. **SSO 单点登录** - 统一的身份认证中心，支持多应用
3. **JWT + RSA** - 使用 RSA 非对称加密，安全性高
4. **多层限流保护** - Nginx + 应用层双重限流，支持用户级/IP级/规则级限流
5. **缓存策略** - Redis 多级缓存，提升性能
6. **全文搜索** - Elasticsearch 实现高效搜索
7. **容器化部署** - Docker 一键部署，环境隔离
8. **自动化运维** - CI/CD 脚本，简化部署流程

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 📝 开发日志

### 2025-11-30

#### Refactor
- **重构 server-auth-service 项目结构，完全对齐 server-blog**
  - 响应/请求结构体重构
    - 移动 `response.go` 到 `internal/model/response/`
    - 更新 `DeviceInfo`、`LogInfo` 到 `internal/model/response/`
    - 添加 `LogQueryParams` 到 `internal/model/request/`
    - 移动 `AuthorizationCode` 到 `internal/model/appTypes/`
    - 修复所有 package 声明（entity → database）
  - Router 结构重构
    - 创建 `enter.go` 聚合所有 Router
    - 按功能模块拆分 router（auth.go、base.go、oauth.go、user.go、device.go、manage.go）
    - 每个 router 使用 Group 组织相同前缀的路由
    - 更新 Setup 函数调用各个 router 的初始化方法
  - 更新 README 项目结构文档
  - 涉及文件：30+ 个文件，完全对齐 server-blog 的项目结构

### 2025-11-29

#### Fixed
- **修复日志和前后端联调** ([c586a4d](https://github.com/Fischlvor/go_blog/commit/c586a4d3864416a93a87d690b90e526d5a9f23b3))
  - 优化SSO认证服务的日志记录功能
  - 修复前后端数据联调问题
  - 完善设备管理界面的交互逻辑
  - 优化管理平台的用户体验
  - 涉及文件：11个文件，新增391行，删除83行

- **优化SSO静默登录，添加SSO管理平台，支持设备、全局下线** ([a62fdc1](https://github.com/Fischlvor/go_blog/commit/a62fdc1c63b64a9880b4192036f3a4eca510241f))
  - 实现SSO静默登录机制，提升用户体验
  - 新增完整的SSO管理平台界面
  - 实现设备管理功能（查看、下线设备）
  - 支持全局下线功能（一键下线所有设备）
  - 增加登录活动监控和安全管理
  - 优化OAuth认证流程和QQ登录
  - 完善前端路由和权限控制
  - 涉及文件：31个文件，新增2135行，删除145行

### 2025-11-23

#### Fixed
- **单点登录后端完成** ([25e3cbc](https://github.com/Fischlvor/go_blog/commit/25e3cbc6d33c27af670b0e6675c1c7897560b4e5))
  - 完善SSO认证服务后端核心功能
  - 实现完整的OAuth认证处理器
  - 增加认证路由和中间件
  - 完善认证服务业务逻辑
  - 涉及文件：6个文件，新增305行，删除11行

- **增加登出日志字段** ([9d7cba8](https://github.com/Fischlvor/go_blog/commit/9d7cba81b32eb88467b93017f0a2d24a1dd99539))
  - 优化用户登出功能的日志记录
  - 增加登出操作的详细日志字段
  - 完善认证服务的审计功能
  - 涉及文件：2个文件，新增14行，删除8行

### 2025-11-22

#### Fixed
- **blog退出登录走SSO** ([007fb09](https://github.com/Fischlvor/go_blog/commit/007fb0957b2c396d0945fef6ec9aaf8761b65cb3))
  - 重构博客系统的登出逻辑
  - 统一登出流程，通过SSO服务处理
  - 确保多应用间登出状态同步
  - 涉及文件：1个文件，新增62行，删除4行

### 2025-11-21

#### Added
- **集成多层限流保护系统**
  - 集成 go-ratelimiter v0.3.0 限流中间件
  - Nginx 层实现全局限流（10000 QPS）
  - 应用层实现用户级限流（1000 QPS/用户）
  - 实现业务规则限流（登录5次/分钟、注册3次/5分钟）
  - 支持设备ID识别和自动降级到IP限流
  - 实现违规记录和自动拉黑机制（15分/5分钟 触发30分钟封禁）
  - 前端复用设备ID生成逻辑，基于浏览器指纹
  - 限流信息通过标准 HTTP Header 返回（X-RateLimit-*）
  - 涉及文件：多个配置和中间件文件

### 2025-11-17

#### Fixed
- **优化OAuth state参数处理和验证码懒加载机制** ([f09ac05](https://github.com/Fischlvor/go_blog/commit/f09ac05a63b1ddd9d4af570216b4c8c0f4f5e384))
  - 优化SSO认证服务的OAuth state参数生成和验证逻辑
  - 增加Redis缓存支持，提升state验证性能
  - 实现验证码懒加载机制，优化用户体验
  - 修复登录、注册、找回密码页面的验证码加载问题
  - 优化前端路由跳转和SSO回调处理逻辑
  - 涉及文件：15个文件，新增518行，删除122行

- **优化emoji系统和封面URL处理逻辑** ([7317515](https://github.com/Fischlvor/go_blog/commit/73175153aca04aeb32edf7aaeb24a7cc0de64451))
  - 优化文章封面URL处理，支持相对路径和绝对路径
  - 重构emoji解析器，提升性能和可维护性
  - 增加emoji全局初始化模块，统一管理emoji资源
  - 优化emoji样式管理器，支持动态加载
  - 修复评论区emoji显示问题
  - 涉及文件：15个文件，新增363行，删除239行

- **移除AI聊天路由的JWT认证中间件并优化登录提示** ([492c37e](https://github.com/Fischlvor/go_blog/commit/492c37edebe7e31395489f5f2633b973d78e8463))
  - 移除AI聊天路由的JWT认证要求，允许游客访问
  - 优化前端路由守卫，提升用户体验
  - 涉及文件：2个文件，新增13行，删除11行

### 2025-11-16

#### Fixed
- **归档脚本** ([8f68c7b](https://github.com/Fischlvor/go_blog/commit/8f68c7b1faaf9126de299f32827ae008326fe9a0))
  - 清理已废弃的emoji处理脚本
  - 删除Python版本的emoji优化器和迁移工具
  - 完善Git双仓库同步工具文档
  - 涉及文件：10个文件，新增128行，删除1309行

#### Added
- **增加GitHub自动同步脚本** ([a4a85bf](https://github.com/Fischlvor/go_blog/commit/a4a85bfbbbbba118ab300fb901bb0a512233b0db))
  - 实现Gitee和GitHub双仓库自动同步
  - 支持增量同步和完整同步
  - 增加冲突检测和处理机制
  - 涉及文件：1个文件，新增287行

#### Fixed
- **优化健康检查配置和QQ登录参数** ([cf531b0](https://github.com/Fischlvor/go_blog/commit/cf531b01499f15b8d22dbd99e887febca7529347))
  - 优化Docker Compose健康检查配置
  - 修复QQ登录参数配置问题
  - 优化API请求基础配置
  - 涉及文件：4个文件，新增10行，删除10行

### 2025-11-15

#### Fixed
- **修复API请求路径和环境变量配置** ([259e144](https://github.com/Fischlvor/go_blog/commit/259e144fdeffd3a533ce821ca42e7dbba847530d))
  - 修复前端API请求路径配置
  - 优化环境变量管理
  - 增加Dockerfile环境变量支持
  - 优化emoji解析器和样式管理器
  - 涉及文件：10个文件，新增43行，删除39行

- **优化构建配置和.gitignore文件** ([08cd8a8](https://github.com/Fischlvor/go_blog/commit/08cd8a89566da2da23f28c64d45042792ef0b970))
  - 为所有服务增加.dockerignore文件
  - 优化.gitignore配置，避免提交敏感文件
  - 优化Dockerfile构建配置
  - 更新前端依赖包
  - 涉及文件：17个文件，新增1268行，删除301行

#### Added
- **增加emoji服务化配置** ([d5d4419](https://github.com/Fischlvor/go_blog/commit/d5d4419de62aad7c80b5691fbf3322a5e663fdf3))
  - 实现完整的emoji管理系统（分组、列表、雪碧图）
  - 增加emoji后台管理界面（分组管理、emoji上传、雪碧图生成）
  - 实现emoji版本控制和CDN部署
  - 优化emoji解析器，支持动态加载
  - 增加emoji选择器组件
  - 重构SSO认证流程，支持QQ登录
  - 增加设备管理和登录日志功能
  - 优化邮件服务和验证码功能
  - 涉及文件：121个文件，新增12587行，删除2462行

### 2025-11-10

#### Added
- **增加自动部署和DockerFile** ([ed938f2](https://github.com/Fischlvor/go_blog/commit/ed938f22fe1ae0aa595dba48aa20f2d049196e55))
  - 完善Docker Compose生产环境配置
  - 优化所有服务的Dockerfile
  - 增加Nginx配置文件
  - 完善自动化部署脚本
  - 支持单服务部署和全量部署
  - 涉及文件：16个文件，新增3127行，删除2957行

#### Fixed
- **优化依赖镜像** ([c0e4fc1](https://github.com/Fischlvor/go_blog/commit/c0e4fc16fccaec91b6d1f2eb7e09a5edc9cc8f88))
  - 优化Docker镜像构建速度
  - 使用国内镜像源加速依赖下载

#### Added
- **增加自动部署** ([21eee2d](https://github.com/Fischlvor/go_blog/commit/21eee2d16e9bb0154943273ffaa041d6b5cf4690))
  - 实现CI/CD自动化部署流程
  - 支持构建、上传、部署分步执行
  - 增加部署脚本和工作流配置

### 2025-11-08

#### Removed
- **删除无效文件** ([5c996c1](https://github.com/Fischlvor/go_blog/commit/5c996c102a2c65b8400ac1b210c3728e73946353))
  - 清理项目中的无效文件和冗余代码

#### Fixed
- **优化新用户自动注册应用** ([552cd41](https://github.com/Fischlvor/go_blog/commit/552cd418410b6fa127857148c14d1b167d889456))
  - 优化SSO新用户注册流程
  - 自动为新用户注册默认应用

#### Added
- **增加SSO单点登录系统** ([76480a8](https://github.com/Fischlvor/go_blog/commit/76480a808add65971dc3374901f297e3fc2e2f40))
  - 实现完整的SSO认证服务（server-auth-service）
  - 支持用户注册、登录、找回密码
  - 支持QQ第三方登录
  - 实现JWT + RSA非对称加密认证
  - 增加设备管理和登录日志
  - 实现SSO登录前端（web-auth-service）
  - 博客服务集成SSO认证
  - 增加SSO回调处理和令牌验证
  - 涉及文件：70个文件，新增6160行，删除434行

### 2025-08-27

#### Fixed
- **优化前端显示** ([a2e35fc](https://github.com/Fischlvor/go_blog/commit/a2e35fccb7231301dfca52badce4b6efc003b1fc))
  - 优化前端UI显示效果
  - 修复样式问题

### 2025-08-26

#### Fixed
- **增加自动部署文件** ([1729a39](https://github.com/Fischlvor/go_blog/commit/1729a39a20ff05e330b7f2eae0edbeaa5353661d))
  - 增加项目自动部署配置文件

### 2025-08-25

#### Fixed
- **修复配置文件路径** ([3410c61](https://github.com/Fischlvor/go_blog/commit/3410c61fd2692f38e93cbc9897bdecc31aa5c60d))
  - 修复后端配置文件路径问题

- **更新项目结构** ([e6cc62d](https://github.com/Fischlvor/go_blog/commit/e6cc62d5a9fabb1e28d09f83da21b9be8304918c))
  - 调整项目目录结构
  - 优化代码组织

### 2025-08-24

#### Fixed
- **nginx配置增加流式响应支持** ([772d872](https://github.com/Fischlvor/go_blog/commit/772d872a0e956b075f6db6d66354bad64a08de29))
  - 增加Nginx对AI聊天流式响应的支持
  - 配置SSE（Server-Sent Events）

- **增加nginx配置文件** ([df48672](https://github.com/Fischlvor/go_blog/commit/df4867207b977aac7466a429c48ab878c78c68ff))
  - 增加Nginx反向代理配置
  - 配置SSL证书和HTTPS

#### Added
- **增加AI助手后台管理** ([b560969](https://github.com/Fischlvor/go_blog/commit/b5609695fd9cc5da49554aaa806d047449af59e9))
  - 实现AI聊天后台管理功能
  - 增加会话管理、消息管理、模型管理
  - 涉及文件：24个文件，新增2691行，删除519行

### 2025-08-23

#### Added
- **上线AI对话助手** ([6f844c1](https://github.com/Fischlvor/go_blog/commit/6f844c1d77ee8691fa6ca316acec6c3fc7166e93))
  - 集成DeepSeek AI API
  - 实现AI聊天功能
  - 支持流式响应
  - 增加上下文管理
  - 实现AI助手前端界面

### 2025-08-17

#### Added
- **增加Docker Compose配置** ([e938bf4](https://github.com/Fischlvor/go_blog/commit/e938bf4595743d215f7a069a2ef89a813fb37914))
  - 完善Docker Compose部署配置
  - 支持容器化部署

#### Fixed
- **删除配置文件** ([26caae4](https://github.com/Fischlvor/go_blog/commit/26caae4b16bd5c3aacf8fe2115e065982cd29010))
  - 清理敏感配置文件

- **代码同步** ([39c5bf7](https://github.com/Fischlvor/go_blog/commit/39c5bf713ac3077964d98fa911452a3ca3304e8f))
  - 初始代码同步到仓库

---

**日志说明**：
- 日志按时间倒序排列（最新在上）
- 使用语义化版本分类：Added（新增）、Changed（变更）、Fixed（修复）、Removed（移除）
- 每个条目包含PR编号和详细的变更说明
- 重大功能更新会标注涉及的文件数量和代码行数

## �� 许可证

本项目采用 MIT 许可证。详情请参考 [LICENSE](./LICENSE) 文件。

## 👨‍💻 作者

- **qiyana423** - [Gitee](https://gitee.com/qiyana423) | [GitHub](https://github.com/Fischlvor)

## 🙏 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- [Vue.js](https://github.com/vuejs/core) - 前端框架
- [Element Plus](https://github.com/element-plus/element-plus) - UI 组件库
- [GORM](https://github.com/go-gorm/gorm) - Go ORM 框架

---

<div align="center">

如果这个项目对你有帮助，请给个 ⭐️ Star 支持一下！

</div>
