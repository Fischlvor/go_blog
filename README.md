# go_blog

`go_blog` 是一个多服务工程，包含博客系统与独立认证系统（SSO）。

本仓库采用“主 README + 专题文档”的文档结构：

- 主 README：给出项目全景、入口与路径
- 专题文档：给出可执行细节（启动、配置、部署、排障）

## 文档导航

- 架构文档：[`docs/PROJECT_STRUCTURE.md`](./docs/PROJECT_STRUCTURE.md)
- 日志文档：[`docs/CHANGELOG.md`](./docs/CHANGELOG.md)

## 一、项目组成

| 目录 | 角色 | 核心技术 |
|---|---|---|
| `server-blog-v2` | 博客后端 | Go 1.24, Fiber v3, GORM, PostgreSQL, Redis, Elasticsearch, Wire |
| `web-blog-v2` | 博客前端 | Next.js 16, React 19, TypeScript, Tailwind CSS 4 |
| `server-auth-service` | 认证后端（SSO） | Go 1.23, Gin, GORM, MySQL, Redis, JWT/RSA |
| `web-auth-service` | 认证前端 | Vue 3, Vite 5, Element Plus |
| `nginx` | 反向代理 | Nginx |

相关编排与脚本：

- `docker-compose.base.yml`：基础设施（MySQL、Redis、Elasticsearch、PostgreSQL）
- `docker-compose.dev.yml`：开发环境业务服务编排
- `docker-compose.prod.yml`：生产环境业务服务编排
- `deploy_dev.sh`：本地开发构建与容器操作脚本
- `deploy.example.sh`：生产部署示例脚本（需按实际环境修改）

## 二、功能范围（基于当前路由）

### 博客公开接口（`/api/v1`）

- 文章：搜索、详情、分类、标签、点赞
- 评论：查询、创建、删除
- AI 聊天：会话、消息、流式响应
- 反馈：提交与查询
- 网站信息：logo、标题、轮播、热点、日历、友链
- 用户：资料、登出、天气、图表
- 文件上传、表情配置、广告信息

### 博客管理接口（`/api/admin`）

- 文章、分类、标签、评论管理
- 用户与登录记录管理
- 文件/图片管理
- 资源分片上传（check/init/upload-chunk/complete/progress/cancel）
- AI 会话/消息/模型管理
- 广告管理
- 网站和系统配置读取

### 回调接口

- 七牛回调：`/api/callback/qiniu/*`

## 三、最短启动路径（摘要）

```bash
# 1) 启动基础设施
docker compose -f docker-compose.base.yml up -d

# 2) 准备配置
cd server-blog-v2 && cp configs/config.example.yaml configs/config.yaml
cd ../server-auth-service && cp configs/config.example.yaml configs/config.yaml

# 3) 启动后端
cd ../server-blog-v2 && go run ./cmd/app
cd ../server-auth-service && go run .

# 4) 启动前端
cd ../web-blog-v2 && pnpm install && pnpm dev
cd ../web-auth-service && npm install && npm run dev
```

完整步骤可根据本仓库 compose 文件与脚本执行（`docker-compose.base.yml`、`docker-compose.dev.yml`、`deploy_dev.sh`）。

## 四、常用命令

```bash
cd server-blog-v2
go run ./cmd/migrate -action up

go run -mod=mod github.com/google/wire/cmd/wire ./internal/app/...

go build ./...
```

## 五、维护约定（建议）

1. `wire_gen.go` 为生成文件，不要手工修改。  
2. 修改 API 路径或返回结构后，需同步检查 `web-blog-v2/lib/api` 与 `web-auth-service/src/api`。  
3. 修改 compose 或脚本后，同步更新对应文档。  
4. 禁止把真实密钥、密码提交到仓库；如果历史上出现过，应立即轮换。  

## 六、状态说明

- 仓库存在历史命名痕迹（旧文档中的 `server-blog`/`web-blog`）。当前以 `*-v2` 目录为主。
- 博客后端与认证后端是独立工程，数据库栈也不同（PostgreSQL vs MySQL）。
- 认证与博客之间通过 SSO/JWT（RSA 公钥验证）协作。
