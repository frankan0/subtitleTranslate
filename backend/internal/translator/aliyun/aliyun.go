package aliyun

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alimt"
	"log"
	"strings"
	"time"
)

// 阿里云翻译限制：每秒最多50个请求，单次最大5000字符
const (
	maxRequestsPerSecond    = 50
	maxCharactersPerRequest = 5000
)

// RateLimiter 实现简单的速率限制
type RateLimiter struct {
	tokens chan struct{}
	ticker *time.Ticker
	stopCh chan struct{}
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, maxRequestsPerSecond),
		ticker: time.NewTicker(time.Second / maxRequestsPerSecond),
		stopCh: make(chan struct{}),
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
		case <-rl.stopCh:
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
	case <-rl.stopCh:
		return fmt.Errorf("rate limiter stopped")
	}
}

// Stop 停止速率限制器
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

// globalRateLimiter 全局速率限制器实例
var globalRateLimiter = NewRateLimiter()

// Cleanup 清理资源
func Cleanup() {
	if globalRateLimiter != nil {
		globalRateLimiter.Stop()
	}
}

// createClient 创建阿里云翻译客户端
func createClient(accessKeyId, accessKeySecret, regionId string) (*alimt.Client, error) {
	return alimt.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
}

// TranslateText 翻译单个文本
func TranslateText(text, targetLang, sourceLang, accessKeyId, accessKeySecret, regionId string) (string, error) {
	// 参数验证
	if text == "" {
		return "", fmt.Errorf("翻译文本不能为空")
	}
	if targetLang == "" {
		return "", fmt.Errorf("目标语言不能为空")
	}
	if accessKeyId == "" || accessKeySecret == "" {
		return "", fmt.Errorf("阿里云API密钥未配置")
	}

	// 等待速率限制
	if err := globalRateLimiter.Wait(); err != nil {
		return "", fmt.Errorf("等待速率限制失败: %w", err)
	}

	log.Printf("[阿里云翻译] 正在翻译文本，长度: %d 字符，源语言: %s，目标语言: %s", len(text), sourceLang, targetLang)

	// 创建客户端
	client, err := createClient(accessKeyId, accessKeySecret, regionId)
	if err != nil {
		return "", fmt.Errorf("创建阿里云翻译客户端失败: %w", err)
	}

	// 创建翻译请求
	request := alimt.CreateTranslateGeneralRequest()
	request.Method = "POST"
	request.FormatType = "text" // 文本格式
	request.SourceLanguage = sourceLang
	request.TargetLanguage = targetLang
	request.SourceText = text
	request.Scene = "general" // 通用场景

	// 发送翻译请求
	response, err := client.TranslateGeneral(request)
	if err != nil {
		log.Printf("[阿里云翻译] 翻译请求失败: %v", err)
		return "", fmt.Errorf("翻译请求失败: %w", err)
	}

	// 检查响应
	if response == nil {
		return "", fmt.Errorf("翻译结果为空")
	}

	if response.Code != 200 {
		errorMsg := response.Message
		if errorMsg == "" {
			errorMsg = "未知错误"
		}
		return "", fmt.Errorf("阿里云翻译API错误: %s - %s", response.Code, errorMsg)
	}

	translatedText := response.Data.Translated

	log.Printf("[阿里云翻译] 翻译成功，结果长度: %d 字符", len(translatedText))
	return translatedText, nil
}

// TranslateTexts 批量翻译文本（逐行翻译）
func TranslateTexts(texts []string, targetLang, sourceLang, accessKeyId, accessKeySecret, regionId string) ([]string, error) {
	results := make([]string, len(texts))

	for i, text := range texts {
		translated, err := TranslateText(text, targetLang, sourceLang, accessKeyId, accessKeySecret, regionId)
		if err != nil {
			return nil, fmt.Errorf("翻译第%d个文本失败: %w", i+1, err)
		}
		results[i] = translated
	}

	return results, nil
}

// TranslateMergedText 翻译合并后的文本，并返回分割后的结果
func TranslateMergedText(mergedText, targetLang, sourceLang, accessKeyId, accessKeySecret, regionId string) ([]string, error) {
	// 翻译合并后的文本
	translated, err := TranslateText(mergedText, targetLang, sourceLang, accessKeyId, accessKeySecret, regionId)
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
func TranslateTextsWithSettings(texts []string, targetLang, accessKeyId, accessKeySecret, sourceLang string) ([]string, error) {
	// 检查必要的参数
	if accessKeyId == "" || accessKeySecret == "" {
		return nil, fmt.Errorf("阿里云API密钥未配置，请在设置中配置AccessKeyId和AccessKeySecret")
	}

	// 设置默认区域
	regionId := "cn-hangzhou"

	return TranslateTexts(texts, targetLang, sourceLang, accessKeyId, accessKeySecret, regionId)
}
