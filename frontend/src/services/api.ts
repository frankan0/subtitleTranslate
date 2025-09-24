import axios from 'axios';
import type { TranslationOptions, SubtitleFile, TranslationResult } from '../types';
import { getCurrentApiName, getApiSettings } from './apiSettingsService';

// 创建axios实例
const createApi = () => {
  const apiName = getCurrentApiName();
  const settings = getApiSettings(apiName);
  
  return axios.create({
    baseURL: settings.apiUrl || import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
    timeout: 120000, // 2分钟超时
    headers: {
      'Content-Type': 'application/json',
    },
  });
};

const getApi = () => createApi();

// 翻译字幕文件
export async function translateSubtitle(
  file: SubtitleFile,
  options: TranslationOptions
): Promise<{
  success: boolean;
  translatedFilename?: string;
  content?: string;
  error?: string;
}> {
  try {
    const api = getApi();
    
    // 获取当前API设置
    const apiName = getCurrentApiName();
    const apiSettings = getApiSettings(apiName);
    
    const response = await api.post('/subtitle/translate', {
      filename: file.name,
      content: file.content,
      targetLanguage: options.targetLanguage,
      sourceLanguage: options.sourceLanguage,
      provider: options.provider,
      outputFormat: options.outputFormat,
      translationPosition: options.translationPosition,
      // 传递API密钥和设置
      apiKey: apiSettings.apiKey,
      apiSecret: apiSettings.apiSecret,
      apiUrl: apiSettings.apiUrl,
    });
    if (response.data.success) {
      return {
        success: true,
        translatedFilename: response.data.data.translatedFilename,
        content: response.data.data.content,
      };
    } else {
      throw new Error(response.data.error || '翻译失败');
    }
  } catch (error) {
    console.error('翻译请求失败:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : '未知错误',
    };
  }
}

// 健康检查
export async function checkHealth(): Promise<boolean> {
  try {
    const api = getApi();
    const response = await api.get('/health');
    return response.data.status === 'ok';
  } catch (error) {
    console.error('健康检查失败:', error);
    return false;
  }
}

export default { getApi };