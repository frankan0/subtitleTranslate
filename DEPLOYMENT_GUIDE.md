# 字幕翻译工具部署指南

本文档提供了字幕翻译工具的多种部署方式，包括本地直接部署、可执行文件部署和Docker部署。

## 目录

- [前提条件](#前提条件)
- [本地开发环境部署](#本地开发环境部署)
- [可执行文件部署](#可执行文件部署)
- [Docker部署](#docker部署)
- [配置说明](#配置说明)
- [常见问题](#常见问题)

## 前提条件

- Go 1.18+ (后端)
- Node.js 16+ (前端)
- Docker 和 Docker Compose (Docker部署)

## 本地开发环境部署

### 后端部署

1. 进入后端目录：
   ```bash
   cd backend
   ```

2. 安装依赖：
   ```bash
   go mod download
   ```

3. 运行后端服务：
   ```bash
   go run main.go
   ```

### 前端部署

1. 进入前端目录：
   ```bash
   cd frontend
   ```

2. 安装依赖：
   ```bash
   npm install
   ```

3. 运行开发服务器：
   ```bash
   npm run dev
   ```

## 可执行文件部署

### 后端打包

1. 进入后端目录：
   ```bash
   cd backend
   ```

2. 构建可执行文件：
   
   Windows:
   ```bash
   go build -o subtitle-translate-backend.exe .
   ```
   
   Linux/macOS:
   ```bash
   go build -o subtitle-translate-backend .
   ```

3. 确保`config.json`与可执行文件在同一目录

4. 运行可执行文件：
   
   Windows:
   ```bash
   .\subtitle-translate-backend.exe
   ```
   
   Linux/macOS:
   ```bash
   ./subtitle-translate-backend
   ```

### 前端打包

1. 进入前端目录：
   ```bash
   cd frontend
   ```

2. 构建生产版本：
   ```bash
   npm run build
   ```

3. 部署`dist`目录下的文件到Web服务器（如Nginx、Apache等）

## Docker部署

### 使用Docker Compose（推荐）

1. 在项目根目录下运行：
   ```bash
   docker-compose up -d
   ```

2. 访问 http://localhost 即可使用应用

### 单独构建和运行容器

#### 后端

1. 构建后端镜像：
   ```bash
   cd backend
   docker build -t subtitle-translate-backend .
   ```

2. 运行后端容器：
   ```bash
   docker run -d -p 8080:8080 -v $(pwd)/config.json:/app/config.json subtitle-translate-backend
   ```

#### 前端

1. 构建前端镜像：
   ```bash
   cd frontend
   docker build -t subtitle-translate-frontend .
   ```

2. 运行前端容器：
   ```bash
   docker run -d -p 80:80 subtitle-translate-frontend
   ```

## 配置说明

### 前端配置（.env）

```
VUE_APP_API_URL=http://localhost:8080
```

## 常见问题

### 1. 翻译API密钥配置问题


### 2. 跨域问题

如果前后端分离部署，确保在后端配置中正确设置了CORS允许的源。

### 3. Docker部署网络问题

如果使用Docker部署时前端无法连接后端API，检查网络配置和容器间通信设置。

### 4. 性能优化

对于大型字幕文件，可以调整后端的并发处理参数以提高翻译效率。