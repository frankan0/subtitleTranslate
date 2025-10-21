package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/frank0/subtitleTranslate/internal/models"
	"github.com/frank0/subtitleTranslate/internal/translator/aliyun"
	"github.com/frank0/subtitleTranslate/internal/translator/google"
	"github.com/frank0/subtitleTranslate/internal/translator/tencent"
	"github.com/frank0/subtitleTranslate/internal/translator/volcengine"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// 最大并发翻译数
const maxConcurrentTranslations = 5

// TranslateWithVolcengine 使用火山引擎翻译字幕文本
func TranslateWithVolcengine(texts []string, targetLanguage string, sourceLanguage string, settings models.ApiSettings) ([]string, error) {
	// 如果文本列表为空，直接返回
	if len(texts) == 0 {
		return nil, nil
	}

	// 获取源语言参数
	srcLang := "auto"
	if len(sourceLanguage) > 0 {
		srcLang = sourceLanguage
	}

	// 设置API密钥和Secret
	if settings.ApiKey != "" {
		// 设置环境变量，供翻译器使用
		os.Setenv("VOLCENGINE_ACCESS_KEY", settings.ApiKey)
	}
	if settings.ApiSecret != "" {
		os.Setenv("VOLCENGINE_SECRET_KEY", settings.ApiSecret)
	}

	// 创建结果切片
	result := make([]string, len(texts))

	// 用于存储需要处理的文本索引和内容的映射
	type textItem struct {
		index int
		text  string
	}

	var itemsToProcess []textItem

	// 预处理：检查每个文本是否需要分割
	for i, text := range texts {
		textLength := len([]rune(text))

		if textLength > 5000 {
			// 超长文本需要分割处理
			runes := []rune(text)
			var combinedResult strings.Builder

			for j := 0; j < len(runes); j += 4000 {
				end := j + 4000
				if end > len(runes) {
					end = len(runes)
				}
				subText := string(runes[j:end])
				translated, err := volcengine.TranslateTextsWithSettings([]string{subText}, targetLanguage, settings.ApiKey, settings.ApiSecret, srcLang)
				if err != nil {
					return nil, fmt.Errorf("翻译超长文本片段失败：%w", err)
				}
				combinedResult.WriteString(translated[0])
			}

			result[i] = combinedResult.String()
		} else {
			// 正常长度的文本加入批量处理队列
			itemsToProcess = append(itemsToProcess, textItem{index: i, text: text})
		}
	}

	// 如果没有需要批量处理的文本，直接返回结果
	if len(itemsToProcess) == 0 {
		return result, nil
	}

	// 按16个一组进行批量处理
	batchSize := 16
	var eg errgroup.Group
	sem := semaphore.NewWeighted(maxConcurrentTranslations)
	ctx := context.Background()

	for i := 0; i < len(itemsToProcess); i += batchSize {
		end := i + batchSize
		if end > len(itemsToProcess) {
			end = len(itemsToProcess)
		}

		batch := itemsToProcess[i:end]

		eg.Go(func() error {
			if err := sem.Acquire(ctx, 1); err != nil {
				return fmt.Errorf("获取信号量失败：%w", err)
			}
			defer sem.Release(1)

			// 提取当前批次的文本
			var batchTexts []string
			indices := make([]int, len(batch))
			for j, item := range batch {
				batchTexts = append(batchTexts, item.text)
				indices[j] = item.index
			}

			// 批量翻译
			translated, err := volcengine.TranslateTextsWithSettings(batchTexts, targetLanguage, settings.ApiKey, settings.ApiSecret, srcLang)
			if err != nil {
				return fmt.Errorf("批量翻译失败：%w", err)
			}

			// 将结果放回到对应位置
			for j, translatedText := range translated {
				result[indices[j]] = translatedText
			}

			return nil
		})
	}

	// 等待所有批次完成
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

// TranslateWithAliyun 使用阿里云翻译字幕文本
func TranslateWithAliyun(texts []string, targetLanguage string, sourceLanguage string, settings models.ApiSettings) ([]string, error) {
	// 如果文本列表为空，直接返回
	if len(texts) == 0 {
		return nil, nil
	}

	// 获取源语言参数
	srcLang := sourceLanguage
	if srcLang == "" {
		srcLang = "auto"
	}

	// 设置API密钥
	accessKeyId := settings.ApiKey
	accessKeySecret := settings.ApiSecret
	regionId := "cn-hangzhou" // 默认区域

	// 创建结果切片
	result := make([]string, len(texts))

	// 用于存储需要处理的文本索引和内容的映射
	type textItem struct {
		index int
		text  string
	}

	var itemsToProcess []textItem

	// 预处理：检查每个文本是否需要分割
	for i, text := range texts {
		textLength := len([]rune(text))

		if textLength > 5000 {
			// 超长文本需要分割处理（阿里云单次请求限制5000字符）
			runes := []rune(text)
			var combinedResult strings.Builder

			for j := 0; j < len(runes); j += 4500 {
				end := j + 4500
				if end > len(runes) {
					end = len(runes)
				}
				subText := string(runes[j:end])
				translated, err := aliyun.TranslateTexts([]string{subText}, targetLanguage, srcLang, accessKeyId, accessKeySecret, regionId)
				if err != nil {
					return nil, fmt.Errorf("翻译超长文本片段失败：%w", err)
				}
				combinedResult.WriteString(translated[0])
			}

			result[i] = combinedResult.String()
		} else {
			// 正常长度的文本加入批量处理队列
			itemsToProcess = append(itemsToProcess, textItem{index: i, text: text})
		}
	}

	// 如果没有需要批量处理的文本，直接返回结果
	if len(itemsToProcess) == 0 {
		return result, nil
	}

	// 尝试将多个短文本合并为长文本进行批量翻译，以减少API调用次数
	// 阿里云更适合处理完整句子，所以我们可以将多个短字幕行合并
	type mergedItem struct {
		indices    []int
		mergedText string
	}

	var mergedItems []mergedItem
	var currentIndices []int
	var currentText strings.Builder
	currentLength := 0

	for _, item := range itemsToProcess {
		textLength := len([]rune(item.text))

		// 如果合并后的文本不超过4500字符，就合并（留出500字符的安全余量）
		if currentLength+textLength < 4500 && currentText.Len() > 0 {
			currentText.WriteString("\n") // 用换行符分隔不同的字幕行
			currentLength++
		}

		if currentLength+textLength < 4500 {
			currentText.WriteString(item.text)
			currentLength += textLength
			currentIndices = append(currentIndices, item.index)
		} else {
			// 当前合并项已满，保存并开始新的合并
			if currentText.Len() > 0 {
				mergedItems = append(mergedItems, mergedItem{
					indices:    currentIndices,
					mergedText: currentText.String(),
				})
			}
			// 开始新的合并
			currentIndices = []int{item.index}
			currentText.Reset()
			currentText.WriteString(item.text)
			currentLength = textLength
		}
	}

	// 处理最后一组
	if currentText.Len() > 0 {
		mergedItems = append(mergedItems, mergedItem{
			indices:    currentIndices,
			mergedText: currentText.String(),
		})
	}

	// 批量翻译合并后的文本
	for _, mergedItem := range mergedItems {
		translatedLines, err := aliyun.TranslateMergedText(mergedItem.mergedText, targetLanguage, srcLang, accessKeyId, accessKeySecret, regionId)
		if err != nil {
			return nil, fmt.Errorf("批量翻译失败：%w", err)
		}

		// 确保翻译结果的行数与原始索引数量匹配
		if len(translatedLines) != len(mergedItem.indices) {
			// 如果不匹配，尝试重新翻译单个文本
			for _, index := range mergedItem.indices {
				originalText := itemsToProcess[index].text
				singleTranslated, err := aliyun.TranslateTexts([]string{originalText}, targetLanguage, srcLang, accessKeyId, accessKeySecret, regionId)
				if err != nil {
					return nil, fmt.Errorf("回退单文本翻译失败：%w", err)
				}
				result[index] = singleTranslated[0]
			}
		} else {
			// 正确映射翻译结果
			for i, index := range mergedItem.indices {
				result[index] = translatedLines[i]
			}
		}
	}

	return result, nil
}

// TranslateWithGoogle 使用Google翻译文本
func TranslateWithGoogle(texts []string, targetLanguage string, sourceLanguage string, settings models.ApiSettings) ([]string, error) {
	srcLang := "auto"
	if len(sourceLanguage) > 0 {
		srcLang = sourceLanguage
	}
	// 设置API密钥
	if settings.ApiKey != "" {
		// 设置环境变量，供翻译器使用
		os.Setenv("GOOGLE_API_KEY", settings.ApiKey)
	}

	// 设置API URL（如果有）
	if settings.ApiUrl != "" {
		os.Setenv("GOOGLE_TRANSLATE_URL", settings.ApiUrl)
	}

	return google.TranslateTexts(texts, targetLanguage, srcLang)
}

// TranslateWithTencent 使用腾讯云翻译字幕文本
func TranslateWithTencent(texts []string, targetLanguage string, sourceLanguage string, settings models.ApiSettings) ([]string, error) {
	// 如果文本列表为空，直接返回
	if len(texts) == 0 {
		return nil, nil
	}

	// 获取源语言参数
	srcLang := sourceLanguage
	if srcLang == "" {
		srcLang = "auto"
	}

	// 设置API密钥
	secretId := settings.ApiKey
	secretKey := settings.ApiSecret
	region := "ap-beijing" // 默认区域

	// 创建结果切片
	result := make([]string, len(texts))

	// 用于存储需要处理的文本索引和内容的映射
	type textItem struct {
		index int
		text  string
	}

	var itemsToProcess []textItem

	// 预处理：检查每个文本是否需要分割
	for i, text := range texts {
		textLength := len([]rune(text))

		if textLength > 6000 {
			// 超长文本需要分割处理（腾讯云单次请求限制约6000字符）
			runes := []rune(text)
			var combinedResult strings.Builder

			for j := 0; j < len(runes); j += 5000 {
				end := j + 5000
				if end > len(runes) {
					end = len(runes)
				}
				subText := string(runes[j:end])
				translated, err := tencent.TranslateTexts([]string{subText}, targetLanguage, srcLang, secretId, secretKey, region)
				if err != nil {
					return nil, fmt.Errorf("翻译超长文本片段失败：%w", err)
				}
				combinedResult.WriteString(translated[0])
			}

			result[i] = combinedResult.String()
		} else {
			// 正常长度的文本加入批量处理队列
			itemsToProcess = append(itemsToProcess, textItem{index: i, text: text})
		}
	}

	// 如果没有需要批量处理的文本，直接返回结果
	if len(itemsToProcess) == 0 {
		return result, nil
	}

	// 尝试将多个短文本合并为长文本进行批量翻译，以减少API调用次数
	// 腾讯云更适合处理完整句子，所以我们可以将多个短字幕行合并
	type mergedItem struct {
		indices    []int
		mergedText string
	}

	var mergedItems []mergedItem
	var currentIndices []int
	var currentText strings.Builder
	currentLength := 0

	for _, item := range itemsToProcess {
		textLength := len([]rune(item.text))

		// 如果合并后的文本不超过5000字符，就合并
		if currentLength+textLength < 5500 && currentText.Len() > 0 {
			currentText.WriteString("\n") // 用换行符分隔不同的字幕行
			currentLength++
		}

		if currentLength+textLength < 5500 {
			currentText.WriteString(item.text)
			currentLength += textLength
			currentIndices = append(currentIndices, item.index)
		} else {
			// 当前合并项已满，保存并开始新的合并
			if currentText.Len() > 0 {
				mergedItems = append(mergedItems, mergedItem{
					indices:    currentIndices,
					mergedText: currentText.String(),
				})
			}
			// 开始新的合并
			currentIndices = []int{item.index}
			currentText.Reset()
			currentText.WriteString(item.text)
			currentLength = textLength
		}
	}

	// 处理最后一组
	if currentText.Len() > 0 {
		mergedItems = append(mergedItems, mergedItem{
			indices:    currentIndices,
			mergedText: currentText.String(),
		})
	}

	// 批量翻译合并后的文本
	for _, mergedItem := range mergedItems {
		translatedLines, err := tencent.TranslateMergedText(mergedItem.mergedText, targetLanguage, srcLang, secretId, secretKey, region)
		if err != nil {
			return nil, fmt.Errorf("批量翻译失败：%w", err)
		}

		// 确保翻译结果的行数与原始索引数量匹配
		if len(translatedLines) != len(mergedItem.indices) {
			// 如果不匹配，尝试重新翻译单个文本
			for _, index := range mergedItem.indices {
				originalText := itemsToProcess[index].text
				singleTranslated, err := tencent.TranslateTexts([]string{originalText}, targetLanguage, srcLang, secretId, secretKey, region)
				if err != nil {
					return nil, fmt.Errorf("回退单文本翻译失败：%w", err)
				}
				result[index] = singleTranslated[0]
			}
		} else {
			// 正确映射翻译结果
			for i, index := range mergedItem.indices {
				result[index] = translatedLines[i]
			}
		}
	}

	return result, nil
}
