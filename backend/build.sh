#!/bin/bash
echo "开始构建一体化应用..."

echo "1. 构建前端..."
cd ../frontend
npm install
npm run build

echo "2. 复制前端构建产物到后端..."
mkdir -p ../backend/internal/static/dist
cp -r dist/* ../backend/internal/static/dist/

echo "3. 构建后端可执行文件..."
cd ../backend
go build -o subtitleTranslate -ldflags="-s -w" .

echo "构建完成! 可执行文件: subtitleTranslate"