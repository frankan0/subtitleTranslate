// 字幕文件类型
export interface SubtitleFile {
  id: string;
  name: string;
  content: string;
  size: number;
  type: string;
  lastModified: number;
  charCount: number;  // 文件字符数
}

// 字幕条目类型
export interface SubtitleEntry {
  id: number;
  startTime: string;
  endTime: string;
  text: string;
}

// 翻译结果类型
export interface TranslationResult {
  fileId: string;
  fileName: string;
  originalName: string;
  translatedContent: string;
  status: 'success' | 'error';
  error?: string;
}

// 翻译API提供商类型
export type TranslationProvider = 'volce' | 'google' | 'other';

// 翻译输出格式类型
export type OutputFormat = 'translation_only' | 'original_and_translation';
export type TranslationPosition = 'top' | 'bottom';

// 翻译选项类型
export interface TranslationOptions {
  targetLanguage: string;
  sourceLanguage?: string;            // 源语言，支持腾讯云等需要明确源语言的API
  provider?: TranslationProvider;
  outputFormat?: OutputFormat;        // 输出格式：仅译文或原文+译文
  translationPosition?: TranslationPosition;  // 译文位置：上方或下方
  bilingual?: boolean;                // 是否生成双语字幕（兼容旧版本）
  keepOriginal?: boolean;             // 是否保留原文（兼容旧版本）
}

// 翻译API响应类型
export interface TranslationApiResponse {
  success: boolean;
  data?: {
    translatedContent: string;
  };
  error?: string;
}