# go-fetch

一个用 Go + PostgreSQL 写的轻量网站分析服务。目标很朴素：把页面访问、访客、来源和自定义事件收好、算清楚、展示出来。至于把项目包装成宇宙级数据平台，这事先放进待办列表的最后一页。

它借鉴了轻量网站分析产品的方向，但不兼容任何现有项目的 API、数据库结构或前端实现。这里是一个从零实现的 Go MVP。

## 它能做什么

- 管理员登录和持久化会话。
- 创建、查看、编辑、删除网站。
- 提供 `/script.js` 采集脚本。
- 采集 pageview、SPA 路由变化和 `window.goFetch.track(name, data)` 自定义事件。
- 查看 pageviews、visitors、visits、bounces、平均访问时长。
- 查看 path、referrer、browser、OS、device、country 和自定义事件排行。
- 自动执行 PostgreSQL 迁移。
- 提供 OpenAPI 文档和前端类型生成链路。

## 技术味道

后端坚持 Go 社区常见的“少一点魔法，多一点清楚”：

- HTTP：`net/http` + `chi`
- API 文档：`Huma`
- 数据库：PostgreSQL + `pgx`
- SQL：`sqlc`
- 迁移：`goose`
- 会话：`scs`
- 日志：`slog` + `httplog`
- 前端：React + TypeScript + Vite + Tailwind CSS + shadcn/ui

更完整的架构、目录、命令和环境变量说明在 [项目详解](docs/project-overview.md)。

## 快速启动

启动 PostgreSQL：

```powershell
docker compose up -d
```

设置环境变量：

```powershell
$env:DATABASE_URL = "postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable"
$env:ADMIN_USERNAME = "admin"
$env:ADMIN_PASSWORD = "admin123"
```

如果之前用旧数据库名启动过本地 PostgreSQL，且不需要保留本地数据，可以用 `docker compose down -v` 后重新 `docker compose up -d`。

安装前端依赖、构建前端、启动后端：

```powershell
npm --prefix frontend install
npm --prefix frontend run build
go run ./cmd/server
```

打开 `http://localhost:8080`。默认管理员账号来自上面的环境变量。

## 开发模式

后端：

```powershell
go run ./cmd/server
```

前端：

```powershell
npm --prefix frontend run dev
```

Vite 通常会输出 `http://localhost:5173`。它已经把 `/api` 和 `/script.js` 代理到后端的 `http://localhost:8080`。

## 常用入口

- Web UI：`http://localhost:8080`
- OpenAPI JSON：`http://localhost:8080/openapi.json`
- API 文档页：`http://localhost:8080/api/docs`
- 健康检查：`http://localhost:8080/healthz`
- 就绪检查：`http://localhost:8080/readyz`

## 常用命令

```powershell
go generate ./...
npm --prefix frontend run api:generate
npm --prefix frontend run lint
npm --prefix frontend run build
.\scripts\check-format.ps1
```

## Tracking Snippet

创建网站后，把站点详情页生成的 snippet 放到目标页面里：

```html
<script
  defer
  src="http://localhost:8080/script.js"
  data-website-id="{website_id}"
></script>
```

自定义事件：

```html
<button onclick="window.goFetch.track('signup', { plan: 'pro' })">注册</button>
```

## 文档

- [项目详解](docs/project-overview.md)
- [前端说明](frontend/README.md)

## 开源说明

发布到 GitHub 前，建议补上 `LICENSE`、`CONTRIBUTING.md` 和 `SECURITY.md`。这几个文件不写代码，但能让贡献者少猜很多事。

当前 Go module path 是 `github.com/kaixianzheng1216-creator/go-fetch`，与 GitHub 仓库地址保持一致。

生成产物会提交进仓库，包括 `api/openapi.json`、`frontend/src/lib/api-types.ts`、`internal/store/db/*` 和 `internal/web/dist/*`。这样普通 `go build ./cmd/server` 就能直接工作，不需要先念一段生成咒语。
