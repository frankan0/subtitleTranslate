package tencent

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
)

// 腾讯云翻译限制：每秒最多5个请求
const maxRequestsPerSecond = 5

// RateLimiter 实现简单的速率限制
type RateLimiter struct {
	tokens chan struct{}
	ticker *time.Ticker
	ctx    context.Context
	cancel context.CancelFunc
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter() *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	rl := &RateLimiter{
		tokens: make(chan struct{}, maxRequestsPerSecond),
		ticker: time.NewTicker(time.Second / maxRequestsPerSecond),
		ctx:    ctx,
		cancel: cancel,
	}

	// 初始化令牌
	for i := 0; i < maxRequestsPerSecond; i++ {
		rl.tokens <- struct{}{}
	}

	// 启动令牌补充协程
	go rl.refillTokens()

	return rl
}

// refillTokens 定期补充令牌
func (rl *RateLimiter) refillTokens() {
	for {
		select {
		case <-rl.ctx.Done():
			rl.ticker.Stop()
			return
		case <-rl.ticker.C:
			select {
			case rl.tokens <- struct{}{}:
				// 成功添加令牌
			default:
				// 令牌桶已满，跳过
			}
		}
	}
}

// Wait 等待获取令牌
func (rl *RateLimiter) Wait() error {
	select {
	case <-rl.tokens:
		return nil
	case <-rl.ctx.Done():
		return fmt.Errorf("rate limiter stopped")
	}
}

// Stop 停止速率限制器
func (rl *RateLimiter) Stop() {
	rl.cancel()
}

// globalRateLimiter 全局速率限制器实例
var globalRateLimiter = NewRateLimiter()

// Cleanup 清理资源
func Cleanup() {
	if globalRateLimiter != nil {
		globalRateLimiter.Stop()
	}
}

// TranslateText 翻译单个文本
func TranslateText(text, targetLang, sourceLang, secretId, secretKey, region string) (string, error) {
	// 参数验证
	if text == "" {
		return "", fmt.Errorf("翻译文本不能为空")
	}
	if targetLang == "" {
		return "", fmt.Errorf("目标语言不能为空")
	}
	if sourceLang == "" {
		return "", fmt.Errorf("源语言不能为空")
	}

	// 等待速率限制
	if err := globalRateLimiter.Wait(); err != nil {
		return "", fmt.Errorf("等待速率限制失败: %w", err)
	}

	log.Printf("[腾讯云翻译] 正在翻译文本，长度: %d 字符，源语言: %s，目标语言: %s", len(text), sourceLang, targetLang)

	// 创建认证对象
	credential := common.NewCredential(secretId, secretKey)

	// 创建客户端配置
	cpf := profile.NewClientProfile()
	//if region != "" {
	//	cpf.HttpProfile.Endpoint = fmt.Sprintf("tmt.%s.tencentcloudapi.com", region)
	//} else {
	cpf.HttpProfile.Endpoint = "tmt.tencentcloudapi.com"
	//}

	// 创建客户端
	client, err := tmt.NewClient(credential, region, cpf)
	if err != nil {
		return "", fmt.Errorf("创建腾讯云翻译客户端失败: %w", err)
	}

	// 创建请求
	request := tmt.NewTextTranslateRequest()
	request.SourceText = common.StringPtr(text)
	request.Source = common.StringPtr(sourceLang)
	request.Target = common.StringPtr(targetLang)
	request.ProjectId = common.Int64Ptr(0)

	// 发送请求
	response, err := client.TextTranslate(request)
	if err != nil {
		if tencentErr, ok := err.(*errors.TencentCloudSDKError); ok {
			log.Printf("[腾讯云翻译] API错误，错误码: %s，错误信息: %s", tencentErr.GetCode(), tencentErr.GetMessage())
			return "", fmt.Errorf("腾讯云翻译API错误: %s - %s", tencentErr.GetCode(), tencentErr.GetMessage())
		}
		log.Printf("[腾讯云翻译] 翻译请求失败: %v", err)
		return "", fmt.Errorf("翻译请求失败: %w", err)
	}

	if response.Response.TargetText == nil {
		return "", fmt.Errorf("翻译结果为空")
	}

	log.Printf("[腾讯云翻译] 翻译成功，结果长度: %d 字符", len(*response.Response.TargetText))

	return *response.Response.TargetText, nil
}

// TranslateTexts 批量翻译文本（逐行翻译）
func TranslateTexts(texts []string, targetLang, sourceLang, secretId, secretKey, region string) ([]string, error) {
	results := make([]string, len(texts))

	for i, text := range texts {
		translated, err := TranslateText(text, targetLang, sourceLang, secretId, secretKey, region)
		if err != nil {
			return nil, fmt.Errorf("翻译第%d个文本失败: %w", i+1, err)
		}
		results[i] = translated
	}

	return results, nil
}

// TranslateMergedText 翻译合并后的文本，并返回分割后的结果
func TranslateMergedText(mergedText, targetLang, sourceLang, secretId, secretKey, region string) ([]string, error) {
	// 翻译合并后的文本
	translated, err := TranslateText(mergedText, targetLang, sourceLang, secretId, secretKey, region)
	if err != nil {
		return nil, fmt.Errorf("翻译合并文本失败: %w", err)
	}

	// 按换行符分割翻译结果
	translatedLines := strings.Split(translated, "\n")

	// 去除每行前后的空白字符
	for i := range translatedLines {
		translatedLines[i] = strings.TrimSpace(translatedLines[i])
	}

	return translatedLines, nil
}

// TranslateTextsWithSettings 使用设置翻译文本（兼容现有接口）
func TranslateTextsWithSettings(texts []string, targetLang, secretId, secretKey, sourceLang string) ([]string, error) {
	// 检查必要的参数
	if secretId == "" || secretKey == "" {
		return nil, fmt.Errorf("腾讯云API密钥未配置，请在设置中配置SecretId和SecretKey")
	}

	// 设置默认区域
	region := "ap-beijing"

	return TranslateTexts(texts, targetLang, sourceLang, secretId, secretKey, region)
}
