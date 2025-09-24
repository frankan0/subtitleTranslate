<template>
  <div class="space-y-4">
    <div>
      <label for="api-key" class="block text-sm font-medium text-gray-700">API 密钥</label>
      <input type="password" id="api-key" v-model="settings.apiKey" placeholder="输入API密钥" class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
      <p class="mt-2 text-sm text-gray-500">API密钥将用于访问翻译服务</p>
    </div>
    <div>
      <label for="api-secret" class="block text-sm font-medium text-gray-700">API Secret</label>
      <input type="password" id="api-secret" v-model="settings.apiSecret" placeholder="输入API Secret" class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
      <p class="mt-2 text-sm text-gray-500">部分API服务需要Secret密钥</p>
    </div>
    <div>
      <label for="chunk-size" class="block text-sm font-medium text-gray-700">分块翻译大小</label>
      <input type="number" id="chunk-size" v-model="settings.chunkSize" class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
      <p class="mt-2 text-sm text-gray-500">针对字幕、Markdown 等多行文本, 将自动合并翻译。</p>
    </div>
    <div>
      <label for="delay" class="block text-sm font-medium text-gray-700">延迟时间 (ms)</label>
      <input type="number" id="delay" v-model="settings.delay" class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
    </div>
    <div>
      <label for="rate" class="block text-sm font-medium text-gray-700">翻译速率</label>
      <input type="number" id="rate" v-model="settings.rate" class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
      <p class="mt-2 text-sm text-gray-500">速率过高可能导致 API 返回空值, 请适当降低速率。</p>
    </div>
    <div class="flex justify-end space-x-2">
      <button @click="clearCache" class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">清除翻译缓存</button>
      <button @click="resetSettings" class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">恢复默认设置</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { getApiSettings, saveApiSettings, clearTranslationCache } from '../services/apiSettingsService';
import type { ApiSettings } from '../services/apiSettingsService';

const props = defineProps<{
  apiName: string;
}>();

const settings = ref<ApiSettings>({
  apiUrl: '',
  apiKey: '',
  apiSecret: '',
  chunkSize: 1000,
  delay: 200,
  rate: 10,
});

const loadSettings = () => {
  settings.value = getApiSettings(props.apiName);
};

const saveSettings = () => {
  saveApiSettings(props.apiName, settings.value);
};

const clearCache = () => {
  clearTranslationCache();
  alert('翻译缓存已清除');
};

const resetSettings = () => {
  settings.value = {
    apiUrl: '',
    apiKey: '',
    apiSecret: '',
    chunkSize: 1000,
    delay: 200,
    rate: 10,
  };
  saveSettings();
};

onMounted(() => {
  loadSettings();
});

watch(settings, saveSettings, { deep: true });

</script>