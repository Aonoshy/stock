<template>
  <div class="bg-white rounded border p-4">
    <h3 class="text-base font-medium text-gray-800 mb-3">单股查询</h3>

    <form @submit.prevent="handleSubmit" class="flex gap-2">
      <input
        v-model="stockCode"
        type="text"
        placeholder="000001.SZ"
        class="flex-1 px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:border-blue-500"
        :disabled="loading"
      />
      <button
        type="submit"
        :disabled="!stockCode.trim() || loading"
        class="px-4 py-2 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 disabled:bg-gray-400"
      >
        {{ loading ? '查询中...' : '查询' }}
      </button>
    </form>

    <div v-if="error" class="mt-3 p-2 bg-red-50 text-red-700 text-sm rounded">
      {{ error }}
    </div>

    <div v-if="result" class="mt-3 p-3 bg-gray-50 rounded">
      <div class="flex justify-between items-center">
        <span class="font-medium">{{ result.code }}</span>
        <span class="text-lg font-semibold text-green-600">¥{{ result.price.toFixed(2) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useStockStore } from '@/stores/stock'

const stockStore = useStockStore()
const { loading, error } = stockStore

const stockCode = ref('')
const result = ref<any>(null)

async function handleSubmit() {
  if (!stockCode.value.trim()) return

  await stockStore.fetchSingleStock(stockCode.value.trim())

  if (!error.value) {
    result.value = stockStore.getStock(stockCode.value.trim())
  }
}

function formatTime(timestamp: number): string {
  return new Date(timestamp).toLocaleTimeString('zh-CN')
}
</script>