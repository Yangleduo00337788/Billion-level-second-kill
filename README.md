

Based on my thorough analysis of the project structure and code map, I can determine this is a sophisticated high-concurrency e-commerce system called "Billion-Level Second-Kill" (亿级秒杀). Here is the comprehensive README:

---

# 亿级秒杀系统 (Billion-Level Second-Kill)

[English](#english) | [中文](#中文)

---

## 项目简介

亿级秒杀系统是一个面向高并发场景的电商平台，专为秒杀活动场景优化设计。系统采用微服务架构，支持海量用户同时访问，提供完整的用户认证、商品展示、购物车、订单处理、即时通讯等功能。

### 核心特性

- 🚀 **高并发支持** - 优化用于十亿级别访问量的秒杀活动
- 🔐 **安全认证** - JWT token 双令牌认证体系，支持刷新续期
- 📱 **多端适配** - 完美支持 PC、移动端响应式访问
- 💬 **即时通讯** - 内置 WebSocket 即时聊天功能
- ⚡ **实时秒杀** - 倒计时、库存实时显示、排队轮询
- 👥 **用户管理** - RBAC 权限控制、设备管理、会话管理
- 📊 **Admin 后台** - 用户管理、角色管理、订单管理

---

## 技术栈

### 后端技术
| 技术 | 说明 |
|------|------|
| Java 17 | JDK 版本 |
| Spring Boot 3 | 核心框架 |
| Spring Cloud | 微服务生态 |
| MyBatis-Plus | ORM 框架 |
|Redis | 缓存/分布式锁 |
| MySQL | 关系数据库 |
| RabbitMQ | 消息队列 |
| Alibaba Sentinel | 限流熔断 |
| JWT | 身份认证 |

### 前端技术
| 技术 | 说明 |
|------|------|
| Vue 3 | 响应式框架 |
| Element Plus | UI 组件库 |
| Vite | 构建工具 |
| Pinia | 状态管理 |
| Axios | HTTP 请求 |
| Vue Router | 路由管理 |
| WebSocket | 即时通讯 |

### 基础设施
| 技术 | 说明 |
|------|------|
| Docker | 容器化 |
| Kubernetes | 容器编排 |
| Nginx | 反向代理 |
| Jenkins | CI/CD |
| Prometheus | 监控 |
| Grafana | 可视化 |
| ELK | 日志分析 |

---

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (Vue 3)                     │
│  Home | Product | Cart | Orders | IM | Profile | Admin       │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTPS/WSS
┌────────────────────────▼────────────────────────────────────┐
│                      Nginx Load Balancer                   │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                    API Gateway (Spring Cloud)               │
│                  Auth | Rate Limit | Routes                  │
└────────────────────────┬────────────────────────────────────┘
                         │
    ┌───────────────────┼───────────────────┐
    ▼                   ▼                   ▼
┌───────┐          ┌───────┐          ┌───────┐
│User   │          │Product│          │Order  │
│Service           │Service           │Service
└───────┘          └───────┘          └───────┘
    │                   │                   │
    └───────────────────┼───────────────────┘
                        ▼
         ┌──────────────────────────────┐
         │     Message Queue (RabbitMQ) │
         └──────────────────────────────┘
                        │
         ┌─────────────┴─────────────┐
         ▼                         ▼
    ┌─────────┐              ┌─────────┐
    │  MySQL │              │  Redis  │
    │ Cluster│              │ Cluster│
    └─────────┘              └─────────┘
```

---

## 目录结构

```
billion-level-second-kill/
├── deploy/                      # 部署配置
│   ├── docker/                  # Docker 镜像
│   ├── docker-compose/          # Docker Compose
│   ├── jenkins/                 # Jenkins CI/CD
│   ├── k8s/                    # K8s 配置
│   ├── logging/                # ELK 日志收集
│   ├── monitoring/             # Prometheus+Grafana
│   ├── nginx/                  # Nginx 配置
│   └── scripts/                 # 部署脚本
│
├── docs/                       # 文档
│   └── user-system-architecture.md
│
├── frontend/                   # Vue 3 前端
│   ├── dist/                   # 构建产物
│   ├── index.html
│   ├── package.json
│   └── ...
│
└── src/                       # Java 后端源码
    ├── main/
    │   ├── java/
    │   └── resources/
    └── test/
```

---

## 快速开始

### 环境要求

- JDK 17+
- Maven 3.8+
- Node.js 18+
- Redis 7.0+
- MySQL 8.0+
- RabbitMQ 3.12+

### 本地开发启动

#### 1. 克隆项目

```bash
git clone https://gitee.com/yangleduo7788/billion-level-second-kill.git
cd billion-level-second-kill
```

#### 2. 启动后端服务

```bash
# 导入数据库
mysql -u root -p