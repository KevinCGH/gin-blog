# gin-blog

基于 Gin + Gorm 的博客系统后端，支持用户认证、文章管理、简单的评论功能。

## 项目介绍

### 技术栈

- **后端服务：** Gin + Gorm + PostgreSQL
- **部署：** Docker Compose

### 系统功能

- **用户系统：** 支持JWT 认证与授权
- **文章管理：** 增删改查
- **评论**

### 目录结构

```
├── app                 # 业务应用
│   ├── handler         # HTTP 处理函数
│   ├── models          # 数据模型
│   └── routers         # 路由管理
├── cmd                 # 入口文件
│   └── main.go         # 主程序入口
├── config              # 配置相关
├── internal            # 内部应用逻辑
│   ├── bootstrap       # 项目启动初始化代码段
│   │   └── init.go     
│   ├── global          # 全局配置
│   ├── middleware      # 中间件
│   └── utils           # 工具函数
├── config.yml          # 配置文件
├── go.mod              # Go 依赖管理文件
├── go.sum              # 依赖校验文件
├── docker-compose.yml
├── Dockerfile
└── README.ms           # 说明文档
```

### 环境要求

- Go 1.20+
- MySQL 8.0+ / PostgreSQL 17 / SQLite 3
- Docker & Docker Compose

### 实现功能
- [x] 数据库设计及模型定义
- [x] 用户认证与授权、注册和登录
- 文章管理功能
  - [x] 创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
  - [x] 读取功能，支持获取所有文章列表和单个文章的详细信息。
  - [x] 更新功能，只有文章的作者才能更新自己的文章。
  - [x] 删除功能，只有文章的作者才能删除自己的文章。
- 评论功能
  - [x] 实现评论的创建功能，已认证的用户可以对文章发表评论。
  - [x] 实现评论的读取功能，支持获取某篇文章的所有评论列表。

## 快速开始

### 启动项目

修改配置

拷贝`config.example.yml`为`config.yml`, 并修改对应的配置项

```bash
# 安装依赖
go mod tidy

# 用 docker compose 启动相关环境PostgreSQL
docker compose up -d

# 运行主程序
go run cmd/main.go
```



