package models

// SubtitleFile 表示一个字幕文件
type SubtitleFile struct {
	Filename string          `json:"filename"` // 文件名
	Entries  []SubtitleEntry `json:"entries"`  // 字幕条目
}

// SubtitleEntry 表示一个字幕条目
type SubtitleEntry struct {
	Index     int    `json:"index"`     // 字幕序号
	TimeRange string `json:"timeRange"` // 时间范围
	Content   string `json:"content"`   // 字幕内容
}

// TranslationResult 表示翻译结果
type TranslationResult struct {
	OriginalFilename   string `json:"originalFilename"`   // 原始文件名
	TranslatedFilename string `json:"translatedFilename"` // 翻译后的文件名
	Content            string `json:"content"`            // 翻译后的内容
}

// TranslationRequest 表示翻译请求
type TranslationRequest struct {
	Filename            string `json:"filename" binding:"required"`       // 文件名
	Content             string `json:"content" binding:"required"`        // 文件内容
	TargetLanguage      string `json:"targetLanguage" binding:"required"` // 目标语言
	SourceLanguage      string `json:"sourceLanguage,omitempty"`          // 源语言，支持腾讯云等需要明确源语言的API
	Provider            string `json:"provider" binding:"required"`       // 翻译提供商 (volcengine 或 google)
	OutputFormat        string `json:"outputFormat" binding:"required"`   // 输出格式: "translation_only" 或 "original_and_translation"
	TranslationPosition string `json:"translationPosition"`               // 翻译位置: "below" 或 "above"
}

// TranslationResponse 表示翻译响应
type TranslationResponse struct {
	Success bool               `json:"success"`         // 是否成功
	Data    *TranslationResult `json:"data,omitempty"`  // 翻译结果
	Error   string             `json:"error,omitempty"` // 错误信息
}
