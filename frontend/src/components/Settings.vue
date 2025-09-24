<template>
  <div class="flex h-screen">
    <div class="w-1/4 bg-gray-800 text-white p-4">
      <h2 class="text-lg font-bold mb-4">API 设置</h2>
      <ul>
        <li v-for="api in apis" :key="api.id" @click="selectedApi = api.id" :class="{ 'bg-gray-700': selectedApi === api.id }" class="cursor-pointer p-2 rounded">
          {{ api.name }}
        </li>
      </ul>
    </div>
    <div class="w-3/4 p-8">
      <div v-if="selectedApi">
        <h2 class="text-2xl font-bold mb-4">{{ selectedApiName }} API 配置</h2>
        <APISettings :api-name="selectedApi" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import APISettings from './APISettings.vue';
import { getCurrentApiName, setCurrentApiName } from '../services/apiSettingsService';

const apis = ref([
  { id: 'volce', name: '火山引擎' },
  { id: 'google', name: 'Google Translate' },
  { id: 'deepl', name: 'DeepL' },
  { id: 'deepseek', name: 'DeepSeek' },
  { id: 'gemini', name: 'Gemini' },
]);

const selectedApi = ref(getCurrentApiName());

const selectedApiName = computed(() => {
  const api = apis.value.find(a => a.id === selectedApi.value);
  return api ? api.name : '';
});

watch(selectedApi, (newValue) => {
  setCurrentApiName(newValue);
});

</script>