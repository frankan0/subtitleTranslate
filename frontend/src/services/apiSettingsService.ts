import axios from 'axios';

// 定义API设置接口
export interface ApiSettings {
  apiUrl: string;
  chunkSize: number;
  delay: number;
  rate: number;
}

// 获取API设置
export const getApiSettings = (apiName: string): ApiSettings => {
  const defaultSettings: ApiSettings = {
    apiUrl: '',
    chunkSize: 1000,
    delay: 200,
    rate: 10,
  };

  const storageKey = `api-settings-${apiName}`;
  const savedSettings = localStorage.getItem(storageKey);
  
  if (savedSettings) {
    try {
      return JSON.parse(savedSettings);
    } catch (error) {
      console.error('Failed to parse API settings:', error);
      return defaultSettings;
    }
  }
  
  return defaultSettings;
};

// 保存API设置
export const saveApiSettings = (apiName: string, settings: ApiSettings): void => {
  const storageKey = `api-settings-${apiName}`;
  localStorage.setItem(storageKey, JSON.stringify(settings));
};

// 清除翻译缓存
export const clearTranslationCache = (): void => {
  // 查找所有与翻译相关的缓存项并删除
  const keysToRemove: string[] = [];
  
  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i);
    if (key && key.startsWith('translation-cache-')) {
      keysToRemove.push(key);
    }
  }
  
  keysToRemove.forEach(key => localStorage.removeItem(key));
};

// 获取当前选择的API名称
export const getCurrentApiName = (): string => {
  const savedApiName = localStorage.getItem('current-api-name');
  return savedApiName || 'deeplx'; // 默认使用DeepLX
};

// 设置当前选择的API名称
export const setCurrentApiName = (apiName: string): void => {
  localStorage.setItem('current-api-name', apiName);
};