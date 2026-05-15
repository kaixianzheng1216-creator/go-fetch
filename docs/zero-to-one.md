# 从 0 到 1 学会 go-fetch 后端

这份教程只讲后端。读完后，你应该能完成三件事：

- 在本地启动 PostgreSQL 和 Go 服务。
- 看懂一次浏览器采集请求如何进入后端、被清洗、写入数据库。
- 知道新增后端接口、修改数据库结构、更新 sqlc 和 OpenAPI 时应该怎么做。

## 1. 这个后端负责什么

`go-fetch` 后端是一个自托管网站分析服务，主要职责是：

1. 启动 HTTP 服务。
2. 连接 PostgreSQL。
3. 执行数据库迁移。
4. 确保管理员用户存在。
5. 提供登录、登出、当前用户接口。
6. 提供站点增删改查接口。
7. 接收公开的事件采集请求。
8. 将采集请求清洗成结构化事件并写入数据库。
9. 提供统计分析接口。
10. 根据 Huma 路由生成 OpenAPI 文档。

后端主线可以概括为：

```text
cmd/server -> config -> migrate -> store -> routes -> handlers -> domain -> database
```

## 2. 本地后端环境

需要安装：

- Go，版本以 `go.mod` 为准，当前是 `1.25.7`。
- Docker，用于启动本地 PostgreSQL。
- 可选：`golangci-lint`，用于完整执行 `make check`。

进入项目根目录：

```powershell
cd D:\workspace\workspace4goland\go-fetch
```

启动 PostgreSQL：

```powershell
docker compose up -d postgres
```

启动 Go 后端：

```powershell
go run ./cmd/server
```

默认后端地址：

```text
http://localhost:8080
```

默认管理员账号：

```text
username: admin
password: change-me
```

非本地环境不要使用默认密码。

## 3. 后端目录结构

后端相关目录如下：

```text
api/                       生成的 OpenAPI 文档和生成入口
cmd/openapi/               OpenAPI JSON 生成器
cmd/server/                Go 服务启动入口
internal/collector/        采集请求清洗和事件输入构建
internal/config/           环境变量配置
internal/event/            事件领域模型和校验
internal/httpapi/          Huma HTTP handlers 和路由注册
internal/middleware/       HTTP 中间件和认证中间件
internal/server/           路由、OpenAPI、静态资源服务
internal/session/          后台登录 session 配置
internal/static/js/        浏览器跟踪脚本源码
internal/store/            数据库访问层
internal/store/migrations/ 数据库迁移
internal/store/query/      sqlc SQL 查询
internal/store/sqlc/       sqlc 生成代码
internal/user/             用户领域模型和校验
internal/website/          站点领域模型和校验
sqlc.yaml                  sqlc 配置
docker-compose.yml         本地 PostgreSQL 服务
```

`reference/umami` 是参考资料，不属于当前后端实现。正常开发不要修改它。

## 4. 后端启动流程

入口文件：

```text
cmd/server/main.go
```

启动流程：

1. `config.Load()` 读取环境变量。
2. `store.Migrate(ctx, cfg.DatabaseURL)` 执行数据库迁移。
3. `store.Open(ctx, cfg.DatabaseURL)` 打开数据库连接。
4. `db.EnsureAdmin(ctx, cfg.AdminUsername, cfg.AdminPassword)` 确保管理员用户存在。
5. `server.New(db)` 创建应用对象。
6. `app.Routes()` 注册 HTTP 路由。
7. `http.Server.ListenAndServe()` 启动服务。

关键代码位置：

```text
cmd/server/main.go
internal/config/config.go
internal/store/migrate.go
internal/server/server.go
internal/server/routes.go
```

## 5. 配置项

配置定义在：

```text
internal/config/config.go
```

当前支持的环境变量：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `DATABASE_URL` | `postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable` | PostgreSQL 连接串 |
| `LISTEN_ADDR` | `:8080` | HTTP 监听地址 |
| `ADMIN_USERNAME` | `admin` | 初始管理员用户名 |
| `ADMIN_PASSWORD` | `change-me` | 初始管理员密码 |

