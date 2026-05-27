# Billion-Level Second-Kill

[English](#english) | [Chinese](#chinese)

---

## Project Overview

The Billion-Level Second-Kill system is an e-commerce platform designed for high-concurrency scenarios, optimized specifically for flash sale events. Built on a microservices architecture, the system supports massive concurrent user access and provides comprehensive features including user authentication, product display, shopping cart, order processing, and real-time messaging.

### Core Features

- 🚀 **High Concurrency Support** - Optimized for flash sales handling up to billion-level traffic
- 🔐 **Secure Authentication** - Dual-token JWT authentication system with refresh token renewal
- 📱 **Multi-Device Compatibility** - Fully responsive design supporting PC and mobile platforms
- 💬 **Real-Time Messaging** - Built-in WebSocket-based chat functionality
- ⚡ **Real-Time Flash Sales** - Countdown timers, live inventory display, and queue polling
- 👥 **User Management** - RBAC permission control, device management, and session management
- 📊 **Admin Dashboard** - User, role, and order management interfaces

---

## Technology Stack

### Backend Technologies
| Technology | Description |
|------------|-------------|
| Java 17 | JDK version |
| Spring Boot 3 | Core framework |
| Spring Cloud | Microservices ecosystem |
| MyBatis-Plus | ORM framework |
| Redis | Caching / Distributed Locking |
| MySQL | Relational Database |
| RabbitMQ | Message Queue |
| Alibaba Sentinel | Rate Limiting & Circuit Breaker |
| JWT | Authentication |

### Frontend Technologies
| Technology | Description |
|------------|-------------|
| Vue 3 | Reactive framework |
| Element Plus | UI component library |
| Vite | Build tool |
| Pinia | State management |
| Axios | HTTP client |
| Vue Router | Routing management |
| WebSocket | Real-time communication |

### Infrastructure
| Technology | Description |
|------------|-------------|
| Docker | Containerization |
| Kubernetes | Container orchestration |
| Nginx | Reverse proxy |
| Jenkins | CI/CD pipeline |
| Prometheus | Monitoring |
| Grafana | Visualization |
| ELK | Log analysis |

---

## System Architecture

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
│Service│          │Service│          │Service│
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

## Directory Structure

```
billion-level-second-kill/
├── deploy/                      # Deployment configurations
│   ├── docker/                  # Docker images
│   ├── docker-compose/          # Docker Compose files
│   ├── jenkins/                 # Jenkins CI/CD
│   ├── k8s/                     # Kubernetes configurations
│   ├── logging/                 # ELK log collection
│   ├── monitoring/              # Prometheus + Grafana
│   ├── nginx/                   # Nginx configurations
│   └── scripts/                 # Deployment scripts
│
├── docs/                        # Documentation
│   └── user-system-architecture.md
│
├── frontend/                    # Vue 3 frontend
│   ├── dist/                    # Build output
│   ├── index.html
│   ├── package.json
│   └── ...
│
└── src/                         # Java backend source code
    ├── main/
    │   ├── java/
    │   └── resources/
    └── test/
```

---

## Quick Start

### Prerequisites

- JDK 17+
- Maven 3.8+
- Node.js 18+
- Redis 7.0+
- MySQL 8.0+
- RabbitMQ 3.12+

### Local Development Setup

#### 1. Clone the Repository

```bash
git clone https://gitee.com/yangleduo7788/billion-level-second-kill.git
cd billion-level-second-kill
```

#### 2. Start the Backend Service

```bash
# Import the database
mysql -u root -p
```