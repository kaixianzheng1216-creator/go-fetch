# 项目详解

这里放更完整的项目介绍、架构说明、开发命令和维护约定。根目录 README 负责招呼人进门，这里负责把厨房、配电箱和工具柜都标清楚。

## 项目定位

`go-fetch` 是一个用 Go + PostgreSQL 实现的最小版网站分析服务。当前版本聚焦核心分析能力：登录、网站管理、采集脚本、事件采集、基础统计和 React 后台页面。

项目借鉴轻量网站分析产品的方向，但不兼容任何现有项目的 API、数据库结构或前端实现。

## 技术栈

- Go `1.25.7`，推荐工具链 `go1.26.3`
- HTTP：标准库 `net/http` + `github.com/go-chi/chi/v5`
- API 框架与文档：`github.com/danielgtaylor/huma/v2` + `humachi`
- 数据库：PostgreSQL + `github.com/jackc/pgx/v5`
- SQL 代码生成：`sqlc v1.31.1`
- 数据库迁移：`github.com/pressly/goose/v3`
- 会话：`github.com/alexedwards/scs/v2` + `pgxstore`，会话表为 `app_sessions`
- 安全头：`github.com/unrolled/secure`
- 限流：`github.com/go-chi/httprate`
- 请求日志：标准库 `log/slog` + `github.com/go-chi/httplog/v3`
- User-Agent/Bot 解析：`github.com/mileusna/useragent`
- 配置解析：`github.com/caarlos0/env/v11`
- API 类型链路：Huma 生成 OpenAPI `3.1`，前端通过 `openapi-typescript` + `openapi-fetch` 使用
- 前端：React `19` + TypeScript `5.9` + Vite + React Router + TanStack Query
- UI：Tailwind CSS `4` + shadcn/ui CLI registry 组件 + lucide-react + Recharts
- 表单：React Hook Form + Zod

## 项目结构

```text
api/openapi.json            Huma 从 Go 后端生成的 OpenAPI 契约
cmd/openapi/                OpenAPI 生成命令
cmd/server/                 服务启动、应用装配和生命周期管理
frontend/                   React + TypeScript + Tailwind + shadcn/ui 前端源码
internal/auth/              密码校验
internal/collector/         采集请求 URL/UTM/session/visit 归一化
internal/config/            环境变量配置
internal/domain/            领域模型、枚举和业务规则
internal/httpapi/           API DTO、响应类型、schema helper 和领域模型转换
internal/server/            chi + Huma 路由、API handler、SPA 托管、SCS 会话接入
internal/store/             PostgreSQL 存储封装
internal/store/db/          sqlc 生成的类型安全查询代码
internal/store/migrations/  goose 迁移文件
internal/store/query/       sqlc 查询定义
internal/web/               采集脚本和 React 构建产物的 embed 入口
scripts/                    后端和前端格式化脚本
```

## 后端组织约定

- `cmd/server/main.go` 放服务启动、依赖装配和生命周期管理。
- `internal/server/app.go` 放 HTTP 应用装配。
- `internal/server/openapi.go`、`errors.go`、`middleware.go`、`static.go` 分别处理接口文档、错误模型、中间件和静态资源。
- `internal/server/{auth,collect,websites,analytics}.go` 按业务域注册路由和放置对应 handler。
- `internal/server/inputs.go`、`params.go`、`responses.go` 放公共 Huma 输入、查询参数和响应包装。
- `internal/domain/rules.go` 集中放置采集类型、时间粒度、指标类型、默认 lookback、metric limit 和 pageview/custom event 映射。
- `internal/store` 按数据职责拆分：连接在 `store.go`，迁移在 `migrate.go`，用户、网站、事件写入和分析查询分别放在独立文件。
- 数据库行到领域模型的转换集中在 `internal/store/mappers.go`。

## MVP 范围

- 单管理员登录，使用 SCS 持久化会话。
- 网站创建、查看、编辑、删除。
- `/script.js` 采集脚本，支持 pageview、SPA 路由变化和 `window.goFetch.track(name, data)` 自定义事件。
- `/api/collect` 事件采集接口。
- 基础统计：pageviews、visitors、visits、bounces、平均访问时长。
- 趋势数据和维度排行：path、referrer、browser、OS、device、country、自定义事件。
- PostgreSQL schema 通过 goose 在服务启动时自动迁移。
- `/healthz` 和 `/readyz` 健康检查。
- 登录和采集接口使用进程内限流，当前 MVP 不引入 Redis 分布式限流。

