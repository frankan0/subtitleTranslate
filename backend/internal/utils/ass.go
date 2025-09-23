package utils

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/frank0/subtitleTranslate/internal/models"
)

// ParseASS 解析ASS格式的字幕文件内容
func ParseASS(content string) ([]models.SubtitleEntry, error) {
	var entries []models.SubtitleEntry
	scanner := bufio.NewScanner(strings.NewReader(content))

	// 跳过头部信息，直接找到[Events]部分
	var inEventsSection bool
	var lineIndex int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(strings.ToLower(line), "[events]") {
			inEventsSection = true
			continue
		}

		if !inEventsSection {
			continue
		}

		// 跳过格式行和注释
		if strings.HasPrefix(strings.ToLower(line), "format:") || strings.HasPrefix(line, ";") || line == "" {
			continue
		}

		// 解析Dialogue行
		if strings.HasPrefix(strings.ToLower(line), "dialogue:") {
			entry, err := parseASSDialogue(line, lineIndex+1)
			if err != nil {
				continue // 跳过解析失败的行
			}
			entries = append(entries, *entry)
			lineIndex++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("解析ASS文件失败: %w", err)
	}

	return entries, nil
}

// parseASSDialogue 解析ASS格式的Dialogue行
func parseASSDialogue(line string, index int) (*models.SubtitleEntry, error) {
	// 移除"Dialogue:"前缀（不区分大小写）
	line = strings.TrimPrefix(strings.ToLower(line), "dialogue:")
	line = strings.TrimSpace(line)

	// 按逗号分割，但需要考虑引号内的逗号
	parts := splitASSLine(line)
	if len(parts) < 10 {
		return nil, fmt.Errorf("ASS Dialogue格式不正确")
	}

	// 提取时间和文本内容
	startTime := parts[1]
	endTime := parts[2]

	// 将ASS时间格式转换为标准时间格式
	startTimeStr := convertASSTime(startTime)
	endTimeStr := convertASSTime(endTime)

	// 提取文本内容（第10个字段开始）
	var textContent strings.Builder
	for i := 9; i < len(parts); i++ {
		if i > 9 {
			textContent.WriteString(",")
		}
		textContent.WriteString(parts[i])
	}

	// 清理文本中的ASS标签
	content := cleanASSTags(textContent.String())

	return &models.SubtitleEntry{
		Index:     index,
		TimeRange: fmt.Sprintf("%s --> %s", startTimeStr, endTimeStr),
		Content:   content,
	}, nil
}

