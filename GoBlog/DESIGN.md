## GoBlog 个人博客系统设计说明

### 1. 概览

- **项目名称**：GoBlog
- **技术栈**：
  - Web 框架：`gin`
  - 数据库：`PostgreSQL`
  - ORM：`gorm`
  - 日志：标准库 `log/slog`
  - Markdown 渲染：`github.com/yuin/goldmark`
- **架构风格**：模块化单体（分层：Transport / Usecase / Domain / Infra）
- **目标场景**：个人或少量作者使用的博客系统，支持 Markdown 写作、前台浏览和基础 SEO。

### 2. 需求设计

#### 2.1 角色与使用场景

- **访客（匿名用户）**
  - 浏览文章列表（按时间倒序）
  - 查看文章详情
  - 按标签/分类筛选（后续扩展）
  - 使用简单搜索（后续扩展）
  - 访问 RSS / Sitemap（后续扩展）

- **作者 / 管理员**
  - 登录后台（鉴权后续补充）
  - 通过 JSON / 后台页面创建、编辑、发布文章
  - 管理标签和分类（后续扩展）
  - 配置站点信息（后续扩展）

#### 2.2 MVP 功能范围

- **必须（MVP）**
  - 文章：支持 Markdown 写作，存储 Markdown 与预渲染 HTML
  - 状态：草稿（draft）、已发布（published）、已归档（archived）
  - 前台功能：
    - 首页展示最近文章列表（JSON / 后续可替换为 HTML 模板）
    - 文章详情页（按 slug 访问）
  - 后台接口：
    - 创建文章（支持是否立即发布）
  - 基础健康检查：`/health`

- **后续扩展（MVP+ / v2）**
  - 标签、分类管理与筛选
  - RSS / Sitemap / robots.txt
  - 前台 SSR 模板渲染
  - 后台管理页面（基于模板或前后端分离）
  - 评论系统（优先考虑第三方）
  - 搜索（SQL LIKE / 外部搜索服务）

#### 2.3 非功能需求

- **性能**
  - 个人博客场景下，单机足够；常见页面 P95 延迟 < 200ms（不含复杂搜索）。
- **可用性**
  - 单实例部署即可；后续可扩展为多实例 + 负载均衡。
- **可维护性**
  - 清晰分层：Transport、Usecase、Domain、Infra
  - 使用 gorm + 明确数据模型，避免隐式迁移副作用。
- **安全**
  - 后台接口需要鉴权（后续接入中间件）
  - 输入校验（gin binding）
  - Markdown 渲染时使用可信库，输出在模板层默认转义（后续 SSR 时考虑 XSS）。

### 3. 架构设计

#### 3.1 分层结构

代码目录结构（简化）：

```text
GoBlog/
  cmd/
    blog/
      main.go                # 入口，组装依赖并启动 HTTP 服务
  internal/
    config/
      config.go              # 配置加载（环境变量）
    domain/
      post.go                # 领域实体：Post, Tag, Category
    usecase/
      post.go                # 用例：创建文章、按 slug 查询、列出最近文章
    infra/
      db/
        postgres.go          # gorm + PostgreSQL 初始化
      logging/
        logger.go            # slog 日志初始化
      markdown/
        goldmark.go          # Markdown 渲染适配
      repository/
        post_gorm.go         # PostRepository 的 gorm 实现
    transport/
      httpserver/
        router.go            # gin 引导与路由注册
        post_handlers.go     # 文章相关 HTTP 处理器
  README.md
  DESIGN.md                  # 本设计文档
```

各层职责说明：

- **Transport 层 (`internal/transport/httpserver`)**
  - 基于 `gin` 提供 HTTP 接入、参数解析、响应转换。
  - 只依赖 Usecase，不直接操作数据库。
  - 注册公共路由和后台接口。

- **Usecase 层 (`internal/usecase`)**
  - 封装业务用例，例如创建文章、查询文章。
  - 通过接口依赖仓储（Repository）和 Markdown 渲染器。

