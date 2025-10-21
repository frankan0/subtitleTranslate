<template>
  <div class="space-y-4">
    <div>
      <label for="api-key" class="block text-sm font-medium text-gray-700">API 密钥</label>
      <input type="text" id="api-key" v-model="settings.apiKey" placeholder="输入API密钥" class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
      <p class="mt-2 text-sm text-gray-500">API密钥将用于访问翻译服务</p>
    </div>
    <div>
      <label for="api-secret" class="block text-sm font-medium text-gray-700">API Secret</label>
      <input type="text" id="api-secret" v-model="settings.apiSecret" placeholder="输入API Secret" class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
      <p class="mt-2 text-sm text-gray-500">部分API服务需要Secret密钥</p>
    </div>

    <div class="flex justify-end space-x-2">
      <button @click="clearCache" class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">清除翻译缓存</button>
      <button @click="resetSettings" class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">恢复默认设置</button>
      <button @click="saveSettings" class="px-4 py-2 bg-blue-600 text-white rounded-md shadow-sm text-sm font-medium hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">保存配置</button>
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

const hasUnsavedChanges = ref(false);

const loadSettings = () => {
  settings.value = getApiSettings(props.apiName);
  hasUnsavedChanges.value = false;
};

const saveSettings = () => {
  saveApiSettings(props.apiName, settings.value);
  hasUnsavedChanges.value = false;
  alert('配置已保存');
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
  hasUnsavedChanges.value = true;
};

onMounted(() => {
  loadSettings();
});

// 监听apiName变化，重新加载配置
watch(() => props.apiName, () => {
  loadSettings();
});

// 监听配置变化，标记为未保存
watch(settings, () => {
  hasUnsavedChanges.value = true;
}, { deep: true });

</script>