// splitASSLine 智能分割ASS行，考虑引号内的逗号
func splitASSLine(line string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false

	for _, char := range line {
		switch char {
		case '"':
			inQuotes = !inQuotes
			current.WriteRune(char)
		case ',':
			if !inQuotes {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
				continue
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	return parts
}

// convertASSTime 将ASS时间格式转换为标准时间格式
func convertASSTime(assTime string) string {
	// ASS格式: 0:00:00.00 (小时:分钟:秒.百分之一秒)
	// 转换为: 00:00:00,000 (小时:分钟:秒,毫秒)
	parts := strings.Split(assTime, ":")
	if len(parts) != 3 {
		return "00:00:00,000"
	}

	hours := parts[0]
	minutes := parts[1]
	secondsParts := strings.Split(parts[2], ".")

	seconds := secondsParts[0]
	if len(seconds) == 1 {
		seconds = "0" + seconds
	}

	var milliseconds string
	if len(secondsParts) > 1 {
		// 将百分之一秒转换为毫秒
		centiseconds := secondsParts[1]
		if len(centiseconds) == 1 {
			centiseconds += "0"
		}
		if len(centiseconds) >= 2 {
			ms, _ := strconv.Atoi(centiseconds[:2])
			milliseconds = fmt.Sprintf("%03d", ms*10)
		} else {
			milliseconds = "000"
		}
	} else {
		milliseconds = "000"
	}

	// 确保小时和分钟都是两位数
	if len(hours) == 1 {
		hours = "0" + hours
	}
	if len(minutes) == 1 {
		minutes = "0" + minutes
	}

	return fmt.Sprintf("%s:%s:%s,%s", hours, minutes, seconds, milliseconds)
}

// cleanASSTags 清理ASS文本中的标签
func cleanASSTags(text string) string {
	// 移除基本的ASS标签
	re := regexp.MustCompile(`\[.*?\]`)
	text = re.ReplaceAllString(text, "")

	// 移除花括号标签
	re = regexp.MustCompile(`\{.*?\}`)
	text = re.ReplaceAllString(text, "")

	// 移除反斜杠转义
	text = strings.ReplaceAll(text, "\\N", "\n")
	text = strings.ReplaceAll(text, "\\n", "\n")
	text = strings.ReplaceAll(text, "\\h", " ")

	// 清理多余的空格
	text = strings.TrimSpace(text)

	return text
}

// BuildASS 构建ASS格式的字幕内容
func BuildASS(entries []models.SubtitleEntry, outputFormat string) string {
	var builder strings.Builder

	// ASS文件头
	builder.WriteString("[Script Info]\n")
	builder.WriteString("Title: Translated Subtitle\n")
	builder.WriteString("ScriptType: v4.00+\n")
	builder.WriteString("WrapStyle: 0\n")
	builder.WriteString("ScaledBorderAndShadow: yes\n")
	builder.WriteString("YCbCr Matrix: None\n\n")

	// 样式部分
	builder.WriteString("[V4+ Styles]\n")
	builder.WriteString("Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding\n")
	builder.WriteString("Style: Default,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,0,0,0,0,100,100,0,0,1,2,2,2,10,10,10,1\n\n")

	// 事件部分
	builder.WriteString("[Events]\n")
	builder.WriteString("Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n")

	for _, entry := range entries {
		// 解析时间范围
		timeRange := strings.Split(entry.TimeRange, " --> ")
		if len(timeRange) != 2 {
			continue
		}

		startTime := convertToASSTime(strings.TrimSpace(timeRange[0]))
		endTime := convertToASSTime(strings.TrimSpace(timeRange[1]))

		// 处理内容格式
		content := entry.Content
		content = strings.ReplaceAll(content, "\n", "\\N")

		if outputFormat == "bilingual" {
			// 双语模式：内容已经包含原始文本和翻译文本
			builder.WriteString(fmt.Sprintf("Dialogue: 0,%s,%s,Default,,0,0,0,,%s\n", startTime, endTime, content))
		} else {
			// 单语模式：内容只包含翻译文本
			builder.WriteString(fmt.Sprintf("Dialogue: 0,%s,%s,Default,,0,0,0,,%s\n", startTime, endTime, content))
		}
	}

	return builder.String()
}

// convertToASSTime 将标准时间格式转换为ASS时间格式
func convertToASSTime(standardTime string) string {
	// 标准格式: 00:00:00,000 (小时:分钟:秒,毫秒)
	// 转换为: 0:00:00.00 (小时:分钟:秒.百分之一秒)
	parts := strings.Split(standardTime, ":")
	if len(parts) != 3 {
		return "0:00:00.00"
	}

	hours := strings.TrimLeft(parts[0], "0")
	if hours == "" {
		hours = "0"
	}

	minutes := parts[1]
	secondsParts := strings.Split(parts[2], ",")

	seconds := secondsParts[0]
	var centiseconds string
	if len(secondsParts) > 1 {
		// 将毫秒转换为百分之一秒
		ms, _ := strconv.Atoi(secondsParts[1])
		centiseconds = fmt.Sprintf("%02d", ms/10)
	} else {
		centiseconds = "00"
	}

	return fmt.Sprintf("%s:%s:%s.%s", hours, minutes, seconds, centiseconds)
}
