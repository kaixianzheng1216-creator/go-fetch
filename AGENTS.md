# Agent 指南

本仓库是一个 Go 后端加 React/Vite 前端的自托管网站分析应用。修改时保持变更范围小，确保生成文件同步，除非用户明确要求，不要修改 `reference/umami` 里的参考资料。

## 项目上下文

- 后端模块名：`github.com/kaixianzheng1216-creator/go-fetch`。
- 服务入口：`cmd/server`。
- HTTP 路由：`internal/server/routes.go`，使用 chi 和 Huma。
- API handlers 位于 `internal/httpapi`。
- 领域校验位于 `internal/event`、`internal/user` 和 `internal/website`。
- 数据存储层位于 `internal/store`；`internal/store/sqlc` 是生成代码。
- 跟踪脚本源码是 `internal/static/js/script.js`。
- 前端应用位于 `frontend`。
- 前端构建产物位于 `internal/static/dist`，由 Go 服务嵌入。
- `reference/umami` 仅作为参考资料，常规应用修改不要动它。

## 开发命令

除特别说明外，从仓库根目录执行：

```powershell
go test ./...
go test -race ./...
go vet ./...
make check
```

前端命令：

```powershell
npm --prefix frontend install
npm --prefix frontend run dev
npm --prefix frontend run lint
npm --prefix frontend run build
npm --prefix frontend run format:check
```

Windows 格式化脚本：

```powershell
.\scripts\format.ps1
.\scripts\check-format.ps1
```

POSIX shell 格式化脚本：

```sh
./scripts/format.sh
./scripts/check-format.sh
```

完整的 `make check` 需要本机安装 `golangci-lint`。

## 本地运行

启动 PostgreSQL：

```powershell
docker compose up -d postgres
```

启动后端：

```powershell
go run ./cmd/server
```

启动前端开发服务：

```powershell
npm --prefix frontend run dev
```

默认配置定义在 `internal/config/config.go`：

- `DATABASE_URL=postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable`
- `LISTEN_ADDR=:8080`
- `ADMIN_USERNAME=admin`
- `ADMIN_PASSWORD=change-me`

## 代码生成

修改 API 路由、请求/响应结构、sqlc 查询或数据库迁移后，运行：

```powershell
go generate ./...
npm --prefix frontend run api:generate
```

生成职责：

- `api/openapi.json` 由 Huma 路由定义生成。
- `frontend/src/lib/api-types.ts` 由 `api/openapi.json` 生成。
- `internal/store/sqlc` 由 `internal/store/query` 和 `internal/store/migrations` 生成。

不要手动编辑生成文件。`api/openapi.json` 必须完全来自生成命令，生成结果是什么就保留什么；不要为了文案或格式对它做手工修改，也不要在生成后增加额外的 OpenAPI 后处理。优先修改源 schema、路由、查询或生成配置。

## 后端约定

- Handler 保持薄层逻辑。校验放在领域包，持久化细节放在 `internal/store`。
- 新增或修改接口时保留 Huma operation 元数据，确保 OpenAPI 输出可用。
- 需要登录的接口必须走 `middleware.RequireAuth`。
- `POST /api/collect` 保持公开，并保持请求大小限制；当前限制为 256 KiB。
- Store 和 handler 路径中持续传递 `context.Context`。
- 优先使用已有领域错误和错误映射模式，避免临时拼接 HTTP 错误处理。
- 修改 collector、领域校验或 store 行为时，添加聚焦测试。

## 前端约定

- 前端技术栈是 React 19、TypeScript、React Router、TanStack Query、React Hook Form、Zod、Tailwind CSS 4、shadcn 风格 UI 组件和 openapi-fetch。
- API 调用统一通过 `frontend/src/lib/api.ts`。
- 不要手写后端 API 类型；通过生成更新 `frontend/src/lib/api-types.ts`。
- Vite 将 `@` 映射到 `frontend/src`。
- 开发环境中，Vite 会把 `/api` 和 `/script.js` 代理到 `http://localhost:8080`。
- 生产构建输出到 `internal/static/dist`。

## 数据库约定

- 数据库迁移位于 `internal/store/migrations`。
- sqlc 查询位于 `internal/store/query`。
- sqlc 配置文件是 `sqlc.yaml`。
- 迁移一旦共享后应保持追加式修改。除非用户确认数据库可丢弃，否则不要重写已有迁移。
- 修改查询或 schema 后，重新生成 sqlc 代码。

## 静态资源

- `internal/static/js/script.js` 是源码，应通过现有脚本使用 Prettier 格式化。
- `internal/static/dist` 是 `frontend` 的构建产物。前端变更影响生产资源时，运行 `npm --prefix frontend run build` 重新生成。

## 交付前验证

按变更类型选择检查：

- 仅文档变更：不需要构建，但要确认路径和命令准确。
- Go 变更：运行 `go test ./...`；如果改到行为或公开 API，再运行 `go vet ./...`。
- 前端变更：运行 `npm --prefix frontend run lint` 和 `npm --prefix frontend run build`。
- API schema 变更：运行 `go generate ./...` 和 `npm --prefix frontend run api:generate`。
- Store schema 或查询变更：运行 `go generate ./internal/store` 和相关 Go 测试。
