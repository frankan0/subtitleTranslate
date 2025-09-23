<template>
  <div class="min-h-screen bg-background text-foreground">
    <header class="bg-background border-b">
      <nav class="container mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-bold">字幕翻译工具</h1>
          </div>
          <div class="flex items-center space-x-4">
            <router-link to="/" class="text-sm font-medium text-muted-foreground hover:text-primary">翻译工作区</router-link>
            <router-link to="/settings" class="text-sm font-medium text-muted-foreground hover:text-primary">API 设置</router-link>
          </div>
        </div>
      </nav>
    </header>
    <main class="container py-10">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
// This component now only handles the main layout and navigation.
</script>