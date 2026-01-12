<script setup lang="ts">
import StockQueryPanel from '@/components/StockQueryPanel.vue'
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-4xl mx-auto px-4 py-8">
      <!-- 页面标题 -->
      <div class="text-center mb-8">
        <h1 class="text-2xl font-light text-gray-700 mb-2">股票价格查询</h1>
        <p class="text-sm text-gray-500">实时A股价格查询系统</p>
      </div>

      <!-- 查询面板 -->
      <StockQueryPanel />
    </div>
  </div>
</template>
