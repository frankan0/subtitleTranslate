import type { SubtitleEntry } from '../types';

/**
 * 解析ASS格式的字幕内容
 * @param content ASS文件内容
 * @returns 解析后的字幕条目数组
 */
export function parseAss(content: string): SubtitleEntry[] {
  const entries: SubtitleEntry[] = [];
  const lines = content.split('\n');
  
  let inEventsSection = false;
  let lineIndex = 1;
  
  for (const line of lines) {
    const trimmedLine = line.trim();
    
    if (trimmedLine.toLowerCase().startsWith('[events]')) {
      inEventsSection = true;
      continue;
    }
    
    if (!inEventsSection) {
      continue;
    }
    
    // 跳过格式行和注释
    if (trimmedLine.toLowerCase().startsWith('format:') || trimmedLine.startsWith(';') || trimmedLine === '') {
      continue;
    }
    
    // 解析Dialogue行
    if (trimmedLine.toLowerCase().startsWith('dialogue:')) {
      const entry = parseAssDialogue(trimmedLine, lineIndex);
      if (entry) {
        entries.push(entry);
        lineIndex++;
      }
    }
  }
  
  return entries;
}

/**
 * 解析ASS格式的Dialogue行
 * @param line Dialogue行内容
 * @param index 字幕索引
 * @returns 解析后的字幕条目
 */
function parseAssDialogue(line: string, index: number): SubtitleEntry | null {
  // 移除"Dialogue:"前缀（不区分大小写）
  const content = line.replace(/^dialogue:\s*/i, '');
  
  // 按逗号分割，但需要考虑引号
  const parts = splitAssLine(content);
  if (parts.length < 10) {
    return null;
  }
  
  const startTime = parts[1];
  const endTime = parts[2];
  
  // 提取文本内容（从第10个字段开始）
  const textParts = parts.slice(9);
  let text = textParts.join(',');
  
  // 清理ASS标签
  text = cleanAssTags(text);
  
  // 转换时间格式
  const startTimeStr = convertAssTimeToStandard(startTime);
  const endTimeStr = convertAssTimeToStandard(endTime);
  
  return {
    index,
    timeRange: `${startTimeStr} --> ${endTimeStr}`,
    content: text.trim()
  };
}

/**
 * 智能分割ASS行，考虑引号内的逗号
 * @param line 行内容
 * @returns 分割后的字段数组
 */
function splitAssLine(line: string): string[] {
  const parts: string[] = [];
  let current = '';
  let inQuotes = false;
  
  for (let i = 0; i < line.length; i++) {
    const char = line[i];
    
    if (char === '"') {
      inQuotes = !inQuotes;
    } else if (char === ',' && !inQuotes) {
      parts.push(current.trim());
      current = '';
      continue;
    }
    
    current += char;
  }
  
  if (current) {
    parts.push(current.trim());
  }
  
  return parts;
}

/**
 * 清理ASS文本中的标签
 * @param text 原始文本
 * @returns 清理后的文本
 */
function cleanAssTags(text: string): string {
  // 移除花括号标签
  text = text.replace(/\{[^}]*\}/g, '');
  
  // 移除方括号标签
  text = text.replace(/\[[^\]]*\]/g, '');
  
  // 处理换行符
  text = text.replace(/\\N/g, '\n');
  text = text.replace(/\\n/g, '\n');
  
  // 处理硬空格
  text = text.replace(/\\h/g, ' ');
  
  return text;
}

/**
 * 将ASS时间格式转换为标准时间格式
 * @param assTime ASS时间格式 (h:mm:ss.cc)
 * @returns 标准时间格式 (HH:MM:SS,mmm)
 */
function convertAssTimeToStandard(assTime: string): string {
  const parts = assTime.split(':');
  if (parts.length !== 3) {
    return '00:00:00,000';
  }
  
  const [hours, minutes, secondsPart] = parts;
  const [seconds, centiseconds = '00'] = secondsPart.split('.');
  
  // 确保两位数的小时和分钟
  const h = hours.padStart(2, '0');
  const m = minutes.padStart(2, '0');
  const s = seconds.padStart(2, '0');
  
  // 将百分之一秒转换为毫秒
  const ms = (centiseconds + '00').slice(0, 3).padStart(3, '0');
  
  return `${h}:${m}:${s},${ms}`;
}

/**
 * 构建ASS格式的字幕内容
 * @param entries 字幕条目数组
 * @returns ASS格式的字幕字符串
 */
export function formatAss(entries: SubtitleEntry[]): string {
  const lines: string[] = [];
  
  // ASS文件头
  lines.push('[Script Info]');
  lines.push('Title: Translated Subtitle');
  lines.push('ScriptType: v4.00+');
  lines.push('WrapStyle: 0');
  lines.push('ScaledBorderAndShadow: yes');
  lines.push('YCbCr Matrix: None');
  lines.push('');
  
  // 样式部分
  lines.push('[V4+ Styles]');
  lines.push('Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding');
  lines.push('Style: Default,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,0,0,0,0,100,100,0,0,1,2,2,2,10,10,10,1');
  lines.push('');
  
  // 事件部分
  lines.push('[Events]');
  lines.push('Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text');
  
  for (const entry of entries) {
    const [startTime, endTime] = parseTimeRange(entry.timeRange);
    const content = entry.content.replace(/\n/g, '\\N');
    
    lines.push(`Dialogue: 0,${startTime},${endTime},Default,,0,0,0,,${content}`);
  }
  
  return lines.join('\n');
}

/**
 * 解析时间范围字符串
 * @param timeRange 时间范围字符串 (HH:MM:SS,mmm --> HH:MM:SS,mmm)
 * @returns [开始时间, 结束时间] 的ASS格式
 */
function parseTimeRange(timeRange: string): [string, string] {
  const [start, end] = timeRange.split(' --> ').map(t => t.trim());
  return [convertStandardTimeToAss(start), convertStandardTimeToAss(end)];
}

/**
 * 将标准时间格式转换为ASS时间格式
 * @param standardTime 标准时间格式 (HH:MM:SS,mmm)
 * @returns ASS时间格式 (h:mm:ss.cc)
 */
function convertStandardTimeToAss(standardTime: string): string {
  const [time, milliseconds] = standardTime.split(',');
  const [hours, minutes, seconds] = time.split(':');
  
  // 将毫秒转换为百分之一秒
  const ms = parseInt(milliseconds || '000', 10);
  const centiseconds = Math.floor(ms / 10).toString().padStart(2, '0');
  
  // 移除小时的前导零
  const h = parseInt(hours, 10).toString();
  
  return `${h}:${minutes}:${seconds}.${centiseconds}`;
}