FROM node:18-alpine as build

WORKDIR /app

# 复制依赖文件
COPY package*.json ./

# 安装依赖
RUN npm install

# 复制源代码
COPY . .

# 创建.env文件（如果不存在）
RUN if [ ! -f .env ]; then cp .env.example .env || echo "VUE_APP_API_URL=http://localhost:8080" > .env; fi

# 构建应用
RUN npm run build

# 生产环境
FROM nginx:alpine

# 复制构建产物到nginx目录
COPY --from=build /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]