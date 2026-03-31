## 📝 开发日志

### 2026-03-30

#### Added
- **前端完整重构，增加本地 Docker 部署脚本** ([01e0d53](https://github.com/Fischlvor/go_blog/commit/01e0d53))
  - 基于 Next.js 完整重构博客前端（`web-blog-v2`），替代原 Vue 版本
  - 新增完整的管理后台页面（文章、评论、AI管理、Emoji、资源、系统配置、用户管理等）
  - 新增 `web-blog-v2/Dockerfile` 及 `nginx.conf`，支持本地 Docker 部署
  - 更新 `deploy_dev.sh`、`docker-compose.dev.yml`、`docker-compose.prod.yml` 适配新前端
  - 更新 Nginx 配置（`nginx/go_blog.conf`、`nginx/go_blog_dev.conf`）
  - 涉及文件：80+ 个文件

- **清除无效文件夹** ([2d593f1](https://github.com/Fischlvor/go_blog/commit/2d593f1))
  - 删除项目中的冗余目录和废弃文件

### 2026-03-25

#### Added
- **前端初步重构** ([786d784](https://github.com/Fischlvor/go_blog/commit/786d784))
  - 初始化 `web-blog-v2` Next.js 前端项目框架
  - 搭建站点页面（首页、文章、归档、关于、友链）
  - 搭建管理后台基础布局（Header、Sidebar）
  - 新增 SSO 回调页面、用户认证 Context
  - 新增公共 UI 组件库（基于 shadcn/ui）
  - 配置 API 封装层（public、admin、user 分层）
  - 后端修复：`server-blog-v2` 用户用例输出字段调整
  - 涉及文件：85+ 个文件

### 2026-03-14

#### Fixed
- **修复文件上传异常** ([3f2b0d8](https://github.com/Fischlvor/go_blog/commit/3f2b0d8))
  - 修复 `server-blog-v2` 文件上传逻辑异常
  - 重构文件相关实体、仓储、用例层（`entity/file.go`、`repo/persistence/file_postgres.go`、`usecase/file/`）
  - 修复 `web-blog` 前端图片上传相关表单组件（文章创建/编辑、友情链接、广告等）
  - 涉及文件：26 个文件

- **修复 tag 配色** ([26783ce](https://github.com/Fischlvor/go_blog/commit/26783ce))
  - 修复文章列表页标签颜色显示异常
  - 影响文件：`web-blog/src/views/dashboard/articles/article-list.vue`

- **修复文章列表 tab 滚动** ([92adaf2](https://github.com/Fischlvor/go_blog/commit/92adaf2))
  - 修复文章列表 tab 切换时滚动位置异常
  - 完善文章创建/编辑表单交互逻辑
  - 影响文件：`ArticleCreateForm.vue`、`ArticleUpdateForm.vue`、`article-list.vue`、`article-publish.vue`

- **修复文章可见性问题** ([edc8a27](https://github.com/Fischlvor/go_blog/commit/edc8a27))
  - 修复文章可见性（公开/私有/仅自己）设置不生效的问题
  - 重构文章实体、分类、标签相关的数据库查询逻辑
  - 新增 `shared/validation.go` 统一参数校验
  - 更新数据库迁移脚本（`migrations/000001_init`）
  - 前端同步更新文章 API 和创建/编辑表单
  - 涉及文件：29 个文件，新增 1020 行，删除 1141 行

### 2026-03-13

#### Fixed
- **修复 AccessToken 不自动刷新** ([53dfe7d](https://github.com/Fischlvor/go_blog/commit/53dfe7d))
  - 新增 `internal/controller/http/middleware/session.go`，实现基于 Redis 的 Session 中间件
  - 实现 AccessToken 过期后自动静默刷新机制
  - 新增 `pkg/redis/redis.go` Redis 工具封装
  - 完善 `v1/auth.go` 认证接口
  - 涉及文件：11 个文件，新增 275 行，删除 35 行

### 2026-03-12

#### Refactor
- **后端切换为 Clean Architecture** ([50ef323](https://github.com/Fischlvor/go_blog/commit/50ef323))
  - 将 `server-blog-v2` 后端完整重构为 Clean Architecture（Controller → UseCase → Repository）
  - 新增 `cmd/app/main.go`、`cmd/gen/main.go`、`cmd/migrate/main.go` 多入口支持
  - 新增 Wire 依赖注入（`internal/app/wire.go`、`wire_gen.go`）
  - 拆分 Controller 层为 `admin/`、`v1/`、`qiniu/` 三个模块
  - 新增完整的 UseCase 层（content、user、file、ai等）
  - 新增 Repository 层（基于 PostgreSQL + GORM Gen）
  - 新增 `REFACTOR_PLAN.md` 重构计划文档
  - 涉及文件：50+ 个文件，新增大量代码

#### Fixed
- **修复前端接收字段** ([c77a0c5](https://github.com/Fischlvor/go_blog/commit/c77a0c5))
  - 修复 `web-blog` 前端友情链接相关 API 字段接收异常
  - 修复导航栏路由跳转问题
  - 影响文件：`src/api/friend-link.ts`、`WebNavbar.vue`、`router/index.ts`、`friend-link/index.vue`

### 2026-02-03

#### Fixed
- **修复 JWT 解析错误日志问题**
  - **问题描述**：日志中频繁出现 `"that's not even a token"` 和 `"couldn't handle this token"` 错误
  - **根本原因**：
    - `RateLimitMiddleware` 全局应用，在 SSO 中间件之前执行
    - `RateLimitKeyGetter` 调用 `utils.GetUUID(c)` 触发 JWT 解析
    - 公开接口没有 token，导致解析失败并记录错误日志
    - `LoginRecord` 中间件在公开路由上调用 `GetUUID`
  - **修复内容**：
    1. **调整限流中间件顺序**：
       - 公开接口：直接限流（按 IP/设备ID）
       - 私有接口：SSO 认证 → 限流（按用户 UUID）
       - 管理员接口：SSO 认证 → Admin 认证 → 限流
    2. **修改 `RateLimitKeyGetter`**：直接从 context 获取 claims，不调用 `GetUUID`
    3. **修改 `LoginRecord`**：直接从 context 获取 `user_uuid`，不调用 `GetUUID`
    4. **简化 `GetClaims`**：只从 context 获取 claims，移除 token 解析逻辑
  - **影响文件**：
    - `internal/initialize/router.go`
    - `internal/middleware/ratelimit.go`
    - `internal/middleware/login_record.go`
    - `pkg/utils/claims.go`

#### Refactor
- **AI 模型配置迁移到配置文件**
  - **问题描述**：`ai_tables.go` 中硬编码了 AI API Key
  - **修复内容**：
    1. 新增 `pkg/config/conf_ai.go` 定义 AI 配置结构
    2. 修改 `InitDefaultAIModels` 从配置文件读取默认模型
    3. 更新 `config.yaml`、`config.prod.yaml`、`config.example.yaml`
    4. 添加环境变量替换支持（`${ENV_VAR}` 格式）
  - **配置示例**：
    ```yaml
    ai:
        default_models:
            - name: deepseek-r1
              display_name: DeepSeek R1 (七牛云)
              provider: qiniu
              endpoint: https://openai.qiniu.com/v1/chat/completions
              api_key: your_api_key_here
              max_tokens: 4096
              temperature: 0.7
              is_active: true
    ```
  - **影响文件**：
    - `pkg/config/conf_ai.go`（新增）
    - `pkg/config/enter.go`
    - `pkg/utils/yaml.go`
    - `internal/initialize/ai_tables.go`
    - `configs/config.yaml`
    - `configs/config.prod.yaml`
    - `configs/config.example.yaml`

### 2026-01-09

#### Fixed
- **修复 QQ 登录 state 参数编码问题**
  - **问题描述**：QQ 登录回调时 state 参数解析失败，报错 "invalid character '%' looking for beginning of value"
  - **根本原因**：
    - 前端 Login.vue 中存在错误的本地 `encodeState` 函数：`btoa(encodeURIComponent(JSON.stringify(data)))`
    - 导致 state 被多次编码：JSON → URL编码 → Base64编码
    - 后端 QQLoginURL 函数对 state 进行了多余的 `url.QueryEscape(state)`
  - **修复内容**：
    - 删除 Login.vue 中错误的本地 `encodeState` 函数
    - 导入并使用 `@/utils/state.js` 中的正确实现
    - 添加 `nonce` 字段到 stateData（符合 OAuth 2.0 规范）
    - 移除后端 `auth.go` 中对 state 的多余 URL 编码
    - 保持 `state.go` 中简单的 Base64 解码逻辑
  - **编码流程优化**：
    - 修复前：JSON → URL编码 → Base64 → URL编码（后端）→ 3层编码
    - 修复后：JSON → Base64 → 1层编码
  - **影响文件**：
    - `web-auth-service/src/views/Login.vue`
    - `server-auth-service/internal/api/auth.go`
    - `server-auth-service/pkg/utils/state.go`

- **完善邮箱验证码冷却机制**
  - **实现内容**：
    - 使用 Go 自定义错误类型 `CooldownError` 代替字符串解析
    - 创建 `pkg/errors/email_errors.go` 定义冷却错误
    - 后端返回业务错误代码 1013（发送频率限制）
    - 前端根据 `remaining_seconds` 字段动态设置倒计时
    - 图形验证码弹窗在确认后立即关闭
  - **用户体验优化**：
    - 点击发送后立即禁用表单和按钮并变灰
    - 成功则保持禁用，失败则恢复可用
    - 刷新页面后仍显示正确的剩余冷却时间
  - **符合 Go 最佳实践**：使用自定义错误类型而非错误码或字符串前缀判断

- **优化登录按钮禁用逻辑**
  - 添加 `canLogin` 计算属性验证表单完整性
  - 密码登录：必须填写邮箱、密码和图形验证码
  - 验证码登录：必须填写邮箱和邮箱验证码
  - 防止提交不完整的表单

#### Added
- **实现 SSO 主页和登录页主题化功能**
  - **SSO 主页 (Home.vue)**：
    - 新增公共应用列表展示功能
    - 实现左右分栏布局（左侧展示区 + 右侧应用卡片）
    - 添加动态渐变背景动画（3色渐变 + 15秒循环）
    - 支持应用图标展示（自适应矩形/正方形）
    - 实现应用点击跳转到 OAuth 授权流程
    - 后端新增 `/api/oauth/applications` API，返回已启用的公共应用列表
    - 根据环境变量 `APP_ENV` 自动选择生产/测试环境的 `redirect_uri`
  - **登录页主题化 (Login.vue)**：
    - 根据 `app_id` 参数动态切换主题（Blog蓝色、MCP绿色、默认紫色）
    - 实现对称渐变背景（浅色 → 深色 → 浅色），配合 `background-position` 动画
    - 主题色方案：
      - **Blog蓝色**：#bfdbfe（浅蓝）→ #a5f3fc（青色）→ #bfdbfe（浅蓝）
      - **MCP绿色**：#a7f3d0（浅绿）→ #fde68a（黄绿）→ #a7f3d0（浅绿）
      - **默认紫色**：#ddd6fe（浅紫）→ #fbcfe8（粉色）→ #ddd6fe（浅紫）
    - 动态调整圆形动画、Logo颜色、标题渐变、按钮样式等元素的主题色
    - 添加 `gradientShift` 动画（15秒循环，背景位置移动）
  - **后端改动**：
    - 新增 `ApplicationService.GetPublicApplications()` 方法
    - 新增 `PublicApplicationInfo` 响应结构体
    - 新增 `IsProduction()` 环境检测工具函数
    - OAuth 路由新增 `GET /api/oauth/applications` 端点
  - **前端改动**：
    - 新增 `web-auth-service/src/api/oauth.js` API 封装
    - 新增 `web-auth-service/src/views/Home.vue` 主页组件
    - 修改 `web-auth-service/src/views/Login.vue` 添加主题化逻辑
    - 修改 `web-auth-service/src/main.js` 添加主页路由
  - **涉及文件**：11个文件（新增3个，修改8个）

### 2026-01-06

#### Fixed
- **修复 QQ 登录时 Email 唯一索引冲突问题**
  - **问题描述**：QQ 登录创建用户时未设置 Email 字段，导致插入空字符串，多个 QQ 用户注册时唯一索引冲突
  - **修复内容**：
    - 将 `SSOUser.Email` 字段类型从 `string` 改为 `*string`
    - 添加 GORM 标签 `default:null`，支持 NULL 值
    - 将 `UserInfo.Email` 响应字段同步改为 `*string`
    - 邮箱注册时传递指针 `&req.Email`
  - **效果**：QQ 登录用户 Email 为 NULL，不参与唯一索引约束
  - **影响文件**：`sso_user.go`、`auth_response.go`、`auth_service.go`

#### Added
- **实现完整的资源上传系统**
  - **核心功能**：
    - ✅ 秒传检测（基于文件 MD5 Hash）
    - ✅ 断点续传（分片上传 + 七牛云 Context）
    - ✅ 并发上传（3 个并发，可配置）
    - ✅ 流式传输（边接收边转发，低内存占用）
    - ✅ 进度跟踪（实时显示上传进度）
    - ✅ 文件类型验证（MIME 类型 + Magic Number 签名）
  - **后端架构**：
    - 新增 `Resource` 模型（资源表，支持秒传）
    - 新增 `ResourceUploadTask` 模型（上传任务表，支持断点续传）
    - 新增 `ResourceUploader` 接口（支持切换云存储提供商）
    - 新增 `QiniuResourceUploader` 实现（七牛云分片上传 V1 API）
    - 新增 `FileValidator` 验证器（MIME 类型 + 文件签名验证）
    - 新增 6 个 API 端点：check、init、upload-chunk、complete、cancel、progress
    - 新增七牛云转码回调处理（`/api/qiniu/callback`）
  - **前端实现**：
    - 新增资源上传组件（拖拽上传、进度显示、状态管理）
    - 新增资源列表页面（分页、搜索、预览、复制链接、删除）
    - 使用 spark-md5 计算文件 Hash（支持大文件分片计算）
    - 视频缩略图显示（使用后端返回的 `thumbnail_url`）
    - 转码状态显示（转码中/已转码/转码失败）
  - **数据库设计**：
    - `resources` 表：存储已上传完成的资源，支持秒传（相同 Hash 复用物理文件）
    - `resource_upload_tasks` 表：存储分片上传临时状态，支持断点续传
    - 新增转码相关字段：`transcode_status`、`transcode_key`、`thumbnail_key`
  - **涉及文件**：24 个文件，新增 3326 行代码
  - **文档**：新增 `docs/资源上传系统设计方案.md`（745 行详细设计文档）

### 2025-12-25

#### Fixed
- **SSO 静默登录缺失设备信息导致的重定向与"SSO 设备"命名问题**
  - 登录与 QQ 登录成功后，将 `device_name`、`device_type` 一并写入 SSO Session，后续静默授权可直接复用浏览器/设备信息
  - OAuth `Authorize` 端点读取 Session 中的设备信息，当 `CheckDeviceExpiry` 报告 `ErrDeviceNotFound` 时，认定为首次访问应用并保持静默登录流程
  - `GenerateTokensForUser` 在自动注册设备时优先使用上下文提供的名称与类型，仅在缺失时回退默认值，避免生成泛化的"SSO 设备"记录
  - 新增 `ErrDeviceNotFound` 哨兵错误，便于区分"设备首次访问"与"设备被移除/过期"，同时停止无意义的 `sso_device_id` 伪造
  - 影响文件：`server-auth-service/internal/api/auth.go`、`server-auth-service/internal/api/oauth.go`、`server-auth-service/internal/service/auth_service.go`
- **QQ 登录服务补充 gin.Context 与会话写入，保持静默登录上下文一致**
  - `AuthApi.QQLogin/QQCallback` 将 `gin.Context` 传入 Service，统一由服务层处理 QQ 登录回调
  - `AuthService.QQLogin` 在成功登录后，将 `device_name`、`device_type`、`user_agent`、`ip_address` 等信息写入 SSO Session，保障后续静默授权直接复用
  - 新增 `github.com/gin-contrib/sessions` 依赖，保存 session 失败时输出结构化日志并阻断异常登录
  - 影响文件：`server-auth-service/internal/api/auth.go`、`server-auth-service/internal/service/auth_service.go`

### 2025-12-17

#### Fixed
- **改进 SSO Refresh Token 存储和安全校验，支持多应用隔离**
  - **问题描述**：
    1. `RefreshToken` 方法查询设备时使用 `claims.AppID`（字符串）而非 `app.ID`（数字）
       - 导致 SQL 变成 `app_id = 'mcp'`，查询失败
    3. 原逻辑先用 `req.ClientID` 查询应用，再解析 token
       - 没有验证 `req.ClientID` 与 `claims.AppID` 是否一致，存在跨应用攻击风险
  
  - **修复内容**：
    1. **修复 app_id 类型错误**：
       - 使用 `app.ID`（数字）替代 `claims.AppID`（字符串）查询设备
       - 添加 OAuth 标准字段与数据库字段的映射注释
    
    2. **增强安全校验**：
       - 调整逻辑顺序：先解析 token → 校验一致性 → 查询应用
       - 添加安全校验：`req.ClientID` 必须与 `claims.AppID` 一致
       - 防止跨应用刷新 token 攻击
  
  - **新的验证流程**：
    ```
    1. 解析 RefreshToken 获取 claims
    2. 校验 req.ClientID == claims.AppID（防止跨应用攻击）
    3. 用 claims.AppID 查询应用并验证 client_secret
    4. 检查用户和设备状态
    5. 生成新 token
    ```
  
  - **字段映射说明**：
    | OAuth 标准字段 | 数据库字段 | 说明 |
    |---------------|-----------|------|
    | `client_id` | `app_key` | 应用标识（字符串，如 "mcp"） |
    | `client_secret` | `app_secret` | 应用密钥 |
    | JWT `AppID` | `app_key` | 字符串，需转换为 `app.ID` |
    | - | `app.ID` | 应用数字 ID（用于数据库关联） |
  
  - **影响范围**：
    - 修改文件：`internal/service/auth_service.go`、`internal/model/request/auth_request.go`、`pkg/jwt/jwt.go`
    - 修改方法：8 个
    - 修改行数：约 30 行

### 2025-12-15

#### Fixed
- **修复 SSO 设备管理跨应用误踢问题**
  - **问题描述**：
    - 用户在应用 A 登录后，应用 B 的同名设备被误踢出
    - 静默登录时，`CheckDeviceExpiry` 只使用 `device_id` 查询，可能查到其他用户的设备
    - 踢出设备时（3个方法）没有区分 `app_id`，导致所有应用的同名设备都被踢出
    - 设备数量限制判断逻辑错误（`>=` 应为 `>`），导致设备数量少于限制时仍被踢出
  
  - **修复内容**：
    1. **修复 `CheckDeviceExpiry` 方法**
       - 添加 `user_uuid` 和 `app_id` 参数
       - 修改查询条件，包含 `user_uuid` 和 `app_id`
       - 更新 `oauth.go` 中的调用，传入正确的参数
    
    2. **修复设备踢出逻辑**（3个方法）
       - `Logout` 方法：添加 `app_id` 条件
       - `KickDevice` 方法：添加 `app_id` 条件
       - `kickDeviceInternal` 方法：添加 `app_id` 条件
    
    3. **修复设备数量限制**
       - `handleDeviceLimit` 方法：修正判断逻辑（`>=` 改为 `>`）
       - 允许 `maxDevices` 个设备同时在线
    
    4. **修复设备活跃时间更新**
       - `CheckDeviceExpiry` 中的 `last_active_at` 更新已包含 `app_id` 条件
  
  - **影响范围**：
    - 修改文件：`internal/service/auth_service.go`、`internal/api/oauth.go`
    - 修改方法：5个核心方法
    - 修复了跨应用设备管理的隔离问题

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
