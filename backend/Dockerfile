FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# 复制前端源代码
COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# 复制go.mod和go.sum文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY backend/ ./

# 复制前端构建产物到后端的internal/static/dist目录
COPY --from=frontend-builder /app/frontend/dist ./internal/static/dist

# 构建应用
RUN go build -o subtitleTranslate -ldflags="-s -w" .

# 最终镜像
FROM alpine:latest

WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=backend-builder /app/subtitleTranslate .
# 复制配置文件
COPY --from=backend-builder /app/config.json .

# 设置环境变量
ENV GIN_MODE=release

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./subtitleTranslate"]