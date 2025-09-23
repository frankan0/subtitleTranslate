import { parseSrt } from './srtParser';
import { parseVtt, formatVtt } from './vttParser';
import { parseAss, formatAss } from './assParser';
import type { SubtitleEntry, TranslationResult } from '../types';

/**
 * 根据文件扩展名解析字幕内容
 * @param content 字幕文件内容
 * @param fileName 文件名
 * @returns 解析后的字幕条目数组
 */
export function parseSubtitle(content: string, fileName: string): SubtitleEntry[] {
  const extension = fileName.split('.').pop()?.toLowerCase();
  
  switch (extension) {
    case 'srt':
      return parseSrt(content);
    case 'vtt':
      return parseVtt(content);
    case 'ass':
    case 'ssa':
      return parseAss(content);
    default:
      // 尝试根据内容格式判断
      if (content.trim().startsWith('WEBVTT')) {
        return parseVtt(content);
      } else if (content.includes('[Events]') && content.includes('Dialogue:')) {
        return parseAss(content);
      } else {
        return parseSrt(content);
      }
  }
}

/**
 * 下载字幕文件，根据原文件名自动选择格式
 * @param content 字幕内容
 * @param originalFileName 原始文件名
 * @param format 指定格式，如果未指定则根据原文件名决定
 */
export function downloadSubtitleFile(content: string, originalFileName: string, format?: 'srt' | 'vtt' | 'ass') {
  const extension = format || originalFileName.split('.').pop()?.toLowerCase() || 'srt';
  const baseName = originalFileName.replace(/\.[^/.]+$/, '');
  const fileName = `${baseName}_translated.${extension}`;
  
  let mimeType = 'text/plain';
  if (extension === 'vtt') mimeType = 'text/vtt';
  else if (extension === 'ass' || extension === 'ssa') mimeType = 'text/plain';
  
  const blob = new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);
  
  const a = document.createElement('a');
  a.href = url;
  a.download = fileName;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

/**
 * 获取字幕内容的格式化字符串
 * @param entries 字幕条目数组
 * @param format 输出格式
 * @returns 格式化的字幕字符串
 */
export function formatSubtitle(entries: SubtitleEntry[], format: 'srt' | 'vtt' | 'ass'): string {
  switch (format) {
    case 'vtt':
      return formatVtt(entries);
    case 'ass':
      return formatAss(entries);
    default:
      // 使用srtParser中的formatSrt函数
      const { formatSrt } = require('./srtParser');
      return formatSrt(entries);
  }
}

/**
 * 获取可下载的翻译结果
 * @param translations 翻译结果数组
 * @returns 可下载的翻译结果
 */
export function getDownloadableTranslations(translations: Array<{
  status: string;
  translatedContent: string;
  fileName: string;
}>): Array<{
  translatedContent: string;
  fileName: string;
}> {
  return translations.filter(t => t.status === 'success');
}

/**
 * 检查是否可以批量下载
 * @param translations 翻译结果数组
 * @returns 是否可以批量下载
 */
export function canBatchDownload(translations: Array<{
  status: string;
}>): boolean {
  return translations.some(t => t.status === 'success');
}

/**
 * 创建ZIP文件并下载多个翻译结果
 * @param translations 翻译结果数组
 */
export async function downloadTranslationsAsZip(translations: Array<{
  translatedContent: string;
  fileName: string;
}>) {
  try {
    const JSZip = (await import('jszip')).default;
    const zip = new JSZip();
    
    // 将所有翻译文件添加到ZIP中
    translations.forEach((translation) => {
      const extension = translation.fileName.split('.').pop()?.toLowerCase() || 'srt';
      const baseName = translation.fileName.replace(/\.[^/.]+$/, '');
      const fileName = `${baseName}_translated.${extension}`;
      zip.file(fileName, translation.translatedContent);
    });
    
    // 生成ZIP文件并下载
    const zipBlob = await zip.generateAsync({ type: 'blob' });
    const url = URL.createObjectURL(zipBlob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `translated_subtitles_${new Date().getTime()}.zip`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  } catch (error) {
    console.error('Failed to create ZIP file:', error);
    // 如果ZIP创建失败，回退到逐个下载
    translations.forEach(translation => {
      downloadSubtitleFile(translation.translatedContent, translation.fileName);
    });
  }
}