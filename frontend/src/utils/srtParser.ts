/**
 * SRT格式字幕解析工具
 */

import type { SubtitleEntry } from '@/types';

/**
 * 解析SRT格式的字幕内容
 * @param content SRT格式的字幕内容
 * @returns 字幕条目数组
 */
export function parseSrt(content: string): SubtitleEntry[] {
  const entries: SubtitleEntry[] = [];
  const blocks = content.trim().split('\n\n');
  
  for (const block of blocks) {
    const lines = block.split('\n');
    if (lines.length >= 3) {
      const id = parseInt(lines[0], 10);
      const timeLine = lines[1];
      const text = lines.slice(2).join('\n');
      
      const [startTime, endTime] = timeLine.split(' --> ');
      
      entries.push({
        id,
        startTime: startTime.trim(),
        endTime: endTime.trim(),
        text: text.trim()
      });
    }
  }
  
  return entries;
}

/**
 * 将字幕条目数组格式化为SRT格式的字符串
 * @param entries 字幕条目数组
 * @returns SRT格式的字符串
 */
export function formatSrt(entries: SubtitleEntry[]): string {
  let srtContent = '';
  
  for (const entry of entries) {
    srtContent += `${entry.id}\n`;
    srtContent += `${entry.startTime} --> ${entry.endTime}\n`;
    srtContent += `${entry.text}\n\n`;
  }
  
  return srtContent.trim();
}

/**
 * 下载字幕文件
 * @param content 文件内容
 * @param fileName 文件名
 */
export function downloadSrtFile(content: string, fileName: string): void {
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  
  const link = document.createElement('a');
  link.href = url;
  link.download = fileName;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

/**
 * 批量下载字幕文件为ZIP
 * @param translations 翻译结果数组
 */
export function downloadTranslationsAsZip(translations: any[]): void {
  // 实现ZIP下载逻辑
  console.log('批量下载功能待实现');
}