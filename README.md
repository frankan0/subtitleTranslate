# 字幕翻译工具 (Subtitle Translator)

一个简单易用的字幕翻译工具，支持SRT格式字幕文件的批量翻译。前端使用Vue 3 + TypeScript + Tailwind CSS构建，后端使用Go语言开发，支持火山引擎和Google翻译API。

## 功能特点

- 支持拖放上传多个SRT字幕文件
- 支持多种目标语言选择
- 支持火山引擎和Google翻译API
- 实时显示翻译进度和状态
- 支持下载翻译后的字幕文件
- 响应式设计，适配各种设备

## 技术栈

### 前端

- Vue 3
- TypeScript
- Vite
- Tailwind CSS
- Shadcn UI

### 后端

- Go
- Gin Web框架

## 快速开始

### 前端开发

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 后端开发

```bash
# 进入后端目录
cd backend

# 复制环境变量示例文件
cp .env.example .env

# 编辑.env文件，填入API密钥

# 运行后端服务
go run main.go
```

## 配置API密钥

### 火山引擎翻译API

1. 注册[火山引擎开发者账号](https://www.volcengine.com/)
2. 创建翻译服务并获取AccessKey和SecretKey
3. 在后端的`.env`文件中配置相关密钥

### Google翻译API

1. 注册[Google Cloud Platform账号](https://cloud.google.com/)
2. 创建项目并启用Cloud Translation API
3. 创建API密钥
4. 在后端的`.env`文件中配置API密钥

## 项目结构

```
├── frontend/               # 前端代码
│   ├── src/                # 源代码
│   │   ├── assets/         # 静态资源
│   │   ├── components/     # 组件
│   │   ├── services/       # 服务
│   │   ├── types/          # 类型定义
│   │   └── utils/          # 工具函数
│   ├── public/             # 公共资源
│   └── index.html          # HTML入口
│
├── backend/                # 后端代码
│   ├── api/                # API相关
│   │   ├── handlers/       # 请求处理器
│   │   ├── middleware/     # 中间件
│   │   └── routes/         # 路由定义
│   ├── config/             # 配置
│   ├── internal/           # 内部包
│   │   ├── models/         # 数据模型
│   │   ├── services/       # 服务
│   │   ├── translator/     # 翻译实现
│   │   └── utils/          # 工具函数
│   └── main.go             # 主入口
│
└── README.md               # 项目说明
```

## 许可证

MIT