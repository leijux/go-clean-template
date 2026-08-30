# Go Clean Template · Frontend

基于 **React 19 + TypeScript + Vite** 的前端，对接本仓库后端的 REST API。
包管理器使用 **npm**。

## 功能

- 认证：注册 / 登录 / 退出，JWT 存于 `localStorage`，自动附带 `Authorization: Bearer <token>`。
- 任务管理：创建、列表、编辑、删除、状态流转
  （`todo → in_progress → done`，进行中可回退为 `todo`）。
- 翻译：调用翻译接口、查看当前用户历史记录。

## 快速开始

先启动后端（见仓库根目录的 [`README.md`](../README.md)，`make compose-up-all` 或本地运行）。
后端 REST 服务默认监听 `http://localhost:8080`，API 前缀为 `/v1`。

```bash
cd frontend
npm install        # 安装依赖
npm run dev        # 启动开发服务，默认 http://localhost:5173
```

开发模式下，Vite 会把 `/v1` 开头的请求代理到后端（默认 `http://localhost:8080`）。
如需指向其他后端地址，设置环境变量 `VITE_API_TARGET`（例如 `http://localhost:9090`）。

## 其他命令

```bash
npm run build    # 类型检查 + 生产构建（输出到 dist/）
npm run preview  # 预览生产构建
npm run lint     # oxlint 静态检查
```

## 技术栈

- **React 19 + TypeScript + Vite** 构建
- **[shadcn/ui](https://ui.shadcn.com)** 组件库（Tailwind CSS v4，Radix UI 基座）
  组件源码位于 `src/components/ui/`，通过 `npx shadcn@latest add <component>` 管理
- 路径别名 `@/*` → `src/*`（`vite.config.ts` 与 `tsconfig` 均已配置）

## 目录结构

```
src/
├── api/               # 类型定义与 fetch 客户端（auth / tasks / translation）
├── auth/              # AuthContext 与路由守卫
├── components/
│   ├── ui/            # shadcn/ui 组件（Button、Card、Dialog、Badge 等）
│   └── Layout.tsx     # 布局组件
├── lib/               # cn() 工具函数
├── pages/             # 登录、注册、任务、翻译页面
├── App.tsx            # 路由配置
└── main.tsx           # 入口
```

## 对接说明

所有请求走相对路径 `/v1/...`，通过 Vite 代理转发到后端，因此无需在前端硬编码后端域名。
接口字段与后端 `internal/entity` 及 `internal/controller/restapi/v1/request` 保持一致，
详见 `src/api/types.ts`。