- **Domain 层 (`internal/domain`)**
  - 定义核心领域实体：Post、Tag、Category 等。
  - 定义状态与基本规则（例如文章状态枚举）。

- **Infra 层 (`internal/infra`)**
  - 提供与外部系统的具体实现：
    - 数据库访问（gorm + PostgreSQL）
    - Markdown 渲染（goldmark）
    - 日志（slog）

#### 3.2 架构示意图

```mermaid
flowchart TD
  Browser[Browser] --> Gin[gin_HTTP_Server]
  Gin --> Handlers[HTTP_Handlers]
  Handlers --> PostUsecase[Post_Usecase]
  PostUsecase --> PostRepo[PostRepository_Interface]
  PostUsecase --> MdRenderer[MarkdownRenderer_Interface]
  PostRepo --> GormDB[gorm.DB]
  GormDB --> Postgres[(PostgreSQL)]
  MdRenderer --> Goldmark[goldmark]
```

### 4. 技术选型

#### 4.1 HTTP 框架：gin

- 原因：
  - 社区成熟、文档丰富、生态完善。
  - 中间件机制简单，易于扩展日志、鉴权。
- 使用方式：
  - 在 `router.go` 中创建 `gin.Engine`，注册全局中间件（`Recovery`、自定义日志）。
  - 使用 `Group("/admin")` 定义后台路由。

#### 4.2 数据库：PostgreSQL

- 原因：
  - SQL 能力强，支持 JSON / 索引丰富。
  - 适用于将来扩展统计、全文检索等功能。
- 连接：
  - 在环境变量中通过 `DATABASE_DSN` 提供 DSN，例如：
    - `postgres://user:pass@localhost:5432/goblog?sslmode=disable`

#### 4.3 ORM：gorm

- 原因：
  - 提供方便的 CRUD API、迁移、关联处理。
  - 社区活跃，遇到问题容易找到资料。
- 使用策略：
  - 在 Domain 层的结构体上使用 gorm tag 以映射数据库字段。
  - 在 Infra 层集中管理 gorm 的 `*gorm.DB` 实例与仓储实现。

#### 4.4 日志：slog

- 原因：
  - Go 标准库，自带结构化日志能力。
  - 使用 JSON 输出方便后续接入日志系统。
- 策略：
  - `production` 环境使用 `JSONHandler`，其他环境使用 `TextHandler`。
  - HTTP 请求日志通过 gin 中间件统一记录。

#### 4.5 Markdown 渲染：goldmark

- 原因：
  - 标准兼容性好，扩展能力强。
  - 被广泛用于 Go 社区项目。
- 使用方式：
  - 在 Infra 中封装 `Renderer`，实现 `RenderToHTML(markdown string) (string, error)` 接口。
  - 在 Usecase 中调用，以在创建文章时生成 HTML 内容。

### 5. 数据模型设计

#### 5.1 实体与关系

- **Post（文章）**
  - `ID`：主键
  - `Title`：标题
  - `Slug`：URL 友好的唯一标识
  - `ContentMD`：Markdown 原文
  - `ContentHTML`：预渲染 HTML
  - `Summary`：摘要
  - `Status`：`draft` / `published` / `archived`
  - `PublishedAt`：发布时间（可空）
  - `CreatedAt` / `UpdatedAt`：时间戳

- **Tag（标签）**（后续扩展）
  - `ID`
  - `Name`：标签名
  - `Slug`：唯一标识

- **Category（分类）**（后续扩展）
  - `ID`
  - `Name`
  - `Slug`

#### 5.2 表结构示意（逻辑）

> 注：具体 DDL 由 gorm 迁移或手写 SQL 控制。

- `posts`
  - `id` BIGSERIAL PRIMARY KEY
  - `title` VARCHAR(255) NOT NULL
  - `slug` VARCHAR(255) UNIQUE NOT NULL
  - `content_md` TEXT NOT NULL
  - `content_html` TEXT
  - `summary` TEXT
  - `status` VARCHAR(32) NOT NULL DEFAULT 'draft'
  - `published_at` TIMESTAMPTZ NULL
  - `created_at` TIMESTAMPTZ NOT NULL
  - `updated_at` TIMESTAMPTZ NOT NULL

