<template>
  <div class="w-full bg-card text-card-foreground rounded-lg border shadow-sm">
    <div class="flex items-center justify-between px-4 py-3 border-b">
      <h3 class="text-base font-medium">文件列表</h3>
      <button
        v-if="hasSuccessfulTranslations"
        @click="downloadAll"
        class="inline-flex items-center px-3 py-1.5 text-sm bg-primary text-primary-foreground rounded hover:bg-primary/90 transition-colors"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-1.5">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
          <polyline points="7 10 12 15 17 10"></polyline>
          <line x1="12" y1="15" x2="12" y2="3"></line>
        </svg>
        批量下载 ({{ successfulTranslationCount }})
      </button>
    </div>
    
    <div class="divide-y">
      <div v-for="file in files" :key="file.id" class="px-4 py-2.5 hover:bg-muted/50 transition-colors">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-muted-foreground flex-shrink-0">
              <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
            </svg>
            <div class="min-w-0">
              <p class="text-sm font-medium truncate">{{ file.name }}</p>
              <p class="text-xs text-muted-foreground flex items-center gap-1">
                <span>{{ formatFileSize(file.size) }}</span>
                <span class="text-muted-foreground/50">•</span>
                <span>{{ formatCharCount(file.charCount) }} 字符</span>
              </p>
            </div>
          </div>
          
          <div class="flex items-center gap-1.5 flex-shrink-0 ml-2">
            <button
              v-if="!getTranslationStatus(file.id)"
              @click="removeFile(file.id)"
              class="p-1 text-destructive hover:bg-destructive/10 rounded transition-colors"
              title="删除文件"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M3 6h18"></path>
                <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"></path>
                <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"></path>
              </svg>
            </button>
            <span v-if="getTranslationStatus(file.id) === 'success'" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
              成功
            </span>
            <span v-else-if="getTranslationStatus(file.id) === 'error'" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
              失败
            </span>
          </div>
        </div>
        
        <!-- 翻译成功后显示下载按钮 -->
        <div v-if="getTranslationStatus(file.id) === 'success'" class="mt-2">
          <button 
            @click="downloadTranslation(file.id)" 
            class="inline-flex items-center px-2.5 py-1 text-xs bg-primary/10 text-primary rounded hover:bg-primary/20 transition-colors"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-1">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
              <polyline points="7 10 12 15 17 10"></polyline>
              <line x1="12" y1="15" x2="12" y2="3"></line>
            </svg>
            下载翻译文件
          </button>
        </div>
      </div>
    </div>
    
    <div v-if="files.length === 0" class="p-8 text-center text-sm text-muted-foreground">
      暂无文件
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { SubtitleFile, TranslationResult } from '../types';
import { downloadSubtitleFile, downloadTranslationsAsZip, getDownloadableTranslations } from '../utils/subtitleUtils';

const props = defineProps<{
  files: SubtitleFile[];
  translations: TranslationResult[];
}>();

const emit = defineEmits<{
  (e: 'remove-file', fileId: string): void
}>();

// 检查是否有成功翻译的文件
const hasSuccessfulTranslations = computed(() => {
  return props.translations.some(t => t.status === 'success');
});

// 成功翻译的文件数量
const successfulTranslationCount = computed(() => {
  return props.translations.filter(t => t.status === 'success').length;
});

// 批量下载所有成功翻译的文件
const downloadAll = async () => {
  const downloadableTranslations = getDownloadableTranslations(props.translations);
  
  if (downloadableTranslations.length === 0) {
    return;
  }
  
  if (downloadableTranslations.length === 1) {
    // 只有一个文件时直接下载
    const translation = downloadableTranslations[0];
    downloadSubtitleFile(translation.translatedContent, translation.fileName);
  } else {
    // 多个文件时打包成ZIP下载
    await downloadTranslationsAsZip(downloadableTranslations);
  }
};

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

// 格式化字符数
const formatCharCount = (count: number): string => {
  if (count < 1000) return count.toString();
  if (count < 10000) return (count / 1000).toFixed(1) + 'k';
  return (count / 10000).toFixed(1) + 'w';
};

// 获取文件的翻译状态
const getTranslationStatus = (fileId: string): 'success' | 'error' | null => {
  const translation = props.translations.find(t => t.fileId === fileId);
  return translation ? translation.status : null;
};

// 获取翻译错误信息
const getTranslationError = (fileId: string): string => {
  const translation = props.translations.find(t => t.fileId === fileId);
  return translation?.error || '未知错误';
};

// 下载翻译结果
const downloadTranslation = (fileId: string) => {
  const translation = props.translations.find(t => t.fileId === fileId);
  if (translation && translation.status === 'success') {
    downloadSubtitleFile(translation.translatedContent, translation.fileName);
  }
};

// 删除文件
const removeFile = (fileId: string) => {
  emit('remove-file', fileId);
};
</script>