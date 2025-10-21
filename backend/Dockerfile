FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o subtitle-translate-backend .

# 最终镜像
FROM alpine:latest

WORKDIR /app

# 复制可执行文件
COPY --from=builder /app/subtitle-translate-backend .
# 复制配置文件示例
COPY config.json.example config.json

# 设置环境变量
ENV GIN_MODE=release

EXPOSE 8080

CMD ["./subtitle-translate-backend"]