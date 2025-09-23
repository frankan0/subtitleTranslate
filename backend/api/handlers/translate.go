package handlers

import (
	"net/http"

	"github.com/frank0/subtitleTranslate/internal/services"
	"github.com/gin-gonic/gin"
)

// TranslateRequest 翻译请求结构
type TranslateRequest struct {
	Texts         []string `json:"texts" binding:"required"`
	TargetLanguage string   `json:"targetLanguage" binding:"required"`
	SourceLanguage string   `json:"sourceLanguage,omitempty"`
}

// TranslateResponse 翻译响应结构
type TranslateResponse struct {
	Success bool                 `json:"success"`
	Data    *TranslateResultData `json:"data,omitempty"`
	Error   string               `json:"error,omitempty"`
}

// TranslateResultData 翻译结果数据
type TranslateResultData struct {
	TranslatedTexts []string `json:"translatedTexts"`
}

// TranslateWithVolcengine 使用火山引擎翻译
func TranslateWithVolcengine(c *gin.Context) {
	var req TranslateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, TranslateResponse{
			Success: false,
			Error:   "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 调用翻译服务
	translatedTexts, err := services.TranslateWithVolcengine(req.Texts, req.TargetLanguage, req.SourceLanguage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, TranslateResponse{
			Success: false,
			Error:   "翻译失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TranslateResponse{
		Success: true,
		Data: &TranslateResultData{
			TranslatedTexts: translatedTexts,
		},
	})
}

// TranslateWithGoogle 使用Google翻译
func TranslateWithGoogle(c *gin.Context) {
	var req TranslateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, TranslateResponse{
			Success: false,
			Error:   "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 调用翻译服务
	translatedTexts, err := services.TranslateWithGoogle(req.Texts, req.TargetLanguage, req.SourceLanguage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, TranslateResponse{
			Success: false,
			Error:   "翻译失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TranslateResponse{
		Success: true,
		Data: &TranslateResultData{
			TranslatedTexts: translatedTexts,
		},
	})
}