标签与分类相关表可在后续版本中补充。

### 6. Usecase 与仓储接口

#### 6.1 PostUsecase

文件：`internal/usecase/post.go`

- **接口依赖**
  - `PostRepository`：
    - `Create(ctx, *Post) error`
    - `FindByID(ctx, id uint) (*Post, error)`
    - `FindBySlug(ctx, slug string) (*Post, error)`
    - `ListRecent(ctx, limit int) ([]Post, error)`
  - `MarkdownRenderer`：
    - `RenderToHTML(markdown string) (string, error)`

- **主要方法**
  - `CreatePost(ctx, CreatePostInput) (*domain.Post, error)`
    - 渲染 Markdown → HTML
    - 根据 `Publish` 决定状态与发布时间
    - 调用仓储保存
  - `GetBySlug(ctx, slug string) (*domain.Post, error)`
  - `ListRecent(ctx, limit int) ([]domain.Post, error)`

#### 6.2 PostRepository（gorm 实现）

文件：`internal/infra/repository/post_gorm.go`

- 使用 `*gorm.DB` 作为底层实现。
- 在 `ListRecent` 中按 `published_at` 与 `created_at` 倒序排序，限制返回条数。

### 7. 路由与接口设计

#### 7.1 公共路由

- `GET /health`
  - 用途：健康检查。
  - 响应示例：
    ```json
    { "status": "ok" }
    ```

- `GET /`
  - 用途：获取最近文章列表。
  - 参数：
    - `limit`（可选，默认 10）
  - 响应：JSON 数组，每个元素为 Post 对象。

- `GET /posts/:slug`
  - 用途：按 slug 获取文章详情。
  - 路径参数：
    - `slug`：文章唯一标识。
  - 响应：
    - 200：Post 对象 JSON
    - 404：`{"error": "post not found"}`（当前简化处理）

#### 7.2 后台路由

- `POST /admin/posts`
  - 用途：创建文章（暂不含鉴权，后续通过中间件补充）。
  - 请求体（JSON）：
    ```json
    {
      "title": "Post title",
      "slug": "post-slug",
      "content_md": "markdown content",
      "summary": "short summary",
      "publish": true
    }
    ```
  - 响应：
    - 201：创建成功，返回完整 Post JSON
    - 400：参数校验失败
    - 500：服务端错误

### 8. 启动与配置

#### 8.1 配置加载

文件：`internal/config/config.go`

- 环境变量：
  - 必填：
    - `DATABASE_DSN`：PostgreSQL 连接字符串
  - 选填：
    - `APP_ENV`：`production` / `development` / `test`，默认 `development`
    - `HTTP_PORT`：监听端口，默认 `8080`

#### 8.2 启动流程

文件：`cmd/blog/main.go`

1. 加载配置。
2. 初始化 slog 日志。
3. 使用 gorm 连接 PostgreSQL。
4. 调用 `httpserver.NewRouter` 创建 gin Engine。
5. 调用 `Run(":"+HTTP_PORT)` 启动 HTTP 服务。

### 9. 部署与后续规划

#### 9.1 部署建议

- 最小部署拓扑：
  - 反向代理（例如 Caddy / Nginx）终止 TLS。
  - 后端部署 GoBlog 服务（可单实例）。
  - PostgreSQL 数据库实例（本地或托管服务）。

#### 9.2 后续演进方向

- 前台改为 SSR 模板渲染（`html/template`），支持主题与 SEO 友好的页面结构。
+- 增加用户认证与后台管理界面。
- 增加标签/分类管理与筛选页面。
- 增加 RSS / Sitemap / robots.txt 与基础 SEO 支持。
- 接入简单评论系统与访问统计。

