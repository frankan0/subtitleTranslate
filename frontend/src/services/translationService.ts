import { translateSubtitle } from './api';
import { parseSubtitle } from '../utils/subtitleUtils';
import type { SubtitleFile, TranslationResult, TranslationProvider, TranslationOptions, OutputFormat, TranslationPosition } from '../types';

// 翻译字幕文件
export async function translateSubtitles(
  files: SubtitleFile[],
  targetLanguage: string,
  outputFormat: OutputFormat = 'translation_only',
  translationPosition: TranslationPosition = 'bottom',
  provider: TranslationProvider = 'volce'
): Promise<TranslationResult[]> {
  const results: TranslationResult[] = [];
  const options: TranslationOptions = {
    targetLanguage,
    sourceLanguage: 'auto',
    provider,
    outputFormat,
    translationPosition
  };

  for (const file of files) {
    try {
      // 解析字幕文件
      const entries = parseSubtitle(file.content, file.name);
      
      // 提取需要翻译的文本
      const textsToTranslate = entries.map(entry => entry.text);
      
      // 调用API服务进行翻译
      const translationResult = await translateSubtitle(file, options);
      console.log(translationResult)
      if (translationResult.success && translationResult.content) {
        let finalContent = translationResult.content||""
        results.push({
          fileId: file.id,
          fileName: translationResult.translatedFilename || file.name,
          originalName: file.name,
          translatedContent: finalContent,
          status: 'success'
        });
      } else {
        throw new Error(translationResult.error || '翻译失败');
      }
    } catch (error) {
      results.push({
        fileId: file.id,
        fileName: '',
        originalName: file.name,
        translatedContent: '',
        status: 'error',
        error: error instanceof Error ? error.message : '翻译失败'
      });
    }
  }
  console.log(results)
  return results;
}

// 生成翻译后的文件名
function generateTranslatedFilename(originalName: string, targetLanguage: string): string {
  const nameParts = originalName.split('.');
  const extension = nameParts.pop();
  const baseName = nameParts.join('.');
  
  const languageMap: Record<string, string> = {
    'zh': '中文',
    'en': '英文',
    'ja': '日文',
    'ko': '韩文',
    'fr': '法文',
    'de': '德文',
    'es': '西班牙文',
    'ru': '俄文',
    'it': '意大利文',
    'pt': '葡萄牙文'
  };
  
  const languageName = languageMap[targetLanguage] || targetLanguage;
  return `${baseName}_${languageName}.${extension}`;
}