# goBlog

#### 介绍
goBlog 是一个基于 Go + Vue3 的现代化博客系统，采用前后端分离架构，提供完整的博客管理功能。

**技术栈：**
- **后端**: Go + Gin + GORM + MySQL + Redis + Elasticsearch
- **前端**: Vue3 + TypeScript + Element Plus + Vite
- **AI功能**: 集成DeepSeek AI聊天助手

#### 软件架构

```
goBlog/
├── server/                    # 后端服务
│   ├── cmd/server/           # 应用程序入口点
│   │   └── main.go          # 主程序入口
│   ├── internal/             # 私有应用程序和库代码
│   │   ├── api/             # API层 - HTTP处理器
│   │   ├── service/         # 业务逻辑层
│   │   ├── model/           # 数据模型
│   │   │   ├── database/    # 数据库模型
│   │   │   ├── request/     # 请求模型
│   │   │   ├── response/    # 响应模型
│   │   │   └── ...
│   │   ├── middleware/      # 中间件
│   │   ├── router/          # 路由配置
│   │   ├── initialize/      # 初始化模块
│   │   └── task/            # 定时任务
│   ├── pkg/                 # 可重用库代码
│   │   ├── config/          # 配置管理
│   │   ├── utils/           # 工具函数
│   │   ├── global/          # 全局变量
│   │   └── core/            # 核心功能
│   ├── configs/             # 配置文件
│   ├── scripts/             # 脚本文件
│   ├── log/                 # 日志文件
│   └── uploads/             # 上传文件
├── web/                     # 前端应用
│   ├── src/
│   │   ├── components/      # Vue组件
│   │   ├── views/           # 页面视图
│   │   ├── api/             # API接口
│   │   └── utils/           # 工具函数
│   └── ...
└── docs/                    # 项目文档
```

#### 主要功能

- **用户管理**: 用户注册、登录、权限管理
- **文章管理**: 文章发布、编辑、分类、标签
- **评论系统**: 文章评论、回复功能
- **AI助手**: 集成DeepSeek AI聊天功能
- **内容管理**: 广告、友情链接、网站配置
- **数据统计**: 文章浏览量、点赞数统计
- **搜索功能**: 基于Elasticsearch的全文搜索

#### 安装教程

1. **克隆项目**
   ```bash
   git clone https://gitee.com/your-repo/go_blog.git
   cd go_blog
   ```

2. **后端配置**
   ```bash
   cd server
   # 安装依赖
   go mod tidy
   # 配置数据库
   cp configs/config.yaml.example configs/config.yaml
   # 修改配置文件中的数据库连接信息
   # 运行项目
   go run main.go
   ```

3. **前端配置**
   ```bash
   cd web
   # 安装依赖
   npm install
   # 开发模式运行
   npm run dev
   # 构建生产版本
   npm run build
   ```

#### 使用说明

1. **开发环境启动**
   - 后端: `cd server && go run main.go`
   - 前端: `cd web && npm run dev`

2. **生产环境部署**
   - 后端: `cd server && go build -o main . && ./main`
   - 前端: `cd web && npm run build`

3. **数据库迁移**
   - 项目会自动创建数据库表结构

#### 项目特色

- **现代化架构**: 采用Go标准项目结构，代码组织清晰
- **AI集成**: 内置AI聊天助手，提升用户体验
- **高性能**: 使用Redis缓存、Elasticsearch搜索
- **响应式设计**: 前端采用Vue3 + Element Plus
- **完整功能**: 涵盖博客系统的所有核心功能

#### 技术亮点

1. **Go标准项目结构**: 遵循Go项目最佳实践
2. **前后端分离**: 清晰的API接口设计
3. **AI功能集成**: 智能聊天助手
4. **高性能架构**: 缓存 + 搜索引擎
5. **现代化前端**: Vue3 + TypeScript + Element Plus
