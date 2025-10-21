# 腾讯云翻译API使用指南

## 概述
本项目已集成腾讯云文本翻译API，支持将字幕文件翻译成多种语言。

## API配置

### 获取腾讯云API密钥
1. 登录[腾讯云控制台](https://console.cloud.tencent.com/)
2. 进入【访问管理】->【API密钥管理】
3. 创建或获取SecretId和SecretKey

### 配置API参数
在API设置页面配置以下参数：
- **API密钥 (SecretId)**: 您的腾讯云SecretId
- **API密钥Secret (SecretKey)**: 您的腾讯云SecretKey
- **API地址**: 可选，默认为`tmt.tencentcloudapi.com`
  - 如果您使用特定区域，可以指定如`tmt.ap-beijing.tencentcloudapi.com`

## 使用限制
- **速率限制**: 每秒最多5个请求
- **文本长度**: 单次请求最多2000个字符
- **语言支持**: 支持多种语言互译

## 支持的语言代码
- `zh`: 中文
- `en`: 英语
- `ja`: 日语
- `ko`: 韩语
- `fr`: 法语
- `de`: 德语
- `es`: 西班牙语
- `ru`: 俄语
- `it`: 意大利语
- `pt`: 葡萄牙语

## 使用示例
1. 在翻译选项中选择"腾讯云"作为翻译提供商
2. 设置源语言和目标语言
3. 上传字幕文件并开始翻译

## 注意事项
1. 确保您的腾讯云账户已开通文本翻译服务
2. API调用会产生费用，请参考[腾讯云定价](https://cloud.tencent.com/document/product/551/17238)
3. 建议设置合理的翻译速率和延迟时间，避免触发API限制
4. 对于大文件，系统会自动分批处理，遵守速率限制

## 错误处理
如果遇到翻译错误，请检查：
- API密钥是否正确配置
- 账户余额是否充足
- 是否超出API调用限制
- 源语言和目标语言代码是否正确

## 相关链接
- [腾讯云文本翻译API文档](https://cloud.tencent.com/document/api/551/15619)
- [腾讯云API Explorer](https://console.cloud.tencent.com/api/explorer?Product=tmt&Version=2018-03-21&Action=TextTranslate)
- [腾讯云定价信息](https://cloud.tencent.com/document/product/551/17238)