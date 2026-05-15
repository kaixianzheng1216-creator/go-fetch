# go-fetch

`go-fetch` 是一个轻量的自托管网站分析应用。项目包含 Go 后端、PostgreSQL 数据库、浏览器跟踪脚本，以及用于管理站点和查看分析数据的 React 管理后台。

## 功能概览

- Go 后端使用 chi、Huma、pgx、goose migrations、sqlc 和 SCS sessions。
- PostgreSQL 存储用户、站点、会话、事件、事件数据和后台登录会话。
- 公开采集接口：`POST /api/collect`。
- 浏览器跟踪脚本通过 `/script.js` 提供。
- 站点管理和分析接口位于 `/api` 下，除采集接口外需要登录。
- React 19 + Vite 前端构建到 `internal/static/dist`，由 Go 使用 `embed` 打包。
- 后端 Huma 路由生成 OpenAPI，再由前端生成 TypeScript API 类型。

## 目录结构

```text
api/                       生成的 OpenAPI 文档和生成入口
cmd/openapi/               OpenAPI JSON 生成器
cmd/server/                应用服务器入口
frontend/                  React/Vite 管理后台
internal/collector/        跟踪数据解析、规范化和事件输入构建
internal/config/           基于环境变量的配置
internal/event/            事件领域模型和校验
internal/httpapi/          Huma HTTP handlers 和路由注册
internal/middleware/       HTTP 与认证中间件
internal/server/           路由、OpenAPI 配置、静态资源服务
internal/session/          Session manager 配置
internal/static/js/        浏览器跟踪脚本源码
internal/static/dist/      前端构建产物，由 Go 服务嵌入
internal/store/            数据存储层、迁移、sqlc 查询和生成代码
internal/user/             用户领域模型和校验
internal/website/          站点领域模型和校验
reference/umami/           参考实现资料，不属于当前应用
scripts/                   跨平台格式化和格式检查脚本
```

## 环境要求

- Go `1.25.7`，或兼容 `go.mod` 中声明版本的 Go 工具链。
- Node.js 和 npm，用于前端开发。
- Docker，用于运行项目自带的 PostgreSQL 服务。
- `golangci-lint`，用于完整执行 `make check`。

## 快速启动

启动 PostgreSQL：

```powershell
docker compose up -d postgres
```

安装前端依赖：

```powershell
npm --prefix frontend install
```

启动 Go 后端：

```powershell
go run ./cmd/server
```

另开一个终端启动 Vite 开发服务：

```powershell
npm --prefix frontend run dev
```

后端默认监听 `http://localhost:8080`。Vite 开发服务会把 `/api` 和 `/script.js` 代理到该后端。

默认管理员账号：

```text
username: admin
password: change-me
```

非本地环境运行前，请通过环境变量设置 `ADMIN_USERNAME` 和 `ADMIN_PASSWORD`。

## 配置项

配置从环境变量读取：

| 变量 | 默认值 |
| --- | --- |
| `DATABASE_URL` | `postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable` |
| `LISTEN_ADDR` | `:8080` |
| `ADMIN_USERNAME` | `admin` |
| `ADMIN_PASSWORD` | `change-me` |

服务启动时会自动执行数据库迁移，并确保配置中的管理员用户存在。

## 跟踪脚本

在管理后台创建站点后，把跟踪脚本加入目标网站：

```html
<script
  defer
  src="http://localhost:8080/script.js"
  data-host-url="http://localhost:8080"
  data-website-id="YOUR_WEBSITE_ID"
></script>
```

自定义事件示例：

```html
<script>
  window.goFetch.track("signup", { plan: "pro" })
</script>
```

生产环境中，将 `src` 和 `data-host-url` 都替换为部署后的 `go-fetch` 地址。如果省略 `data-host-url`，事件会发送到当前页面来源下的 `/api/collect`。

## 常用命令

```powershell
# 运行全部 Go 测试
go test ./...

# 运行带 race 检查的 Go 测试
go test -race ./...

# 运行 Go vet
go vet ./...

# 运行前端 lint 和构建
npm --prefix frontend run lint
npm --prefix frontend run build

# 运行完整检查流程
make check
```

Windows 环境可使用：

```powershell
.\scripts\format.ps1
.\scripts\check-format.ps1
```

POSIX shell 环境可使用：

```sh
./scripts/format.sh
./scripts/check-format.sh
```

## 生成文件

后端路由或接口结构变更后，需要更新 OpenAPI 和前端类型：

```powershell
go generate ./...
npm --prefix frontend run api:generate
```

数据库查询或迁移变更后，需要更新 sqlc 生成代码：

```powershell
go generate ./internal/store
```

前端生产构建会输出到 `internal/static/dist`：

```powershell
npm --prefix frontend run build
```

## API 文档

服务运行后，Huma API 文档地址为：

```text
http://localhost:8080/api/docs
```