## API 文档和类型生成

后端以 Go/Huma operation 为接口契约来源，不维护手写 OpenAPI map 或 YAML。

- 运行时 OpenAPI：`http://localhost:8080/openapi.json`
- 运行时文档页：`http://localhost:8080/api/docs`
- 提交到仓库的契约文件：`api/openapi.json`
- 前端类型文件：`frontend/src/lib/api-types.ts`

重新生成后端 OpenAPI 和前端 API 类型：

```powershell
go generate ./...
npm --prefix frontend run api:generate
```

也可以单独生成 OpenAPI：

```powershell
go run ./cmd/openapi -out api/openapi.json
```

新增或修改接口时遵循这个顺序：

1. 在 `internal/httpapi` 定义或调整请求/响应 DTO。
2. 在 `internal/server` 对应业务文件注册 Huma operation，并把业务逻辑接到 `store`、`collector` 或其他领域服务。
3. 运行 `go generate ./...` 重新生成 `api/openapi.json`。
4. 运行 `npm --prefix frontend run api:generate` 重新生成前端 API 类型。
5. 前端通过 `frontend/src/lib/api.ts` 使用 `openapi-fetch` 调用接口，不手写接口响应类型。

## 后端标准化边界

- API 路由、请求体 schema、查询参数、响应 schema 和 OpenAPI 文档由 Huma 注册生成。
- Huma 通过 `humachi` 接入现有 chi router，继续复用 chi middleware。
- 统一错误响应保持 `{"error":{"message":"..."}}`，通过 Huma 自定义 error model 实现。
- 登录、网站创建、采集等简单业务校验保留在 handler 中，因为这是业务规则，不是通用校验框架问题。
- 查询参数由 Huma 解析到 input struct，不再维护独立 form decoder。
- Huma 不支持 query/path/header/cookie 字段直接使用指针类型；需要区分“未传”和“传 0”的查询参数使用 `optionalParam[T]` 模式。
- Huma operation 是后端接口文档的单一来源；不要再新增手写 OpenAPI YAML、Swagger 注释或单独维护的接口文档。
- 采集类型、时间粒度、指标类型和默认 limit 以 `internal/domain/rules.go` 为单一来源；OpenAPI schema、handler 校验和 store 查询都应引用这里的规则。
- 请求日志使用 `httplog` + `slog`；`APP_ENV=production` 时输出 JSON 日志。
- User-Agent 和 Bot 判断使用 `mileusna/useragent`，不维护自写 UA token 表。
- 安全响应头使用 `unrolled/secure`。
- PostgreSQL 查询使用 sqlc 生成代码；UUID 通过 sqlc override 映射到 `google/uuid.UUID`。

## 本地运行

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

构建前端并启动服务：

```powershell
npm --prefix frontend install
npm --prefix frontend run build
go run ./cmd/server
```

开发时也可以前后端分开运行。后端监听 `:8080`：

```powershell
go run ./cmd/server
```

另开一个终端启动 Vite：

```powershell
npm --prefix frontend run dev
```

Vite 通常输出 `http://localhost:5173`。前端开发服务器代理 `/api` 和 `/script.js` 到 `http://localhost:8080`。

## Tracking Snippet

创建网站后，在目标网站页面中加入：

```html
<script
  defer
  src="http://localhost:8080/script.js"
  data-website-id="{website_id}"
></script>
```

生产环境中把 `src` 替换为实际部署域名。

自定义事件示例：

```html
<button onclick="window.goFetch.track('signup', { plan: 'pro' })">注册</button>
```

## 开发命令

开发前端：

```powershell
npm --prefix frontend run dev
```

构建前端：

```powershell
npm --prefix frontend run build
```

检查前端 lint：

```powershell
npm --prefix frontend run lint
```

格式化全部 Go 和前端代码：

```powershell
./scripts/format.ps1
```

检查格式化：

```powershell
./scripts/check-format.ps1
```

单独生成 sqlc：

```powershell
go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1 generate
```

生成 OpenAPI 和前端 API 类型：