配置加载后会检查空值。任何必填配置为空，服务都会启动失败。

## 6. HTTP 路由如何注册

路由入口：

```text
internal/server/routes.go
```

`Routes()` 中做了几件事：

1. 创建 chi router。
2. 挂载通用 HTTP middleware。
3. 创建 Huma API。
4. 注册 API 路由。
5. 注册静态资源和 SPA 兜底路由。

核心代码路径：

```text
internal/server/routes.go
internal/server/openapi.go
internal/httpapi/auth/routes.go
internal/httpapi/events/routes.go
internal/httpapi/websites/routes.go
internal/httpapi/analytics/routes.go
```

当前 API 大致分为：

- `auth`：登录、登出、当前用户。
- `events`：公开事件采集。
- `websites`：站点管理，需要登录。
- `analytics`：站点统计分析，需要登录。

需要登录的接口通过：

```text
middleware.RequireAuth
```

公开采集接口 `POST /api/collect` 不需要登录。

## 7. 认证和 session

认证相关目录：

```text
internal/httpapi/auth/
internal/middleware/auth.go
internal/session/session.go
```

后端使用 SCS 管理登录 session。session 数据存储在 PostgreSQL 的：

```text
app_sessions
```

登录流程：

1. 用户请求 `POST /api/login`。
2. handler 校验用户名和密码。
3. 密码验证通过后，把用户 ID 写入 session。
4. 后续需要登录的接口通过 cookie 找回当前用户。

登出流程：

1. 用户请求 `POST /api/logout`。
2. 后端销毁当前 session。

当前用户接口：

```text
GET /api/me
```

用于判断当前请求是否已登录。

## 8. 数据库结构

初始迁移文件：

```text
internal/store/migrations/00001_init.sql
```

主要表：

| 表 | 作用 |
| --- | --- |
| `users` | 后台用户 |
| `websites` | 用户创建的站点 |
| `sessions` | 访客会话 |
| `events` | 页面访问和自定义事件 |
| `event_data` | 自定义事件携带的数据 |
| `app_sessions` | 后台登录 session |

服务启动时会自动执行迁移：

```text
store.Migrate(ctx, cfg.DatabaseURL)
```

迁移原则：

- 已共享的迁移尽量不要改。
- 新增字段或表时，新建迁移文件。
- 只有确认数据库可以丢弃重建时，才考虑改已有迁移。

## 9. Store 层和 sqlc

Store 层目录：

```text
internal/store/
```

SQL 查询写在：

```text
internal/store/query/store.sql
```

sqlc 配置：

```text
sqlc.yaml
```

sqlc 生成代码：

```text
internal/store/sqlc/
```

生成命令：

```powershell
go generate ./internal/store
```

规则：

- 不要手动修改 `internal/store/sqlc`。
- 修改 SQL 查询或数据库 schema 后，重新生成。
- handler 不直接拼 SQL，优先通过 `internal/store` 提供的方法访问数据库。

## 10. 一次事件采集的完整后端链路

这是后端最核心的业务链路。

### 10.1 跟踪脚本发起请求

脚本源码：

```text
internal/static/js/script.js
```

目标网站加载脚本后，脚本会发送：

```text
POST /api/collect
```

请求体中包含：

- website ID。
- 当前 URL。
- referrer。
- 页面标题。
- 屏幕尺寸。
- 浏览器语言。
- 自定义事件名。
- 自定义事件数据。

### 10.2 路由接收请求

路由注册：

```text
internal/httpapi/events/routes.go
```

handler：

```text
internal/httpapi/events/handler.go
```

该接口是公开接口，但限制请求体最大 256 KiB。

### 10.3 清洗采集数据

核心文件：

```text
internal/collector/collector.go
```

`BuildEventInput` 会处理：

- URL 解析和兜底。
- referrer 解析。
- UTM 参数提取。
- User-Agent 解析。
- 设备类型判断。
- session ID 生成。
- visit ID 生成。
- 字段长度截断。
- 自定义事件数据传递。

这一步的目标是把外部输入变成内部稳定的事件结构。

### 10.4 写入数据库

