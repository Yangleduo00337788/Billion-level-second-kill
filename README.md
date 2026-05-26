# Inference Engine (推理引擎)

AI 驱动的开发者社区平台，提供智能问答、内容推荐和代码分析等功能。

## 技术栈

- **后端**: Go 1.24, Gin, GORM, MySQL, Redis, Elasticsearch
- **前端**: Vue 3, TypeScript, Vite
- **AI**: OpenAI API, RAG 检索增强生成
- **存储**: MySQL, Redis, Elasticsearch, MinIO
- **部署**: Docker, Docker Compose

## 快速启动

```bash
docker compose -f docker/docker-compose.yml up -d
```

访问 http://localhost 查看前端，http://localhost:8080/api 访问后端 API。

## 本地开发

### 后端

```bash
cd backend
cp .env.example .env
go run ./cmd/server/
```

### 前端

```bash
cd frontend
npm install
npm run dev
```

## 项目结构

```
├── backend/          # Go 后端服务
│   ├── cmd/server/   # 入口
│   ├── internal/     # 内部包
│   │   ├── ai/       # AI 集成
│   │   ├── article/  # 文章模块
│   │   ├── chat/     # 聊天模块
│   │   ├── comment/  # 评论模块
│   │   ├── config/   # 配置
│   │   ├── middleware/ # 中间件
│   │   ├── router/   # 路由
│   │   └── user/     # 用户模块
│   └── pkg/          # 公共包
├── docker/           # Docker 配置
├── frontend/         # 前端应用
└── README.md
```
