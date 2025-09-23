<template>
  <div class="w-full">
    <div
      class="border-2 border-dashed border-border rounded-lg p-8 text-center hover:border-primary transition-colors"
      :class="{ 'border-primary bg-primary/5': isDragging }"
      @dragover.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @drop.prevent="onDrop"
    >
      <div class="flex flex-col items-center justify-center gap-4">
        <div class="rounded-full bg-primary/10 p-3">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-primary">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="17 8 12 3 7 8"></polyline>
            <line x1="12" y1="3" x2="12" y2="15"></line>
          </svg>
        </div>
        <div class="space-y-2">
          <h3 class="text-lg font-medium">拖放文件到此处或点击上传</h3>
          <p class="text-sm text-muted-foreground">支持 .srt, .vtt, .ass, .ssa 格式的字幕文件</p>
        </div>
        <input
          type="file"
          ref="fileInput"
          multiple
          accept=".srt,.vtt,.ass,.ssa,text/plain"
          class="hidden"
          @change="onFileChange"
        />
        <button
          type="button"
          class="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors"
          @click="$refs.fileInput.click()"
        >
          选择文件
        </button>
      </div>
    </div>

    <div v-if="errorMessage" class="mt-4 p-4 bg-destructive/10 text-destructive rounded-md">
      {{ errorMessage }}
    </div>


  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { SubtitleFile } from '../types';

interface FileStats {
  name: string;
  charCount: number;
}

const emit = defineEmits<{
  (e: 'files-uploaded', files: SubtitleFile[]): void
}>();

const fileInput = ref<HTMLInputElement | null>(null);
const isDragging = ref(false);
const errorMessage = ref('');
const fileStats = ref<FileStats[]>([]);

const onDrop = (event: DragEvent) => {
  isDragging.value = false;
  if (!event.dataTransfer) return;
  
  const files = Array.from(event.dataTransfer.files);
  processFiles(files);
};

const onFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (!target.files) return;
  
  const files = Array.from(target.files);
  processFiles(files);
  
  // 重置文件输入，允许重新选择相同的文件
  if (fileInput.value) {
    fileInput.value.value = '';
  }
};

const processFiles = async (files: File[]) => {
  errorMessage.value = '';
  fileStats.value = [];
  
  // 验证文件类型
  const invalidFiles = files.filter(file => {
    const extension = file.name.split('.').pop()?.toLowerCase();
    return !['srt', 'vtt', 'ass', 'ssa'].includes(extension || '');
  });
  
  if (invalidFiles.length > 0) {
    errorMessage.value = `以下文件不是有效的字幕格式 (仅支持 .srt, .vtt, .ass, .ssa): ${invalidFiles.map(f => f.name).join(', ')}`;
    return;
  }
  
  // 读取文件内容
  const subtitleFiles: SubtitleFile[] = [];
  
  for (const file of files) {
    try {
      const content = await readFileAsText(file);
      // 统计字符数
      const charCount = content.replace(/\s/g, '').length;
      fileStats.value.push({
        name: file.name,
        charCount
      });

      subtitleFiles.push({
        id: generateId(),
        name: file.name,
        content,
        size: file.size,
        type: file.type || 'text/plain',
        lastModified: file.lastModified,
        charCount: content.replace(/\s/g, '').length
      });
    } catch (error) {
      console.error(`读取文件 ${file.name} 失败:`, error);
      errorMessage.value = `读取文件 ${file.name} 失败`;
    }
  }
  
  if (subtitleFiles.length > 0) {
    emit('files-uploaded', subtitleFiles);
  }
};

const readFileAsText = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = () => reject(new Error(`读取文件 ${file.name} 失败`));
    reader.readAsText(file);
  });
};

const generateId = (): string => {
  return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
};

const formatCharCount = (count: number): string => {
  if (count < 1000) {
    return `${count} 字符`;
  } else if (count < 10000) {
    return `${(count / 1000).toFixed(1)}k 字符`;
  } else {
    return `${(count / 10000).toFixed(1)}w 字符`;
  }
};
</script>