# 字幕翻译工具部署指南 (更新版)

本文档提供了字幕翻译工具的多种部署方式，包括本地开发环境、**一体化可执行文件部署**和Docker部署。

## 目录

- [前提条件](#前提条件)
- [本地开发环境部署](#本地开发环境部署)
- [一体化可执行文件部署](#一体化可执行文件部署)
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

## 一体化可执行文件部署

我们提供了自动化构建脚本，可以将前端静态文件嵌入到Go程序中，生成单一可执行文件：

### Windows

1. 在项目根目录下运行：
   ```bash
   cd backend
   .\build.bat
   ```

2. 构建完成后，运行生成的可执行文件：
   ```bash
   .\subtitleTranslate.exe
   ```

### Linux/macOS

1. 在项目根目录下运行：
   ```bash
   cd backend
   chmod +x build.sh
   ./build.sh
   ```

2. 构建完成后，运行生成的可执行文件：
   ```bash
   ./subtitleTranslate
   ```

应用将在 http://localhost:8080 上运行，同时提供API和前端界面。

## Docker部署

### 使用Docker Compose（推荐）

1. 确保已安装Docker和Docker Compose
2. 在项目根目录下运行：
   ```bash
   docker-compose up -d
   ```

这将启动一体化服务，可通过 http://localhost:8080 访问。

### 手动构建容器

1. 在项目根目录下运行：
   ```bash
   docker build -t subtitle-translate -f backend/Dockerfile .
   ```

2. 运行容器：
   ```bash
   docker run -p 8080:8080 -v $(pwd)/backend/config.json:/app/config.json subtitle-translate
   ```

## 配置说明

### config.json

后端配置文件`config.json`包含以下设置：

```json
{
  "server": {
    "port": 8080,
    "host": "0.0.0.0"
  },
  "translator": {
    "default": "aliyun",
    "aliyun": {
      "accessKeyId": "YOUR_ACCESS_KEY_ID",
      "accessKeySecret": "YOUR_ACCESS_KEY_SECRET",
      "region": "cn-hangzhou"
    }
  },
  "logging": {
    "level": "info",
    "file": "logs/app.log"
  }
}
```

## 常见问题

### 1. 翻译API密钥配置

确保在`config.json`中正确配置了翻译服务的API密钥。

### 2. 静态文件嵌入问题

如果遇到静态文件访问问题，请确保：
- 构建脚本正确执行，前端静态文件已成功构建
- 后端程序能够正确访问嵌入的静态文件

### 3. Docker网络问题

如果使用Docker部署时遇到网络问题，请检查容器网络设置和端口映射。