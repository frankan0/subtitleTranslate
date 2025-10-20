<template>
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
    <!-- 左侧：翻译选项 -->
    <div class="lg:col-span-1">
      <div class="sticky top-10 space-y-6">
        <TranslationOptions 
          v-model:target-language="targetLanguage"
          v-model:output-format="outputFormat"
          v-model:translation-position="translationPosition"
          v-model:provider="provider"
        />
        <div v-if="files.length > 0" class="mt-4">
          <button 
            @click="translateFiles" 
            class="w-full px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors"
            :disabled="isTranslating"
          >
            {{ isTranslating ? '翻译中...' : '开始翻译' }}
          </button>
        </div>
      </div>
    </div>
    
    <!-- 右侧：拖放框 + 文件列表 -->
    <div class="lg:col-span-2 space-y-6">
      <FileUploader @files-uploaded="handleFilesUploaded" />
      <FileList 
        v-if="files.length > 0" 
        :files="files" 
        :translations="translations" 
        @remove-file="handleFileRemove"
      />
      <div v-else class="flex items-center justify-center h-64 bg-muted/30 rounded-lg border-2 border-dashed">
        <p class="text-muted-foreground">请上传文件开始翻译</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import FileUploader from './FileUploader.vue';
import TranslationOptions from './TranslationOptions.vue';
import FileList from './FileList.vue';
import { translateSubtitles } from '../services/translationService';
import { getCurrentApiName, getApiSettings } from '../services/apiSettingsService';
import type { SubtitleFile, TranslationResult, OutputFormat, TranslationPosition, TranslationProvider } from '../types';

const files = ref<SubtitleFile[]>([]);
const translations = ref<TranslationResult[]>([]);
const targetLanguage = ref('zh');
const outputFormat = ref<OutputFormat>('translation_only');
const translationPosition = ref<TranslationPosition>('bottom');
const provider = ref<TranslationProvider>('volce');
const isTranslating = ref(false);

const handleFilesUploaded = (uploadedFiles: SubtitleFile[]) => {
  // 检查是否有重复文件
  const newFiles = uploadedFiles.filter(uploadedFile => {
    return !files.value.some(existingFile => 
      existingFile.name === uploadedFile.name && 
      existingFile.size === uploadedFile.size
    );
  });

  // 将新文件追加到现有文件列表
  files.value = [...files.value, ...newFiles];
};

const handleFileRemove = (fileId: string) => {
  files.value = files.value.filter(file => file.id !== fileId);
  translations.value = translations.value.filter(translation => translation.fileId !== fileId);
};

const translateFiles = async () => {
  if (files.value.length === 0) return;
  
  // 获取当前选择的翻译提供商的API配置
  const apiSettings = getApiSettings(provider.value);
  
  // 校验API配置是否存在
  if (!apiSettings.apiKey ) {
    alert(`请先在API设置中配置API密钥`);
    return;
  }
  
  isTranslating.value = true;
  try {
    // 只翻译没有翻译结果的文件
    const untranslatedFiles = files.value.filter(file => 
      !translations.value.some(t => t.fileId === file.id)
    );
    
    if (untranslatedFiles.length > 0) {
      const newTranslations = await translateSubtitles(
        untranslatedFiles, 
        targetLanguage.value,
        outputFormat.value,
        translationPosition.value,
        provider.value
      );
      translations.value = [...translations.value, ...newTranslations];
    }
  } catch (error) {
    console.error('Translation failed:', error);
    // 可以添加错误处理逻辑
  } finally {
    isTranslating.value = false;
  }
};
</script>