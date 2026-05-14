# 代码 Review 路线

这份文档用于从零开始 Review 当前项目。建议按「先跑通、再读链路、最后审风险」的顺序进行，不要一上来逐文件扫。

## 1. 建立基线

先确认当前代码能正常构建和测试。

```powershell
go test ./...
go vet ./...
go build -o $env:TEMP\go-fetch-build-check.exe ./cmd/server
npm --prefix frontend install
npm --prefix frontend run build
```

先读这些文件，建立项目全局印象：

- `README.md`
- `go.mod`
- `frontend/package.json`
- `docker-compose.yml`
- `.gitignore`

Review 重点：

- 技术栈是否清晰。
- 启动命令是否完整。
- 生成命令是否明确。
- 后端、前端、数据库边界是否容易理解。

## 2. 读后端启动链路

按真实启动顺序阅读：

1. `cmd/server/main.go`
2. `internal/config/config.go`
3. `internal/store/store.go`
4. `internal/store/migrate.go`
5. `internal/server/app.go`

Review 重点：

- `main.go` 是否只负责组装和启动。
- 配置项是否足够少，默认值是否合理。
- 数据库连接、迁移、管理员初始化是否顺序清楚。
- 是否存在为了“看起来工程化”而增加的间接层。

## 3. 读后端业务链路

不要按包逐个看，按请求链路看。

### 登录和会话

阅读顺序：

1. `internal/server/auth.go`
2. `internal/server/context.go`
3. `internal/store/users.go`
4. `internal/auth`

Review 重点：

- 登录失败是否统一返回，避免泄露用户是否存在。
- session 是否只保存必要字段。
- `/api/me` 是否依赖统一的登录态读取。
- 首个管理员初始化逻辑是否清楚。

### 网站 CRUD

阅读顺序：

1. `internal/server/websites.go`
2. `internal/store/websites.go`
3. `internal/store/query/store.sql`

Review 重点：

- 用户只能访问自己的 website。
- 删除是否为软删除。
- 更新是否过滤 `deleted_at is null`。
- handler 是否只处理 HTTP 层职责。

### 埋点采集

阅读顺序：

1. `internal/server/collect.go`
2. `internal/collector/collector.go`
3. `internal/store/events.go`
4. `internal/web/static/script.js`

Review 重点：

- `/api/collect` 是否校验必要字段。
- Bot UA 是否会被忽略。
- URL、referrer、UTM、UA 解析是否集中。
- `sessions` 和 `events` 的字段边界是否清晰。
- 自定义事件数据是否正确落入 `event_data`。

### 统计查询

阅读顺序：

1. `internal/server/analytics.go`
2. `internal/store/analytics.go`
3. `internal/domain/rules.go`
4. `internal/store/query/store.sql`

Review 重点：

- 时间范围默认值是否合理。
- metric 类型是否集中维护。
- pageview、visitor、visit、bounce 的定义是否一致。
- `browser/os/device/country` 是否从 `sessions` 查询，而不是从 `events` 快照字段查询。

## 4. Review 数据库设计

核心文件：

- `internal/store/migrations/00001_init.sql`
- `internal/store/query/store.sql`
- `internal/store/db`
- `sqlc.yaml`

Review 重点：

- 表字段是否真的被业务使用。
- 字段长度是否和输入截断规则一致。
- `events` 是否只保存事件事实。
- `sessions` 是否只保存会话维度。
- `event_data` 是否能表达字符串、数字、日期等自定义数据。
- 索引是否服务当前查询。
- sqlc 生成代码是否只由 SQL 驱动，不手改。

数据库表建议重点看：

- `users`：后台用户。
- `websites`：被追踪网站。
- `sessions`：访客会话维度。
- `events`：页面访问和自定义事件。
- `event_data`：自定义事件数据。
- `app_sessions`：后台登录会话。

## 5. Review API 契约

核心文件：

- `internal/httpapi/types.go`
- `internal/server/inputs.go`
- `api/openapi.json`
- `frontend/src/lib/api-types.ts`
- `frontend/src/lib/api.ts`

Review 重点：

- HTTP DTO 是否和领域模型分开。
- OpenAPI 是否由后端生成。
- 前端 API 类型是否由 OpenAPI 生成。
- API 响应结构是否稳定。
- 错误响应是否统一。

生成命令：

```powershell
go run ./cmd/openapi
npm --prefix frontend run api:generate
```

## 6. Review 前端结构

阅读顺序：

1. `frontend/src/app`
2. `frontend/src/lib/api.ts`
3. `frontend/src/features/auth`
4. `frontend/src/features/websites`
5. `frontend/src/components/ui`

Review 重点：

- 页面是否按 feature 拆分。
- API 调用是否集中在 `lib/api.ts` 和 feature query 文件中。
- 表单、错误、加载态是否统一。
- Shadcn UI 组件是否来自 CLI 引入，而不是手写假组件。
- 是否存在大单文件回潮。
- 组件是否只承担一种清晰职责。

前端验证命令：

```powershell
npm --prefix frontend run lint
npm --prefix frontend run format:check
npm --prefix frontend run build
```

## 7. Review 安全和稳定性

重点检查：

- 未登录 API 返回 401。
- 未登录页面能正确重定向。
- 登录失败不暴露用户是否存在。
- 输入长度是否不会超过数据库字段长度。
- 数据库迁移是否能在空库初始化。

## 8. Review 测试覆盖

先列出已有测试：

```powershell
rg -n "func Test" internal
```

优先补这些测试：

- 登录成功。
- 登录失败。
- 未登录访问 API。
- 网站创建、更新、删除。
- 采集事件写入 `sessions/events/event_data`。
- Bot UA 被忽略。
- `browser/os/device/country` 指标从 `sessions` 查询。
- 无效 JSON 或缺少必填字段返回 400。

## 9. Review 是否造轮子或过度设计

逐项问这些问题：

- 这个抽象是否有两个以上真实调用方？
- 这个 helper 是否只是隐藏一行标准库代码？
- 这个字段是否已经被业务使用？
- 这个配置是否真的需要用户理解？
- 这个库是否替代了成熟社区方案？
- 这个功能是否属于当前 MVP 必须能力？

重点检查：

- `internal/domain/rules.go`
- `internal/store/mappers.go`
- `internal/server/params.go`
- `internal/httpapi/types.go`
- `internal/collector/collector.go`

判断标准：

- 能删除而不影响表达，就删。
- 能用成熟库表达，就不要手写。
- 能靠数据库、sqlc、OpenAPI 生成约束，就不要复制维护。
- 业务边界可以保留，通用样板优先减少。

## 10. Review 输出模板

建议最终按下面格式输出：

```text
P0 / P1 / P2 问题
- 文件:行号
- 问题是什么
- 为什么有风险
- 建议怎么改

开放问题
- 哪些业务决策还不确定

整体评价
- 架构是否清晰
- 是否符合 Go 社区习惯
- 是否存在造轮子
- 是否存在过度设计
- 下一步优先级
```

## 推荐 Review 顺序

如果只想走最短路线，按这个顺序：

```text
README
-> cmd/server/main.go
-> internal/config/config.go
-> internal/store/store.go
-> internal/store/migrations/00001_init.sql
-> internal/store/query/store.sql
-> internal/server/auth.go
-> internal/server/collect.go
-> internal/server/analytics.go
-> internal/httpapi/types.go
-> frontend/src/lib/api.ts
-> frontend/src/features
-> tests
```
