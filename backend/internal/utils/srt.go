package utils

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/frank0/subtitleTranslate/internal/models"
)

// ParseSRT 解析SRT格式的字幕文件内容
func ParseSRT(content string) ([]models.SubtitleEntry, error) {
	var entries []models.SubtitleEntry
	scanner := bufio.NewScanner(strings.NewReader(content))

	var currentEntry *models.SubtitleEntry
	var contentLines []string

	for scanner.Scan() {
		line := scanner.Text()

		// 尝试解析为索引号
		if currentEntry == nil {
			index, err := strconv.Atoi(strings.TrimSpace(line))
			if err == nil && index > 0 {
				currentEntry = &models.SubtitleEntry{Index: index}
				continue
			}
		}

		// 尝试解析为时间范围
		if currentEntry != nil && currentEntry.TimeRange == "" {
			// 时间格式: 00:00:00,000 --> 00:00:00,000
			timeRangePattern := regexp.MustCompile(`\d{2}:\d{2}:\d{2},\d{3} --> \d{2}:\d{2}:\d{2},\d{3}`)
			if timeRangePattern.MatchString(line) {
				currentEntry.TimeRange = line
				continue
			}
		}

		// 收集字幕内容
		if currentEntry != nil && currentEntry.TimeRange != "" {
			if line == "" {
				// 空行表示一个条目的结束
				if len(contentLines) > 0 {
					currentEntry.Content = strings.Join(contentLines, "\n")
					entries = append(entries, *currentEntry)
					currentEntry = nil
					contentLines = nil
				}
			} else {
				contentLines = append(contentLines, line)
			}
		}
	}

	// 处理最后一个条目
	if currentEntry != nil && currentEntry.TimeRange != "" && len(contentLines) > 0 {
		currentEntry.Content = strings.Join(contentLines, "\n")
		entries = append(entries, *currentEntry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("解析SRT文件失败: %w", err)
	}

	return entries, nil
}

// FormatSRT 将字幕条目格式化为SRT字符串
func FormatSRT(entries []models.SubtitleEntry) string {
	var builder strings.Builder

	for i, entry := range entries {
		// 添加索引
		builder.WriteString(strconv.Itoa(entry.Index))
		builder.WriteString("\n")

		// 添加时间范围
		builder.WriteString(entry.TimeRange)
		builder.WriteString("\n")

		// 添加内容
		builder.WriteString(entry.Content)
		builder.WriteString("\n")

		// 在条目之间添加空行，最后一个条目后不添加
		if i < len(entries)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// BuildSRT 构建SRT格式字幕内容
func BuildSRT(entries []models.SubtitleEntry, outputFormat string) string {
	var builder strings.Builder

	for i, entry := range entries {
		// 写入索引
		builder.WriteString(strconv.Itoa(entry.Index))
		builder.WriteString("\n")

		// 写入时间范围
		builder.WriteString(entry.TimeRange)
		builder.WriteString("\n")

		// 根据输出格式写入内容
		if outputFormat == "bilingual" {
			// 双语模式：内容已经包含原始文本和翻译文本
			builder.WriteString(entry.Content)
		} else {
			// 单语模式：内容只包含翻译文本
			// 由于内容已经更新为翻译文本，直接写入
			builder.WriteString(entry.Content)
		}

		// 在条目之间添加空行，最后一个条目后不添加
		if i < len(entries)-1 {
			builder.WriteString("\n\n")
		}
	}

	return builder.String()
}
