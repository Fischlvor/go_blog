# 项目结构与架构说明（PROJECT_STRUCTURE）

本文档仅包含两部分内容：

1. 当前仓库的文件结构（按实际目录）
2. 当前系统的架构关系（按实际服务与依赖）

---

## 1. 仓库文件结构

> 以下为主干结构（省略大量业务细节文件），用于快速定位模块职责。

```text
go_blog/
├── server-blog-v2/                  # 博客后端（Go + Fiber）
│   ├── cmd/
│   │   ├── app/                     # 服务启动入口
│   │   ├── migrate/                 # 数据库迁移入口
│   │   └── gen/                     # gorm/gen 代码生成入口
│   ├── config/                      # 配置结构定义
│   ├── configs/                     # 配置文件（example/dev/prod）
│   ├── internal/
│   │   ├── app/                     # 应用装配（wire）
│   │   ├── controller/              # HTTP 路由与处理器
│   │   ├── entity/                  # 领域实体
│   │   ├── repo/                    # 仓储接口与实现
│   │   └── usecase/                 # 业务用例层
│   ├── migrations/                  # SQL 迁移脚本
│   ├── pkg/                         # 基础设施封装（postgres/redis/logger/httpserver 等）
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── web-blog-v2/                     # 博客前端（Next.js + React）
│   ├── app/                         # 路由页面（App Router）
│   ├── components/                  # UI 组件
│   ├── context/                     # React 上下文
│   ├── lib/api/                     # 前端 API 请求层
│   ├── styles/
│   ├── package.json
│   ├── next.config.ts
│   └── Dockerfile
│
├── server-auth-service/             # SSO 认证后端（Go + Gin）
│   ├── configs/
│   ├── internal/
│   │   ├── api/                     # HTTP API
│   │   ├── initialize/              # 初始化（gorm/redis）
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── router/
│   │   └── service/
│   ├── pkg/
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
│
├── web-auth-service/                # 认证前端（Vue + Vite）
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   ├── views/
│   │   └── utils/
│   ├── package.json
│   ├── vite.config.js
│   └── Dockerfile
│
├── nginx/                           # Nginx 反向代理配置
│   ├── Dockerfile
│   ├── go_blog_dev.conf
│   └── go_blog.conf
│
├── docs/
│   ├── PROJECT_STRUCTURE.md         # 本文档
│   ├── CHANGELOG.md
│   ├── REFACTOR_PLAN.md
│   └── 资源上传系统设计方案.md
│
├── docker-compose.base.yml          # 基础设施编排
├── docker-compose.dev.yml           # 开发环境业务编排
├── docker-compose.prod.yml          # 生产环境业务编排
├── deploy_dev.sh                    # 本地开发部署脚本
├── deploy.example.sh                # 生产部署脚本示例
└── README.md
```

---

## 2. 系统架构

## 2.1 服务拓扑

```text
用户浏览器
   │
   ▼
 Nginx（统一入口与转发）
   ├── /api/v1, /api/admin, /api/callback -> server-blog-v2
   ├── 博客页面请求                      -> web-blog-v2
   ├── 认证页面请求                      -> web-auth-service
   └── 认证 API 请求                     -> server-auth-service
```

## 2.2 后端服务职责划分

### `server-blog-v2`（博客后端）

- 对外提供博客业务 API（公开侧与管理侧）
- 管理文章、评论、资源上传、AI 对话、广告、配置读取等能力
- 通过 SSO/JWT 中间件与认证系统协作
- 对接对象存储（七牛）与回调处理

主要分层：

- `controller`：路由与 HTTP 适配
- `usecase`：业务编排
- `repo`：数据访问与外部依赖抽象
- `pkg`：基础设施组件

### `server-auth-service`（认证后端）

- 提供统一认证、授权、令牌相关能力
- 负责登录、注册、OAuth、设备与安全相关流程
- 使用 JWT + RSA 提供可验证令牌

## 2.3 前端职责

### `web-blog-v2`

- 博客站点与管理后台 UI
- 通过 `lib/api` 对接博客后端与认证相关接口

### `web-auth-service`

- 认证交互页面（登录、注册、找回密码、OAuth 回调等）
- 对接 `server-auth-service`

## 2.4 数据与中间件分工

- `PostgreSQL`：博客主业务数据（`server-blog-v2`）
- `MySQL`：认证业务数据（`server-auth-service`）
- `Redis`：缓存、会话、热点等临时数据
- `Elasticsearch`：搜索能力
- `Qiniu`：对象存储（上传资源）

## 2.5 部署分层

- `docker-compose.base.yml`：先启动基础设施层
- `docker-compose.dev.yml`：开发环境业务层
- `docker-compose.prod.yml`：生产环境业务层

---

## 3. 维护说明

1. `server-blog-v2/internal/app/wire_gen.go` 是生成文件，修改依赖注入请改 `wire.go` 后重新生成。  
2. 目录命名存在历史痕迹（旧文档里的 `server-blog` / `web-blog`），当前以 `*-v2` 目录为准。  
3. 本文档只维护结构与架构，不包含详细部署步骤与业务说明。
