# go-fetch frontend

这是 `go-fetch` 的 React 前端，由 Vite CLI 从零创建，并通过 shadcn CLI 初始化和下载 UI 组件。

## 技术栈

- React `19`
- TypeScript `5.9`
- Vite `8`
- React Router
- TanStack Query
- React Hook Form + Zod
- Tailwind CSS `4`
- shadcn/ui registry 组件，配置见 `components.json`
- Recharts，通过 shadcn `chart` 组件接入
- OpenAPI 类型：`openapi-typescript`
- API 调用：`openapi-fetch`

## 结构

```text
src/app/                  应用 Provider、路由、QueryClient、Error Boundary
src/components/layout/    通用页面布局
src/components/ui/        shadcn CLI 下载的 registry 组件
src/features/auth/        登录、认证查询、受保护布局
src/features/websites/    站点页面、表单 schema、查询 hooks 和业务组件
src/lib/                  API client、格式化和工具函数
```

## API 类型来源

前端不手写后端接口类型。后端通过 Huma 从 Go operation 生成 `../api/openapi.json`，前端再用 `openapi-typescript` 生成 `src/lib/api-types.ts`。

接口调用集中在 `src/lib/api.ts`，内部使用 `openapi-fetch`。业务页面只从这里拿封装后的 `api` 方法和导出的类型。

后端接口变更后，在仓库根目录运行：

```powershell
go generate ./...
npm --prefix frontend run api:generate
```

只重新生成前端类型时，在 `frontend` 目录运行：

```powershell
npm run api:generate
```

## 命令

```powershell
npm install
npm run dev
npm run lint
npm run format
npm run format:check
npm run api:generate
npm run build
npm run ui:info
```

添加 shadcn/ui 组件：

```powershell
npm run ui:add -- button card input --overwrite --yes
```

开发服务会把 `/api` 和 `/script.js` 代理到 `http://localhost:8080`。生产构建输出到 `../internal/web/dist`，由 Go `embed` 打包进服务端。