事件最终通过 store 层写入：

- `sessions`
- `events`
- `event_data`

你阅读时可以从 handler 进入，再跟到 store 方法和 sqlc 查询。

推荐阅读顺序：

```text
internal/httpapi/events/handler.go
internal/collector/collector.go
internal/event/model.go
internal/event/rules.go
internal/store/events.go
internal/store/query/store.sql
```

## 11. 分析接口如何工作

分析接口目录：

```text
internal/httpapi/analytics/
```

它提供站点统计、趋势和指标查询。

路由注册：

```text
internal/httpapi/analytics/routes.go
```

store 查询：

```text
internal/store/analytics.go
internal/store/query/store.sql
```

分析接口通常需要：

1. 从 session 中确认当前用户。
2. 确认目标 website 属于当前用户。
3. 查询 events、sessions 或 event_data。
4. 映射成 API response。

这类接口都应该走认证中间件。

## 12. OpenAPI 如何生成

后端使用 Huma 定义接口并生成 OpenAPI。

相关文件：

```text
api/generate.go
cmd/openapi/main.go
internal/server/openapi.go
```

生成命令：

```powershell
go generate ./...
```

生成结果：

```text
api/openapi.json
```

修改接口路径、请求结构、响应结构后，都应该重新生成 OpenAPI。

规则：

- OpenAPI 来源是后端 Huma route 定义。
- 不要手动改 `api/openapi.json`。
- 修改源代码后重新生成。

## 13. 从 0 到 1 新增一个后端接口

下面用一个小任务串起后端开发流程：新增健康检查接口。

目标：

```text
GET /api/health
```

响应：

```json
{
  "status": "ok"
}
```

### 13.1 新建目录

建议新增：

```text
internal/httpapi/health/
```

文件可以拆成：

```text
routes.go
handler.go
response.go
```

### 13.2 定义 response

在 `response.go` 中定义响应结构：

```go
type HealthResponse struct {
	Body struct {
		Status string `json:"status"`
	}
}
```

### 13.3 实现 handler

在 `handler.go` 中实现：

```go
type Handler struct{}

func New() Handler {
	return Handler{}
}

func (h Handler) Health(ctx context.Context, input *struct{}) (*HealthResponse, error) {
	resp := &HealthResponse{}
	resp.Body.Status = "ok"
	return resp, nil
}
```

### 13.4 注册路由

在 `routes.go` 中用 Huma 注册 `GET /api/health`。

健康检查通常不需要登录，所以不要套 `middleware.RequireAuth`。

### 13.5 接入总路由

修改：

```text
internal/server/routes.go
```

在 `registerAPIRoutes` 中创建 handler 并调用 `health.Register`。

### 13.6 生成 OpenAPI

运行：

```powershell
go generate ./...
```

确认：

```text
api/openapi.json
```

中出现 `/api/health`。

### 13.7 验证

运行：

```powershell
go test ./...
go vet ./...
```

启动服务后请求：

```powershell
Invoke-RestMethod http://localhost:8080/api/health
```

你应该看到：

```text
status
------
ok
```

## 14. 从 0 到 1 修改数据库

再看一个更真实的后端改动：给 `websites` 增加 `description` 字段。

后端修改顺序：

1. 新增迁移文件，在 `websites` 表增加 `description`。
2. 修改 `internal/store/query/store.sql` 中站点相关 SQL。
3. 运行 `go generate ./internal/store`。
4. 修改 `internal/website` 里的模型和校验。
5. 修改 `internal/httpapi/websites` 的 request/response。
6. 修改 `internal/store/mappers.go` 或相关 mapper。
7. 运行 `go generate ./...` 更新 OpenAPI。
8. 运行 `go test ./...`。
9. 视情况运行 `go vet ./...`。

这个流程覆盖了后端最常见的跨层修改。

## 15. 后端测试怎么看

当前已有测试主要覆盖：

```text
internal/collector/collector_test.go
internal/event/rules_test.go
internal/event/data_test.go
```

测试重点：

- 外部输入是否被正确清洗。
- 领域校验是否稳定。
- 事件数据类型推断是否正确。

