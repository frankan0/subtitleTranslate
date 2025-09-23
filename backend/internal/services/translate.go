package services

import (
	"context"
	"fmt"
	"github.com/frank0/subtitleTranslate/internal/translator/google"
	"github.com/frank0/subtitleTranslate/internal/translator/volcengine"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"strings"
)

// 最大并发翻译数
const maxConcurrentTranslations = 5

// TranslateWithVolcengine 使用火山引擎翻译文本
func TranslateWithVolcengine(texts []string, targetLanguage string, sourceLanguage ...string) ([]string, error) {
	// 如果文本列表为空，直接返回
	if len(texts) == 0 {
		return nil, nil
	}

	// 获取源语言参数
	srcLang := ""
	if len(sourceLanguage) > 0 {
		srcLang = sourceLanguage[0]
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
				translated, err := volcengine.TranslateTexts([]string{subText}, targetLanguage, srcLang)
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
				translated, err := volcengine.TranslateTexts(batchTexts, targetLanguage, srcLang)
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

// TranslateWithGoogle 使用Google翻译文本
func TranslateWithGoogle(texts []string, targetLanguage string, sourceLanguage ...string) ([]string, error) {
	srcLang := ""
	if len(sourceLanguage) > 0 {
		srcLang = sourceLanguage[0]
	}
	return google.TranslateTexts(texts, targetLanguage, srcLang)
}
