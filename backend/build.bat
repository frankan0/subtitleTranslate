@echo off
echo 开始构建一体化应用...

echo 1. 构建前端...
cd ..\frontend
call npm install
call npm run build

echo 2. 复制前端构建产物到后端...
if not exist "..\backend\internal\static\dist" mkdir "..\backend\internal\static\dist"
xcopy /E /Y "dist\*" "..\backend\internal\static\dist\"

echo 3. 构建后端可执行文件...
cd ..\backend
go build -o subtitleTranslate.exe -ldflags="-s -w" .

echo 构建完成! 可执行文件: subtitleTranslate.exe