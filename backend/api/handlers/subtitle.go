package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/frank0/subtitleTranslate/internal/models"
	"github.com/frank0/subtitleTranslate/internal/services"
	"github.com/frank0/subtitleTranslate/internal/subtitle"
	"github.com/frank0/subtitleTranslate/internal/utils"
	"github.com/gin-gonic/gin"
)

// TranslateSubtitle 处理字幕翻译请求
func TranslateSubtitle(c *gin.Context) {
	var req models.TranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.TranslationResponse{
			Success: false,
			Error:   "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 将API设置传递给翻译服务
	apiSettings := models.ApiSettings{
		ApiKey:    req.ApiKey,
		ApiSecret: req.ApiSecret,
		ApiUrl:    req.ApiUrl,
	}

	// 创建解析器工厂
	factory := subtitle.NewParserFactory()
	
	// 获取文件扩展名
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(req.Filename), "."))
	if ext == "" {
		c.JSON(http.StatusBadRequest, models.TranslationResponse{
			Success: false,
			Error:   "文件名必须包含扩展名",
		})
		return
	}

	// 获取合适的解析器
	parser, err := factory.GetParser(req.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.TranslationResponse{
			Success: false,
			Error:   "不支持的文件格式: " + ext,
		})
		return
	}

	// 解析字幕文件
	entries, err := parser.Parse(req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.TranslationResponse{
			Success: false,
			Error:   "解析字幕文件失败: " + err.Error(),
		})
		return
	}

	// 提取所有字幕文本
	texts := make([]string, len(entries))
	for i, entry := range entries {
		texts[i] = entry.Content
	}

	// 根据提供商选择翻译服务
	var translatedTexts []string
	var translateErr error

	switch strings.ToLower(req.Provider) {
	case "volcengine":
		translatedTexts, translateErr = services.TranslateWithVolcengine(texts, req.TargetLanguage, req.SourceLanguage, apiSettings)
	case "google":
		translatedTexts, translateErr = services.TranslateWithGoogle(texts, req.TargetLanguage, req.SourceLanguage, apiSettings)
	default:
		c.JSON(http.StatusBadRequest, models.TranslationResponse{
			Success: false,
			Error:   "不支持的翻译提供商: " + req.Provider,
		})
		return
	}

	if translateErr != nil {
		c.JSON(http.StatusInternalServerError, models.TranslationResponse{
			Success: false,
			Error:   "翻译失败: " + translateErr.Error(),
		})
		return
	}

	// 更新字幕内容
	for i, entry := range entries {
		translated := translatedTexts[i]

		switch req.OutputFormat {
		case "original_and_translation":
			if req.TranslationPosition == "above" {
				entry.Content = translated + "\n" + entry.Content
			} else {
				entry.Content = entry.Content + "\n" + translated
			}
		default: // "translation_only"
			entry.Content = translated
		}
		entries[i] = entry
	}

	// 根据文件扩展名选择构建器
	lowerFilename := strings.ToLower(req.Filename)
	var translatedContent string
	if strings.HasSuffix(lowerFilename, ".vtt") {
		translatedContent = utils.BuildVTT(entries, req.OutputFormat)
	} else if strings.HasSuffix(lowerFilename, ".ass") || strings.HasSuffix(lowerFilename, ".ssa") {
		translatedContent = utils.BuildASS(entries, req.OutputFormat)
	} else {
		translatedContent = utils.BuildSRT(entries, req.OutputFormat)
	}

	// 生成翻译后的文件名
	fileExt := filepath.Ext(req.Filename)
	fileBase := strings.TrimSuffix(req.Filename, fileExt)
	var translatedFilename string

	if req.OutputFormat == "original_and_translation" {
		// 双语字幕
		if req.TranslationPosition == "above" {
			translatedFilename = fmt.Sprintf("%s_%s_bilingual_above%s", fileBase, req.TargetLanguage, fileExt)
		} else {
			translatedFilename = fmt.Sprintf("%s_%s_bilingual%s", fileBase, req.TargetLanguage, fileExt)
		}
	} else {
		// 仅译文
		translatedFilename = fmt.Sprintf("%s_%s%s", fileBase, req.TargetLanguage, fileExt)
	}

	// 返回翻译结果
	c.JSON(http.StatusOK, models.TranslationResponse{
		Success: true,
		Data: &models.TranslationResult{
			OriginalFilename:   req.Filename,
			TranslatedFilename: translatedFilename,
			Content:            translatedContent,
		},
	})
}
