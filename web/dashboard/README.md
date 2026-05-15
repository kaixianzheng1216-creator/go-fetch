# go-fetch dashboard

React/Vite 管理后台，作为主仓库的前端应用维护。

## Commands

```powershell
npm install
npm run dev
npm run lint
npm run format
npm run format:check
npm run api:generate
npm run build
```

从仓库根目录运行时：

```powershell
npm --prefix web/dashboard run dev
npm --prefix web/dashboard run build
```

## API Types

后端通过 Huma 生成 `../../api/openapi.json`，前端通过以下命令生成 API 类型：

```powershell
npm run api:generate
```

API client 集中在 `src/lib/api.ts`。

## Build Output

生产构建输出到 `../dist`，由仓库根目录的 `web/embed.go` 嵌入 Go 服务。开发服务会把 `/api` 和 `/script.js` 代理到 `http://localhost:8080`。
