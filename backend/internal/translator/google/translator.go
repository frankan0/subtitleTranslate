package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TranslateRequest Google翻译请求结构
type TranslateRequest struct {
	Q      []string `json:"q"`
	Target string   `json:"target"`
	Source string   `json:"source,omitempty"`
	Format string   `json:"format,omitempty"`
}

// TranslateResponse Google翻译响应结构
type TranslateResponse struct {
	Data struct {
		Translations []struct {
			TranslatedText         string `json:"translatedText"`
			DetectedSourceLanguage string `json:"detectedSourceLanguage,omitempty"`
		} `json:"translations"`
	} `json:"data"`
}

// TranslateTexts 使用Google翻译多个文本
func TranslateTexts(texts []string, targetLanguage string, sourceLanguage ...string) ([]string, error) {
	//// 获取配置
	//cfg, err := config.GetConfig()
	//if err != nil {
	//	return nil, fmt.Errorf("获取配置失败: %w", err)
	//}
	//
	//// 检查API密钥是否配置
	//if cfg.Google.APIKey == "" {
	//	return nil, fmt.Errorf("Google API密钥未配置")
	//}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	reqBody := TranslateRequest{
		Q:      texts,
		Target: mapLanguageCode(targetLanguage),
		Format: "text",
	}

	// 如果提供了源语言且不是自动检测
	if len(sourceLanguage) > 0 && sourceLanguage[0] != "" && sourceLanguage[0] != "auto" {
		reqBody.Source = mapLanguageCode(sourceLanguage[0])
	}

	// 构建请求URL
	// apiKey := cfg.Google.APIKey
	// url := fmt.Sprintf("https://translation.googleapis.com/language/translate/v2?key=%s", apiKey)
	url := "https://translation.googleapis.com/language/translate/v2?key=demo"

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("请求翻译API失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("翻译API返回错误: %s, 响应: %s", resp.Status, string(body))
	}

	var translateResp TranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&translateResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	var results []string
	for _, translation := range translateResp.Data.Translations {
		results = append(results, translation.TranslatedText)
	}

	return results, nil
}

// mapLanguageCode 将通用语言代码映射到Google支持的语言代码
func mapLanguageCode(language string) string {
	// Google翻译支持的语言代码映射
	languageMap := map[string]string{
		"zh":    "zh",    // 中文
		"zh-CN": "zh-CN", // 简体中文
		"zh-TW": "zh-TW", // 繁体中文
		"en":    "en",    // 英语
		"ja":    "ja",    // 日语
		"ko":    "ko",    // 韩语
		"fr":    "fr",    // 法语
		"de":    "de",    // 德语
		"es":    "es",    // 西班牙语
		"it":    "it",    // 意大利语
		"ru":    "ru",    // 俄语
		"pt":    "pt",    // 葡萄牙语
		"ar":    "ar",    // 阿拉伯语
		"th":    "th",    // 泰语
		"vi":    "vi",    // 越南语
	}

	// 如果找到映射，返回映射后的代码，否则返回原始代码
	if code, ok := languageMap[language]; ok {
		return code
	}
	return language
}
