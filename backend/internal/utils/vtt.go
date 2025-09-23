package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/frank0/subtitleTranslate/internal/models"
)

// BuildVTT 构建VTT格式字幕内容
func BuildVTT(entries []models.SubtitleEntry, outputFormat string) string {
	var builder strings.Builder

	// VTT文件头
	builder.WriteString("WEBVTT\n\n")

	for _, entry := range entries {
		// 写入序号（可选）
		builder.WriteString(fmt.Sprintf("%d\n", entry.Index))

		// 直接使用原始的时间范围
		builder.WriteString(fmt.Sprintf("%s\n", entry.TimeRange))

		// 根据输出格式写入内容
		if outputFormat == "bilingual" {
			// 双语模式：内容已经包含原始文本和翻译文本
			builder.WriteString(fmt.Sprintf("%s\n\n", entry.Content))
		} else {
			// 单语模式：内容只包含翻译文本
			builder.WriteString(fmt.Sprintf("%s\n\n", entry.Content))
		}
	}

	return builder.String()
}

// ParseVTT 解析VTT格式字幕内容
func ParseVTT(content string) ([]models.SubtitleEntry, error) {
	var entries []models.SubtitleEntry

	// 移除BOM和WEBVTT头部
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "WEBVTT")
	content = strings.TrimSpace(content)

	// 按空行分割
	blocks := strings.Split(content, "\n\n")

	index := 1
	for _, block := range blocks {
		lines := strings.Split(strings.TrimSpace(block), "\n")
		if len(lines) >= 2 {
			// 第一行是时间范围
			timeRange := strings.TrimSpace(lines[0])

			// 验证时间范围格式
			timeRegex := regexp.MustCompile(`^(\d{2}:\d{2}:\d{2}\.\d{3})\s+-->\s+(\d{2}:\d{2}:\d{2}\.\d{3})`)
			if timeRegex.MatchString(timeRange) {
				// 剩余的是内容
				content := strings.Join(lines[1:], "\n")
				content = strings.TrimSpace(content)

				if content != "" {
					entries = append(entries, models.SubtitleEntry{
						Index:     index,
						TimeRange: timeRange,
						Content:   content,
					})
					index++
				}
			}
		}
	}

	return entries, nil
}

// SecondsToVTTTime 将秒数转换为VTT时间格式 (HH:MM:SS.mmm)
func SecondsToVTTTime(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60
	millis := int((seconds - float64(int(seconds))) * 1000)

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, secs, millis)
}
