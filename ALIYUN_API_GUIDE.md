# 阿里云翻译API使用指南

## 概述

阿里云机器翻译服务提供高质量的文本翻译功能，支持多种语言之间的互译。本指南将帮助您在字幕翻译项目中配置和使用阿里云翻译API。

## API配置

### 1. 获取AccessKey

1. 登录阿里云控制台：https://account.aliyun.com
2. 进入"访问控制" -> "AccessKey管理"
3. 创建新的AccessKey，获取AccessKeyId和AccessKeySecret

### 2. 开通机器翻译服务

1. 访问阿里云机器翻译产品页：https://www.aliyun.com/product/ai/alimt
2. 点击"立即开通"，选择适合的计费方式
3. 开通成功后，即可开始使用API

### 3. 参数设置

在字幕翻译项目中配置以下参数：

- **API密钥 (ApiKey)**: 您的AccessKeyId
- **API密钥Secret (ApiSecret)**: 您的AccessKeySecret
- **API地址 (ApiUrl)**: 可留空，使用默认地址

## 使用限制

### 速率限制
- **QPS限制**: 50次/秒（如有扩展需求，请联系阿里云）

### 文本长度限制
- **单次请求**: 最大5000字符
- **超出处理**: 超过5000字符的文本会自动分割处理

### 语言支持
支持以下主要语言：
- 中文 (zh)
- 英语 (en)
- 日语 (ja)
- 韩语 (ko)
- 法语 (fr)
- 德语 (de)
- 西班牙语 (es)
- 俄语 (ru)
- 意大利语 (it)
- 葡萄牙语 (pt)

更多语言代码请参考阿里云官方文档。

## 使用示例

### 基本使用
```json
{
  "filename": "example.srt",
  "content": "Hello world\nHow are you?",
  "targetLanguage": "zh",
  "sourceLanguage": "en",
  "provider": "aliyun",
  "outputFormat": "translation_only",
  "apiKey": "your-access-key-id",
  "apiSecret": "your-access-key-secret"
}
```

### 双语字幕
```json
{
  "filename": "example.srt",
  "content": "Hello world\nHow are you?",
  "targetLanguage": "zh",
  "sourceLanguage": "en",
  "provider": "aliyun",
  "outputFormat": "original_and_translation",
  "translationPosition": "below",
  "apiKey": "your-access-key-id",
  "apiSecret": "your-access-key-secret"
}
```

## 注意事项

1. **费用控制**: 阿里云机器翻译按字符计费，请注意使用量
2. **安全保护**: 妥善保管AccessKey，避免泄露
3. **错误处理**: 系统会自动处理API错误并提供回退机制
4. **文本格式**: 支持纯文本格式，自动处理HTML标签

## 错误处理

### 常见错误码
- **10110001**: 请求超时
- **10210002**: 系统错误
- **10310003**: 原文解码失败
- **10410004**: 参数缺失
- **10510005**: 语言对不支持
- **10610006**: 语种识别失败
- **10710007**: 翻译失败
- **10810008**: 字符串过长

### 处理建议
1. 检查API密钥是否正确配置
2. 确认源语言和目标语言代码是否正确
3. 确保文本长度在限制范围内
4. 如遇速率限制，系统会自动重试

## 相关链接

- [阿里云机器翻译产品页](https://www.aliyun.com/product/ai/alimt)
- [阿里云API文档](https://help.aliyun.com/zh/machine-translation/)
- [阿里云控制台](https://home.console.aliyun.com/)
- [阿里云定价信息](https://www.aliyun.com/price/product/?spm=a2c4g.11186623.0.0.3d5c4be6x6DZfP#/alimt/detail)

## 技术支持

如遇到问题，可通过以下方式获取帮助：
- 阿里云官方文档和帮助中心
- 阿里云工单系统
- 相关技术社区和论坛