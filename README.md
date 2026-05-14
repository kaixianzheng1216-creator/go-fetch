# go-fetch

`go-fetch` 是一个自托管的网站访问分析系统。后端使用 Go 和 PostgreSQL，前端使用 React/Vite，并通过一个轻量级的 `/script.js` 采集脚本记录页面访问和自定义事件。

## 功能概览

- 用户登录和基于 Cookie 的后台会话
- 网站创建、编辑、删除和列表管理
- 页面访问采集：URL、标题、referrer、UTM、语言、屏幕尺寸、浏览器、系统和设备类型
- 自定义事件采集：`window.goFetch.track(name, data)`
- 统计看板：PV、访客、访问次数、跳出次数、平均访问时长
- 维度排行：路径、来源、浏览器、自定义事件
- Huma 生成 OpenAPI，前端通过 `openapi-typescript` 生成类型
- React 管理前端构建后嵌入 Go 服务

## 技术栈

后端：

- Go `1.26` toolchain
- chi
- Huma
- pgx
- sqlc
- goose
- scs session
- PostgreSQL 16

前端：

- React 19
- TypeScript
- Vite
- React Router
- TanStack Query
- openapi-fetch
- Tailwind CSS 4
- shadcn 风格组件
- Recharts

## 目录结构

```text
api/                         OpenAPI 生成入口和 openapi.json
cmd/server/                  服务启动入口
cmd/openapi/                 OpenAPI 文件生成命令
internal/auth/               密码校验
internal/collector/          采集数据解析和事件构建
internal/config/             环境变量配置
internal/domain/             业务模型和规则
internal/httpapi/            HTTP DTO 和 domain 映射
internal/server/             路由、处理器、中间件、静态资源服务
internal/store/              PostgreSQL 访问、迁移、sqlc 查询封装
internal/web/static/         追踪脚本 script.js
internal/web/dist/           前端构建产物，供 Go embed 使用
frontend/                    React 管理前端
reference/umami/             参考项目代码，不属于主应用运行路径
scripts/                     格式化和格式检查脚本
```

## 配置

本地开发默认配置已经和 `docker-compose.yml` 对齐，启动 PostgreSQL 后可以直接运行后端：

```text
DATABASE_URL=postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable
LISTEN_ADDR=:8080
ADMIN_USERNAME=admin
ADMIN_PASSWORD=change-me
PRODUCTION=false
```

需要改数据库、端口或管理员账号时，再用环境变量覆盖对应默认值。

生产环境可设置 `PRODUCTION=true`，用于切换 JSON 日志并启用 Secure Cookie；本地运行不用设置。

`ADMIN_PASSWORD` 用于第一次启动时创建初始管理员；数据库已有用户后不会重置密码。生产环境必须覆盖默认密码。

## 本地启动

启动 PostgreSQL：

```powershell
docker compose up -d postgres
```

启动 Go 后端：

```powershell
go run ./cmd/server
```

服务默认监听：

```text
http://localhost:8080
```

健康检查：

```text
GET /healthz
GET /readyz
```

API 文档：

```text
http://localhost:8080/api/docs
http://localhost:8080/openapi.json
```

## 前端开发

安装依赖并启动 Vite：

```powershell
cd frontend
npm install
npm run dev
```

Vite 开发服务会把 `/api` 和 `/script.js` 代理到 `http://localhost:8080`。开发时通常同时运行：

- Go 后端：`go run ./cmd/server`
- 前端开发服务：`npm run dev`

## 生产式本地运行

先构建前端，产物会输出到 `internal/web/dist`：

```powershell
npm --prefix frontend run build
```

再启动后端：

```powershell
go run ./cmd/server
```

访问 `http://localhost:8080` 即可打开嵌入后的管理前端。

## 追踪脚本

在需要统计的网站中加入：

```html
<script defer src="http://localhost:8080/script.js" data-website-id="YOUR_WEBSITE_ID"></script>
```

如果脚本文件和采集 API 不在同一域名，可以显式指定采集端地址：

```html
<script
  defer
  src="https://analytics.example.com/script.js"
  data-website-id="YOUR_WEBSITE_ID"
  data-host-url="https://analytics.example.com"
></script>
```

自定义事件：

```js
window.goFetch.track("signup", {
  plan: "pro",
  source: "pricing",
})
```

采集端点：

```text
POST /api/collect
```

## 常用命令

后端测试：

```powershell
go test ./api ./cmd/... ./internal/...
```

如果本地存在 `frontend/node_modules`，不建议直接使用 `go test ./...`，因为 Go 会扫描到 `node_modules` 里的 Go 文件。

前端检查：

```powershell
npm --prefix frontend run lint
npm --prefix frontend run format:check
```

格式化：

```powershell
.\scripts\format.ps1
```

检查格式：

```powershell
.\scripts\check-format.ps1
```

重新生成 OpenAPI：

```powershell
go generate ./...
```

根据 OpenAPI 重新生成前端类型：

```powershell
npm --prefix frontend run api:generate
```

数据库 SQL 变更后重新生成 sqlc 代码：

```powershell
sqlc generate
```

## 数据模型

核心表：

- `users`：后台用户
- `websites`：被追踪的网站
- `sessions`：访客会话
- `events`：页面访问和自定义事件
- `event_data`：自定义事件数据
- `app_sessions`：后台登录会话

迁移文件位于：

```text
internal/store/migrations/
```

SQL 查询定义位于：

```text
internal/store/query/store.sql
```

## 开发约定

- 后端 API 变更后，先更新 OpenAPI，再更新前端类型。
- 前端 API 调用集中放在 `frontend/src/lib/api.ts`。
- 业务查询 hook 放在对应 feature 目录中，例如 `frontend/src/features/websites/website-queries.ts`。
- 采集字段长度和标准化规则集中在 `internal/collector/collector.go`。
- 统计口径和枚举规则集中在 `internal/domain/rules.go`。
- `reference/umami` 是参考代码，不参与主应用构建和运行。
