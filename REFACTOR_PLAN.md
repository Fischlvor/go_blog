# go_blog 重构方案

> 参考项目: [nimbus-blog-api](https://github.com/scc749/nimbus-blog-api)
> 目标: 采用 Clean Architecture 提升代码可维护性、可测试性

## 一、项目定位

| 服务 | 定位 | 重构策略 |
|------|------|----------|
| **server-auth-service** | 独立 SSO 认证服务，供多应用使用 | 独立重构，保持独立性 |
| **server-blog** | 博客核心服务，依赖 auth-service | 彻底重构为 Clean Architecture |

---

## 二、技术选型

| 组件 | 当前 | 重构后 | 理由 |
|------|------|--------|------|
| **Web 框架** | Gin | **Fiber v3** | 高性能，基于 fasthttp |
| **ORM** | GORM (反射) | **GORM Gen** | 类型安全，编译时检查 |
| **数据库** | MySQL | **PostgreSQL** | 功能更强，JSONB 支持 |
| **搜索** | Elasticsearch | Elasticsearch | 保持，接口化 |
| **缓存** | Redis | Redis | 保持，接口化 |
| **存储** | 七牛云 | 七牛云 | 保持，接口化 |
| **日志** | Zap | **Zerolog** | 更轻量，JSON 友好 |
| **配置** | 自定义 | **Viper** | 标准化，支持多格式 |
| **依赖注入** | 全局变量 | **Wire** | 编译时注入，类型安全 |
| **API 文档** | 无 | **Swagger** | 自动生成 |

---

## 三、目标架构

### 3.1 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│  Controller (HTTP)                                          │
│  - 只处理 HTTP 请求/响应                                      │
│  - 参数验证、错误码映射                                        │
│  - 调用 UseCase，不包含业务逻辑                                │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  UseCase (业务用例)                                          │
│  - 纯业务逻辑                                                 │
│  - 编排多个 Repo 完成业务                                     │
│  - 不关心数据如何存储                                         │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  Repo (数据仓库接口)                                         │
│  - 定义数据访问契约                                           │
│  - 可替换实现 (PostgreSQL/MySQL/Mock)                        │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  Entity (领域实体)                                           │
│  - 纯数据结构，无依赖                                         │
│  - 业务核心概念                                               │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 目录结构

```
server-blog/
├── cmd/
│   ├── app/main.go              # 应用入口
│   ├── migrate/main.go          # 数据库迁移工具
│   └── gen/main.go              # GORM Gen 代码生成
├── config/
│   └── config.go                # 配置结构体
├── config.yaml                  # 配置文件
├── internal/
│   ├── app/                     # 应用启动 + Wire 依赖注入
│   │   ├── app.go               # 应用生命周期管理
│   │   ├── wire.go              # Wire 依赖声明
│   │   └── wire_gen.go          # Wire 生成的代码
│   ├── controller/
│   │   └── http/
│   │       ├── admin/           # 后台管理 API
│   │       │   ├── article.go
│   │       │   ├── comment.go
│   │       │   └── ...
│   │       ├── v1/              # 公开 API
│   │       │   ├── article.go
│   │       │   ├── comment.go
│   │       │   ├── ai_chat.go
│   │       │   └── ...
│   │       ├── middleware/      # 中间件
│   │       │   ├── jwt.go
│   │       │   ├── logger.go
│   │       │   └── ratelimit.go
│   │       ├── shared/          # 共享响应/错误码
│   │       │   ├── response.go
│   │       │   └── errors.go
│   │       └── router.go        # 路由注册
│   ├── usecase/                 # 业务用例层
│   │   ├── contracts.go         # 所有 UseCase 接口定义
│   │   ├── article/
│   │   │   └── article.go
│   │   ├── comment/
│   │   │   └── comment.go
│   │   ├── ai_chat/
│   │   │   └── ai_chat.go
│   │   ├── feedback/
│   │   │   └── feedback.go
│   │   └── input/output/        # DTO
│   │       ├── article.go
│   │       └── ...
│   ├── repo/                    # 数据仓库层
│   │   ├── contracts.go         # 所有 Repo 接口定义
│   │   ├── persistence/         # PostgreSQL 实现 (GORM Gen)
│   │   │   ├── article_repo.go
│   │   │   ├── comment_repo.go
│   │   │   └── gen/             # GORM Gen 生成的代码
│   │   ├── search/              # Elasticsearch 实现
│   │   │   └── article_search.go
│   │   ├── cache/               # Redis 缓存实现
│   │   │   ├── article_cache.go
│   │   │   └── view_buffer.go
│   │   └── storage/             # 七牛云实现
│   │       └── qiniu.go
│   └── entity/                  # 领域实体
│       ├── article.go
│       ├── comment.go
│       ├── user.go
│       └── ...
├── pkg/                         # 公共库
│   ├── httpserver/              # HTTP 服务器封装
│   ├── logger/                  # Zerolog 封装
│   ├── postgres/                # PostgreSQL 连接
│   ├── redis/                   # Redis 连接
│   ├── elasticsearch/           # ES 客户端
│   └── sso/                     # SSO 客户端 SDK
├── migrations/                  # 数据库迁移文件
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
└── docs/                        # Swagger 文档
    └── swagger.json
```

---

## 四、核心接口设计

### 4.1 Repo 接口 (repo/contracts.go)

```go
package repo

import (
    "context"
    "time"
    "server/internal/entity"
)

// ArticleRepo 文章数据仓库 (PostgreSQL)
type ArticleRepo interface {
    List(ctx context.Context, offset, limit int, filters ArticleFilters) ([]*entity.Article, int64, error)
    GetByID(ctx context.Context, id uint) (*entity.Article, error)
    Create(ctx context.Context, article *entity.Article) (uint, error)
    Update(ctx context.Context, article *entity.Article) error
    Delete(ctx context.Context, id uint) error
}

// ArticleSearchRepo 文章搜索仓库 (Elasticsearch)
type ArticleSearchRepo interface {
    Index(ctx context.Context, article *entity.Article) error
    Search(ctx context.Context, query string, filters SearchFilters) ([]*entity.Article, int64, error)
    Delete(ctx context.Context, id string) error
    BulkIndex(ctx context.Context, articles []*entity.Article) error
}

// ArticleCacheRepo 文章缓存仓库 (Redis)
type ArticleCacheRepo interface {
    GetViews(ctx context.Context, id string) (int, error)
    IncrViews(ctx context.Context, id string) error
    GetHotList(ctx context.Context) ([]string, error)
    SetHotList(ctx context.Context, ids []string, ttl time.Duration) error
}

// CommentRepo 评论数据仓库
type CommentRepo interface {
    ListByArticleID(ctx context.Context, articleID uint, offset, limit int) ([]*entity.Comment, int64, error)
    Create(ctx context.Context, comment *entity.Comment) (uint, error)
    Delete(ctx context.Context, id uint) error
}

// StorageRepo 对象存储仓库 (七牛云)
type StorageRepo interface {
    Upload(ctx context.Context, key string, data io.Reader, size int64) (string, error)
    Delete(ctx context.Context, key string) error
    GetURL(ctx context.Context, key string) string
}

// AIChatRepo AI 聊天仓库
type AIChatRepo interface {
    CreateSession(ctx context.Context, session *entity.ChatSession) (uint, error)
    GetSession(ctx context.Context, id uint) (*entity.ChatSession, error)
    ListSessions(ctx context.Context, userID string, offset, limit int) ([]*entity.ChatSession, int64, error)
    SaveMessage(ctx context.Context, msg *entity.ChatMessage) error
    ListMessages(ctx context.Context, sessionID uint) ([]*entity.ChatMessage, error)
}
```

### 4.2 UseCase 接口 (usecase/contracts.go)

```go
package usecase

import (
    "context"
    "server/internal/usecase/input"
    "server/internal/usecase/output"
)

// ArticleUseCase 文章业务用例
type ArticleUseCase interface {
    // 公开 API
    List(ctx context.Context, params input.ListArticles) (*output.ArticlePage, error)
    GetByID(ctx context.Context, id uint) (*output.ArticleDetail, error)
    Search(ctx context.Context, params input.SearchArticles) (*output.ArticlePage, error)
    RecordView(ctx context.Context, id uint, ip, userAgent string)
    
    // 管理 API
    Create(ctx context.Context, params input.CreateArticle) (uint, error)
    Update(ctx context.Context, params input.UpdateArticle) error
    Delete(ctx context.Context, id uint) error
}

// CommentUseCase 评论业务用例
type CommentUseCase interface {
    List(ctx context.Context, articleID uint, params input.ListComments) (*output.CommentPage, error)
    Create(ctx context.Context, params input.CreateComment) (uint, error)
    Delete(ctx context.Context, id uint) error
}

// AIChatUseCase AI 聊天业务用例
type AIChatUseCase interface {
    CreateSession(ctx context.Context, userID string, params input.CreateSession) (*output.Session, error)
    SendMessage(ctx context.Context, sessionID uint, params input.SendMessage) (<-chan string, error)
    ListSessions(ctx context.Context, userID string, params input.ListSessions) (*output.SessionPage, error)
    GetMessages(ctx context.Context, sessionID uint) ([]*output.Message, error)
}

// FeedbackUseCase 反馈业务用例
type FeedbackUseCase interface {
    Create(ctx context.Context, params input.CreateFeedback) error
    List(ctx context.Context, params input.ListFeedback) (*output.FeedbackPage, error)
}
```

---

## 五、Fiber 迁移指南

### 5.1 中间件签名变化

```go
// Gin
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ...
        c.Next()
    }
}

// Fiber
func JWTAuth() fiber.Handler {
    return func(c fiber.Ctx) error {
        // ...
        return c.Next()
    }
}
```

### 5.2 常用 API 对照

| 操作 | Gin | Fiber |
|------|-----|-------|
| 路径参数 | `c.Param("id")` | `c.Params("id")` |
| 查询参数 | `c.Query("page")` | `c.Query("page")` |
| JSON 绑定 | `c.ShouldBindJSON(&req)` | `c.Bind().JSON(&req)` |
| JSON 响应 | `c.JSON(200, data)` | `c.JSON(data)` |
| 设置状态码 | `c.JSON(404, ...)` | `c.Status(404).JSON(...)` |
| 获取 Header | `c.GetHeader("X-Token")` | `c.Get("X-Token")` |
| 设置 Header | `c.Header("key", "val")` | `c.Set("key", "val")` |
| 终止请求 | `c.Abort()` | `return err` |
| Context 存值 | `c.Set("key", val)` | `c.Locals("key", val)` |
| Context 取值 | `c.Get("key")` | `c.Locals("key")` |
| 获取 IP | `c.ClientIP()` | `c.IP()` |

### 5.3 Session 替换

```go
// Gin: github.com/gin-contrib/sessions
// Fiber: github.com/gofiber/fiber/v2/middleware/session
```

---

## 六、GORM Gen 使用

### 6.1 代码生成配置

```go
// cmd/gen/main.go
package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gen"
    "gorm.io/gorm"
)

func main() {
    g := gen.NewGenerator(gen.Config{
        OutPath:      "./internal/repo/persistence/gen",
        ModelPkgPath: "./internal/entity",
        Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
    })

    db, _ := gorm.Open(postgres.Open("..."))
    g.UseDB(db)

    // 生成所有表的模型和查询方法
    g.ApplyBasic(
        g.GenerateModel("articles"),
        g.GenerateModel("comments"),
        g.GenerateModel("users"),
        // ...
    )

    g.Execute()
}
```

### 6.2 类型安全查询

```go
// 反射方式 (当前)
db.Where("title LIKE ?", "%"+keyword+"%").Find(&articles)

// GORM Gen (类型安全)
q := query.Article
articles, _ := q.Where(q.Title.Like("%" + keyword + "%")).Find()
```

---

## 七、Wire 依赖注入

### 7.1 Provider 定义

```go
// internal/app/wire.go
//go:build wireinject

package app

import (
    "github.com/google/wire"
    "server/internal/repo/persistence"
    "server/internal/repo/search"
    "server/internal/repo/cache"
    "server/internal/usecase/article"
    "server/internal/controller/http"
)

var ProviderSet = wire.NewSet(
    // 基础设施
    NewPostgres,
    NewRedis,
    NewElasticsearch,
    NewQiniuClient,
    
    // Repo 实现
    persistence.NewArticleRepo,
    search.NewArticleSearchRepo,
    cache.NewArticleCacheRepo,
    
    // UseCase
    article.NewArticleUseCase,
    
    // HTTP Server
    http.NewRouter,
    NewHTTPServer,
)

func InitializeApp(cfg *config.Config) (*App, func(), error) {
    wire.Build(ProviderSet)
    return nil, nil, nil
}
```

### 7.2 生成代码

```bash
# 安装 Wire
go install github.com/google/wire/cmd/wire@latest

# 生成依赖注入代码
wire ./internal/app
```

---

## 八、重构步骤

### 第一阶段：基础设施 (1-2 周)

- [ ] 创建新项目骨架目录结构
- [ ] 配置 Wire 依赖注入
- [ ] 配置 GORM Gen 代码生成
- [ ] 配置 Swagger 文档生成
- [ ] 迁移 pkg/ 公共库
  - [ ] pkg/postgres (新建)
  - [ ] pkg/redis
  - [ ] pkg/elasticsearch
  - [ ] pkg/logger (Zerolog)
- [ ] 创建 migrations/ 目录，使用 golang-migrate

### 第二阶段：核心模块 (2-3 周)

- [ ] Article 模块 (最复杂，先做)
  - [ ] entity/article.go
  - [ ] repo/persistence/article_repo.go
  - [ ] repo/search/article_search.go
  - [ ] repo/cache/article_cache.go
  - [ ] usecase/article/
  - [ ] controller/http/v1/article.go
- [ ] Comment 模块
- [ ] User 模块 (对接 SSO)
- [ ] Feedback/FriendLink 模块

### 第三阶段：特殊功能 (1-2 周)

- [ ] AI Chat 模块
  - [ ] 流式响应处理 (SSE)
  - [ ] 会话管理
- [ ] 定时任务
  - [ ] 抽象为 Scheduler UseCase
  - [ ] 依赖注入 Repo
- [ ] 文件上传
  - [ ] StorageRepo 接口
  - [ ] 七牛云实现

### 第四阶段：收尾 (1 周)

- [ ] Swagger 文档完善
- [ ] 单元测试 (目标覆盖率 70%+)
- [ ] 集成测试
- [ ] Docker 配置更新
- [ ] README 更新

---

## 九、数据迁移 (后续)

重构阶段使用 PostgreSQL 空库开发，数据迁移后续单独处理：

1. MySQL 导出 JSON/CSV
2. 数据转换 (时间格式、布尔值等)
3. PostgreSQL 导入
4. ES 数据重建 (可选)

---

## 十、预期收益

| 指标 | 当前 | 重构后 |
|------|------|--------|
| **单元测试覆盖率** | ~0% | 70%+ |
| **代码可读性** | 中等 | 高 |
| **新功能开发速度** | 慢 | 快 |
| **Bug 定位时间** | 长 | 短 |
| **技术债务** | 高 | 低 |
| **可替换性** | 低 | 高 (接口化) |

---

## 十一、参考资料

- [nimbus-blog-api](https://github.com/scc749/nimbus-blog-api) - 参考项目
- [Fiber v3 文档](https://docs.gofiber.io/)
- [GORM Gen 文档](https://gorm.io/gen/)
- [Google Wire](https://github.com/google/wire)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