新增后端逻辑时，优先给以下位置补测试：

- collector：输入清洗、URL、referrer、User-Agent、session/visit 规则。
- domain rules：字段长度、必填项、枚举值。
- store：复杂查询或映射逻辑。
- handler：请求/响应结构或错误映射。

## 16. 常用后端命令

启动数据库：

```powershell
docker compose up -d postgres
```

启动后端：

```powershell
go run ./cmd/server
```

运行测试：

```powershell
go test ./...
```

运行 race 测试：

```powershell
go test -race ./...
```

运行 vet：

```powershell
go vet ./...
```

生成 sqlc 和 OpenAPI：

```powershell
go generate ./...
```

只生成 store：

```powershell
go generate ./internal/store
```

格式化 Go 代码：

```powershell
gofmt -w <files>
```

使用项目脚本检查格式：

```powershell
.\scripts\check-format.ps1
```

完整检查：

```powershell
make check
```

注意：`make check` 需要安装 `golangci-lint`，并且会执行完整项目检查。如果你只看后端，通常先跑 `go test ./...` 和 `go vet ./...`。

## 17. 推荐源码阅读顺序

只看后端时，按这个顺序读：

1. `cmd/server/main.go`：服务如何启动。
2. `internal/config/config.go`：配置从哪里来。
3. `internal/server/routes.go`：路由怎么挂载。
4. `internal/server/openapi.go`：OpenAPI 怎么生成。
5. `internal/httpapi/auth`：登录和当前用户。
6. `internal/session/session.go`：session 怎么存。
7. `internal/middleware/auth.go`：认证中间件怎么拦请求。
8. `internal/httpapi/websites`：站点管理接口。
9. `internal/httpapi/events`：采集接口入口。
10. `internal/collector/collector.go`：采集数据怎么清洗。
11. `internal/event`：事件领域模型和校验。
12. `internal/store/migrations/00001_init.sql`：数据库结构。
13. `internal/store/query/store.sql`：SQL 查询。
14. `internal/store`：业务 store 方法和 mapper。
15. `internal/httpapi/analytics`：统计接口。

读完这些后，你基本就掌握了后端主链路。

## 18. 常见问题

### 后端启动失败，提示数据库连不上

确认数据库容器状态：

```powershell
docker compose ps
```

再确认 `DATABASE_URL` 是否和 `docker-compose.yml` 中的账号、密码、数据库名一致。

### 启动后登录失败

确认服务启动时没有迁移或建管理员失败。默认账号密码来自：

```text
internal/config/config.go
```

如果改过 `ADMIN_USERNAME` 或 `ADMIN_PASSWORD`，用新的环境变量登录。

### 修改 SQL 后 Go 编译失败

重新生成 sqlc：

```powershell
go generate ./internal/store
```

然后根据生成代码的类型变化修正 store 层。

### 修改接口后 OpenAPI 没变

重新生成：

```powershell
go generate ./...
```

如果仍然没变，检查接口是否真的通过 Huma 注册到了 `registerAPIRoutes`。

### 采集接口没有数据

按顺序检查：

1. 请求是否打到 `POST /api/collect`。
2. `data-website-id` 是否存在且有效。
3. User-Agent 是否被判定为 bot。
4. handler 是否返回错误。
5. 数据库中的 `sessions`、`events`、`event_data` 是否写入。

## 19. 后端开发原则

- Handler 保持薄层逻辑，复杂规则放到领域包或 store。
- 需要登录的接口必须走认证中间件。
- 公开接口要控制请求体大小和输入校验。
- 数据库访问集中在 store 层。
- 不手动修改 sqlc 生成代码。
- 不手动修改 OpenAPI 生成文件。
- 已共享迁移保持追加式修改。
- 每次修改后跑匹配范围的测试和生成命令。

到这里，你已经能从后端视角理解 `go-fetch`：服务如何启动、接口如何注册、认证如何工作、事件如何采集、数据如何入库、OpenAPI 和 sqlc 如何生成。后续开发时，先判断改动落在哪个后端层，再按生成和验证规则收尾即可。
