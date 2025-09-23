import type { SubtitleEntry } from '../types';

/**
 * 解析VTT格式的字幕文件内容
 * @param content VTT文件内容
 * @returns 解析后的字幕条目数组
 */
export function parseVtt(content: string): SubtitleEntry[] {
  // 规范化换行符
  const normalizedContent = content.replace(/\r\n/g, '\n').replace(/\r/g, '\n');
  
  // 移除WEBVTT头部和BOM
  let cleanContent = normalizedContent.trim();
  if (cleanContent.startsWith('\uFEFF')) {
    cleanContent = cleanContent.slice(1);
  }
  cleanContent = cleanContent.replace(/^WEBVTT.*\n?/, '').trim();
  
  // 按空行分割字幕块
  const blocks = cleanContent.split(/\n\n+/);
  const entries: SubtitleEntry[] = [];
  let entryId = 1;
  
  for (const block of blocks) {
    if (!block.trim()) continue;
    
    const lines = block.split('\n').filter(line => line.trim());
    if (lines.length < 2) continue;
    
    let timeLineIndex = 0;
    let currentId = entryId;
    
    // 检查第一行是否是数字ID
    const firstLine = lines[0].trim();
    const numericId = parseInt(firstLine, 10);
    if (!isNaN(numericId) && numericId > 0) {
      currentId = numericId;
      timeLineIndex = 1;
    }
    
    if (timeLineIndex >= lines.length) continue;
    
    // 解析时间行
    const timeLine = lines[timeLineIndex].trim();
    const timeMatch = timeLine.match(/([\d:.]+)\s+--?>\s+([\d:.]+)/);
    if (!timeMatch) continue;
    
    const [, startTime, endTime] = timeMatch;
    
    // 剩余行是字幕文本
    const text = lines.slice(timeLineIndex + 1).join('\n').trim();
    if (!text) continue;
    
    entries.push({
      id: currentId,
      startTime,
      endTime,
      text
    });
    
    entryId++;
  }
  
  return entries;
}

/**
 * 将字幕条目数组格式化为VTT格式的字符串
 * @param entries 字幕条目数组
 * @returns VTT格式的字符串
 */
export function formatVtt(entries: SubtitleEntry[]): string {
  let vttContent = 'WEBVTT\n\n';
  
  for (const entry of entries) {
    vttContent += `${entry.id}\n`;
    vttContent += `${entry.startTime} --> ${entry.endTime}\n`;
    vttContent += `${entry.text}\n\n`;
  }
  
  return vttContent.trim();
}

/**
 * 验证内容是否为有效的VTT格式
 * @param content 文件内容
 * @returns 是否为有效VTT格式
 */
export function isValidVtt(content: string): boolean {
  const normalizedContent = content.trim();
  
  // 检查是否包含WEBVTT头部
  if (!normalizedContent.startsWith('WEBVTT') && !normalizedContent.includes('WEBVTT')) {
    return false;
  }
  
  // 检查时间格式
  const timeRegex = /\d{2}:\d{2}:\d{2}\.\d{3}\s+--?>\s+\d{2}:\d{2}:\d{2}\.\d{3}/;
  return timeRegex.test(normalizedContent);
}