```powershell
go generate ./...
npm --prefix frontend run api:generate
```

运行测试：

```powershell
go test $(go list ./... | Where-Object { $_ -notlike "*/frontend/node_modules/*" })
```

静态检查：

```powershell
$packages = go list ./... | Where-Object { $_ -notlike "*/frontend/node_modules/*" }
go vet $packages
go run honnef.co/go/tools/cmd/staticcheck@latest $packages
```

构建服务：

```powershell
go build -o go-fetch.exe ./cmd/server
```

完整验证建议：

```powershell
$packages = go list ./... | Where-Object { $_ -notlike "*/frontend/node_modules/*" }
go test $packages
go vet $packages
go run honnef.co/go/tools/cmd/staticcheck@latest $packages
npm --prefix frontend run lint
npm --prefix frontend run build
./scripts/check-format.ps1
```

## 开源前检查

- 补充 `LICENSE`，明确项目许可证。
- 补充 `CONTRIBUTING.md`，说明如何安装依赖、生成代码、跑测试和提交 PR。
- 补充 `SECURITY.md`，说明如何报告安全问题。
- 配置 CI，至少运行代码生成、Go 测试、Go vet、Staticcheck、前端 lint/build 和格式检查。
- 生成产物：`api/openapi.json`、`frontend/src/lib/api-types.ts`、`internal/store/db/*` 和 `internal/web/dist/*` 建议提交进仓库，保证普通 `go build ./cmd/server` 可用。
- 当前 Go module path 是 `github.com/kaixianzheng1216-creator/go-fetch`，与 GitHub 仓库地址保持一致。

## 常见问题

### npm install 报 ERESOLVE

前端使用 `openapi-typescript@7.x`，它要求 TypeScript `5.x`。项目已经把 `typescript` 锁到 `~5.9.3`，正常情况下直接运行：

```powershell
npm --prefix frontend install
```

不要用 `--force` 或 `--legacy-peer-deps` 绕过依赖冲突；如果仍然报错，先删除 `frontend/node_modules` 后重新安装。

### 修改接口后前端类型没有更新

后端接口契约来自 Huma operation。修改 `internal/httpapi` 或 `internal/server` 后，按顺序运行：

```powershell
go generate ./...
npm --prefix frontend run api:generate
```

然后再运行前端 lint/build。

## 环境变量

| 变量                            | 必填           | 默认值                                 | 说明                                            |
| ------------------------------- | -------------- | -------------------------------------- | ----------------------------------------------- |
| `DATABASE_URL`                  | 是             | 无                                     | PostgreSQL 连接字符串                           |
| `LISTEN_ADDR`                   | 否             | `:8080`                                | HTTP 监听地址                                   |
| `ADMIN_USERNAME`                | 否             | `admin`                                | 首次初始化管理员用户名                          |
| `ADMIN_PASSWORD`                | 首次初始化必填 | 无                                     | 首次初始化管理员密码；已有用户后重启不再需要    |
| `APP_ENV`                       | 否             | `development`                          | 运行环境；为 `production` 时默认启用安全 Cookie |
| `COOKIE_SECURE`                 | 否             | production 中 `true`，其他环境 `false` | 是否为会话 Cookie 添加 `Secure`                 |
| `HTTP_READ_TIMEOUT`             | 否             | `10s`                                  | HTTP 读取请求超时                               |
| `HTTP_WRITE_TIMEOUT`            | 否             | `30s`                                  | HTTP 写响应超时                                 |
| `HTTP_IDLE_TIMEOUT`             | 否             | `60s`                                  | HTTP keep-alive 空闲超时                        |
| `HTTP_HANDLER_TIMEOUT`          | 否             | `30s`                                  | 单个 handler 的 chi middleware 超时             |
| `HTTP_SHUTDOWN_TIMEOUT`         | 否             | `10s`                                  | 优雅关闭等待时间                                |
| `DB_MAX_CONNS`                  | 否             | `10`                                   | PostgreSQL 连接池最大连接数                     |
| `LOGIN_RATE_LIMIT_PER_MINUTE`   | 否             | `10`                                   | 每 IP 每分钟登录请求上限                        |
| `COLLECT_RATE_LIMIT_PER_MINUTE` | 否             | `120`                                  | 每 IP 每分钟采集请求上限                        |
