<template>
  <div class="w-full p-6 bg-card text-card-foreground rounded-lg border shadow-sm">
    <h3 class="text-lg font-medium mb-4">翻译选项</h3>
    
    <div class="space-y-4">
      <div class="grid gap-2">
        <label for="source-language" class="text-sm font-medium">源语言</label>
        <select
          id="source-language"
          v-model="localSourceLanguage"
          class="w-full px-3 py-2 border rounded-md bg-background"
        >
          <option value="auto">自动检测</option>
          <option v-for="lang in languages" :key="lang.code" :value="lang.code">
            {{ lang.name }}
          </option>
        </select>
      </div>
      
      <div class="grid gap-2">
        <label for="target-language" class="text-sm font-medium">目标语言</label>
        <select
          id="target-language"
          v-model="localTargetLanguage"
          class="w-full px-3 py-2 border rounded-md bg-background"
        >
          <option v-for="lang in languages" :key="lang.code" :value="lang.code">
            {{ lang.name }}
          </option>
        </select>
      </div>
      
      <div class="grid gap-2">
        <label for="translation-provider" class="text-sm font-medium">翻译提供商</label>
        <select
          id="translation-provider"
          v-model="provider"
          class="w-full px-3 py-2 border rounded-md bg-background"
        >
          <option value="volce">火山引擎</option>
          <option value="tencent">腾讯云</option>
        </select>
      </div>

      <div class="grid gap-2">
        <label class="text-sm font-medium">输出格式</label>
        <div class="space-y-2">
          <label class="flex items-center space-x-2">
            <input
              type="radio"
              v-model="outputFormat"
              value="translation_only"
              class="form-radio"
            />
            <span>仅译文</span>
          </label>
          <label class="flex items-center space-x-2">
            <input
              type="radio"
              v-model="outputFormat"
              value="original_and_translation"
              class="form-radio"
            />
            <span>原文+译文</span>
          </label>
        </div>
      </div>

      <div v-if="outputFormat === 'original_and_translation'" class="grid gap-2">
        <label class="text-sm font-medium">译文位置</label>
        <div class="space-y-2">
          <label class="flex items-center space-x-2">
            <input
              type="radio"
              v-model="translationPosition"
              value="bottom"
              class="form-radio"
            />
            <span>译文在原文下方</span>
          </label>
          <label class="flex items-center space-x-2">
            <input
              type="radio"
              v-model="translationPosition"
              value="top"
              class="form-radio"
            />
            <span>译文在原文上方</span>
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { setCurrentApiName } from '../services/apiSettingsService';
import type { TranslationProvider, OutputFormat, TranslationPosition } from '../types';

const props = defineProps<{
  targetLanguage: string;
  sourceLanguage?: string;
  provider?: TranslationProvider;
}>();

const emit = defineEmits<{
  (e: 'update:target-language', value: string): void;
  (e: 'update:source-language', value: string): void;
  (e: 'update:output-format', value: OutputFormat): void;
  (e: 'update:translation-position', value: TranslationPosition): void;
  (e: 'update:provider', value: TranslationProvider): void;
}>();

const localTargetLanguage = computed({
  get: () => props.targetLanguage || 'zh',
  set: (value) => emit('update:target-language', value)
});

const localSourceLanguage = computed({
  get: () => props.sourceLanguage || 'auto',
  set: (value) => emit('update:source-language', value)
});

const provider = computed({
  get: () => props.provider || 'volce',
  set: (value) => {
    emit('update:provider', value);
    setCurrentApiName(value);
  }
});
const outputFormat = ref<OutputFormat>('translation_only');
const translationPosition = ref<TranslationPosition>('bottom');

watch(outputFormat, (newValue) => {
  emit('update:output-format', newValue);
});

watch(translationPosition, (newValue) => {
  emit('update:translation-position', newValue);
});

// 支持的语言列表
const languages = [
  { code: 'zh', name: '中文' },
  { code: 'en', name: '英语' },
  { code: 'ja', name: '日语' },
  { code: 'ko', name: '韩语' },
  { code: 'fr', name: '法语' },
  { code: 'de', name: '德语' },
  { code: 'es', name: '西班牙语' },
  { code: 'ru', name: '俄语' },
  { code: 'it', name: '意大利语' },
  { code: 'pt', name: '葡萄牙语' }
];
</script>