package volcengine

import (
	"encoding/json"
	"fmt"
	"github.com/volcengine/volc-sdk-golang/base"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// TranslateRequest 火山引擎翻译请求结构
type TranslateRequest struct {
	SourceLanguage string   `json:"SourceLanguage"`
	TargetLanguage string   `json:"TargetLanguage"`
	TextList       []string `json:"TextList"`
}

// TranslateResponse 火山引擎翻译响应结构
type TranslateResponse struct {
	ResponseMetadata struct {
		RequestID string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"Region"`
		Error     string `json:"Error"`
	} `json:"ResponseMetadata"`
	TranslationList []struct {
		Translation            string `json:"Translation"`
		DetectedSourceLanguage string `json:"DetectedSourceLanguage,omitempty"`
	} `json:"TranslationList"`
}
type Req struct {
	SourceLanguage string   `json:"SourceLanguage"`
	TargetLanguage string   `json:"TargetLanguage"`
	TextList       []string `json:"TextList"`
}

const (
	kServiceVersion = "2020-06-01"
)

var (
	client     *base.Client
	clientOnce sync.Once

	ServiceInfo = &base.ServiceInfo{
		Timeout: 10 * time.Second,
		Host:    "translate.volcengineapi.com",
		Header: http.Header{
			"Accept": []string{"application/json"},
		},
		Credentials: base.Credentials{Region: base.RegionCnNorth1, Service: "translate"},
	}
	ApiInfoList = map[string]*base.ApiInfo{
		"TranslateText": {
			Method: http.MethodPost,
			Path:   "/",
			Query: url.Values{
				"Action":  []string{"TranslateText"},
				"Version": []string{kServiceVersion},
			},
		},
	}
)

// getClient 获取火山引擎客户端，根据提供的API设置创建
func getClient(accessKey, secretKey string) *base.Client {
	// 如果没有提供API密钥，则从环境变量获取
	if accessKey == "" {
		accessKey = getEnvWithDefault("VOLCENGINE_ACCESS_KEY", "")
	}
	if secretKey == "" {
		secretKey = getEnvWithDefault("VOLCENGINE_SECRET_KEY", "")
	}

	// 如果仍然没有密钥，使用默认密钥（仅开发环境）
	if accessKey == "" || secretKey == "" {
		accessKey = "AKLTMWU5ZThiYzFlNDBjNDY1NzhjOTg3ODhjMjlmNDBiMGM"
		secretKey = "WkdJeU9USmlNek5rWTJabU5HTm1ZVGxoWVdJMVl6UTJPVFkxWkdZd1pUVQ=="
	}

	// 创建新的客户端
	client := base.NewClient(ServiceInfo, ApiInfoList)
	client.SetAccessKey(accessKey)
	client.SetSecretKey(secretKey)
	return client
}

// getEnvWithDefault 获取环境变量，不存在则返回默认值
func getEnvWithDefault(key, defaultValue string) string {
	if value := getEnv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnv 获取环境变量（包装函数，便于测试）
func getEnv(key string) string {
	// 实际实现中可以替换为更复杂的环境变量获取逻辑
	return os.Getenv(key)
}

// TranslateTexts 使用火山引擎翻译多个文本，支持重试机制
func TranslateTexts(texts []string, targetLanguage string, sourceLanguage ...string) ([]string, error) {
	return TranslateTextsWithSettings(texts, targetLanguage, "", "", sourceLanguage...)
}

// TranslateTextsWithSettings 使用火山引擎翻译多个文本，支持自定义API设置，支持重试机制
func TranslateTextsWithSettings(texts []string, targetLanguage, accessKey, secretKey string, sourceLanguage ...string) ([]string, error) {
	if len(texts) == 0 {
		return []string{}, nil
	}

	client := getClient(accessKey, secretKey)

	req := Req{
		TargetLanguage: mapLanguageCode(targetLanguage),
		TextList:       texts,
	}

	// 如果提供了源语言且不是自动检测
	if len(sourceLanguage) > 0 && sourceLanguage[0] != "" && sourceLanguage[0] != "auto" {
		req.SourceLanguage = mapLanguageCode(sourceLanguage[0])
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 重试配置
	maxRetries := 3
	retryDelay := time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// 指数退避重试
			time.Sleep(retryDelay * time.Duration(attempt))
		}

		resp, code, err := client.Json("TranslateText", nil, string(body))
		if err != nil {
			lastErr = fmt.Errorf("翻译请求失败: %w", err)
			continue
		}

		if code != http.StatusOK {
			lastErr = fmt.Errorf("翻译服务返回非200状态码: %d, 响应: %s", code, string(resp))
			if code >= 500 && code < 600 {
				// 服务器错误，可以重试
				continue
			}
			return nil, lastErr
		}

		var response TranslateResponse
		if err := json.Unmarshal([]byte(resp), &response); err != nil {
			lastErr = fmt.Errorf("解析响应失败: %w", err)
			continue
		}

		if response.ResponseMetadata.Error != "" {
			lastErr = fmt.Errorf("翻译服务错误: %s", response.ResponseMetadata.Error)
			// 业务错误不重试
			return nil, lastErr
		}

		if len(response.TranslationList) != len(texts) {
			lastErr = fmt.Errorf("翻译结果数量不匹配: 请求%d条，返回%d条", len(texts), len(response.TranslationList))
			return nil, lastErr
		}

		translations := make([]string, len(response.TranslationList))
		for i, translation := range response.TranslationList {
			translations[i] = translation.Translation
		}

		return translations, nil
	}

	return nil, fmt.Errorf("翻译失败，重试%d次后仍无法完成: %w", maxRetries, lastErr)
}

// mapLanguageCode 将通用语言代码映射到火山引擎支持的语言代码
func mapLanguageCode(language string) string {
	// 火山引擎支持的语言代码映射
	languageMap := map[string]string{
		"zh":    "zh",      // 中文
		"zh-CN": "zh",      // 简体中文
		"zh-TW": "zh-Hant", // 繁体中文
		"en":    "en",      // 英语
		"ja":    "ja",      // 日语
		"ko":    "ko",      // 韩语
		"fr":    "fr",      // 法语
		"de":    "de",      // 德语
		"es":    "es",      // 西班牙语
		"it":    "it",      // 意大利语
		"ru":    "ru",      // 俄语
		"pt":    "pt",      // 葡萄牙语
		"ar":    "ar",      // 阿拉伯语
		"th":    "th",      // 泰语
		"vi":    "vi",      // 越南语
	}

	// 如果找到映射，返回映射后的代码，否则返回原始代码
	if code, ok := languageMap[language]; ok {
		return code
	}
	return language
}
