package subtitle

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/frank0/subtitleTranslate/internal/models"
	"github.com/frank0/subtitleTranslate/internal/utils"
)

// Parser 字幕解析器接口
type Parser interface {
	Parse(content string) ([]models.SubtitleEntry, error)
	SupportedExtensions() []string
}

// ParserFactory 解析器工厂
type ParserFactory struct {
	parsers map[string]Parser
}

// NewParserFactory 创建新的解析器工厂
func NewParserFactory() *ParserFactory {
	factory := &ParserFactory{
		parsers: make(map[string]Parser),
	}
	
	// 注册支持的解析器
	factory.Register(".srt", &SRTParser{})
	factory.Register(".vtt", &VTTParser{})
	factory.Register(".ass", &ASSParser{})
	
	return factory
}

// Register 注册新的解析器
func (f *ParserFactory) Register(extension string, parser Parser) {
	f.parsers[strings.ToLower(extension)] = parser
}

// GetParser 根据文件扩展名获取解析器
func (f *ParserFactory) GetParser(filename string) (Parser, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	parser, exists := f.parsers[ext]
	if !exists {
		return nil, fmt.Errorf("不支持的文件格式: %s", ext)
	}
	return parser, nil
}

// GetSupportedExtensions 获取所有支持的扩展名
func (f *ParserFactory) GetSupportedExtensions() []string {
	var extensions []string
	for _, parser := range f.parsers {
		extensions = append(extensions, parser.SupportedExtensions()...)
	}
	return extensions
}

// SRTParser SRT格式解析器
type SRTParser struct{}

func (p *SRTParser) Parse(content string) ([]models.SubtitleEntry, error) {
	return utils.ParseSRT(content)
}

func (p *SRTParser) SupportedExtensions() []string {
	return []string{".srt"}
}

// VTTParser VTT格式解析器
type VTTParser struct{}

func (p *VTTParser) Parse(content string) ([]models.SubtitleEntry, error) {
	return utils.ParseVTT(content)
}

func (p *VTTParser) SupportedExtensions() []string {
	return []string{".vtt"}
}

// ASSParser ASS格式解析器
type ASSParser struct{}

func (p *ASSParser) Parse(content string) ([]models.SubtitleEntry, error) {
	return utils.ParseASS(content)
}

func (p *ASSParser) SupportedExtensions() []string {
	return []string{".ass", ".ssa"}
}