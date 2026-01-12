<template>
  <div class="bg-white rounded border p-4">
    <h3 class="text-base font-medium text-gray-800 mb-3">批量查询</h3>

    <div class="space-y-3">
      <textarea
        v-model="stockCodesText"
        rows="3"
        placeholder="000001.SZ&#10;600036.SH"
        class="w-full px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:border-blue-500 resize-none"
        :disabled="loading"
      />

      <div class="flex gap-2 text-xs">
        <button
          v-for="preset in presetStocks"
          :key="preset.name"
          type="button"
          @click="loadPreset(preset.codes)"
          class="px-2 py-1 bg-gray-100 text-gray-600 rounded hover:bg-gray-200"
        >
          {{ preset.name }}
        </button>
      </div>

      <button
        @click="handleSubmit"
        :disabled="!stockCodes.length || loading"
        class="w-full py-2 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 disabled:bg-gray-400"
      >
        {{ loading ? '查询中...' : `批量查询 (${stockCodes.length}只)` }}
      </button>
    </div>

    <div v-if="error" class="mt-3 p-2 bg-red-50 text-red-700 text-sm rounded">
      {{ error }}
    </div>

    <div v-if="results.length > 0" class="mt-4">
      <div class="space-y-1">
        <div
          v-for="stock in results"
          :key="stock.code"
          class="flex justify-between items-center py-1 px-2 bg-gray-50 rounded text-sm"
        >
          <span>{{ stock.code }}</span>
          <span class="font-medium text-green-600">¥{{ stock.price.toFixed(2) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useStockStore, type StockData } from '@/stores/stock'

const stockStore = useStockStore()
const { loading, error } = stockStore

const stockCodesText = ref('')
const results = ref<StockData[]>([])

const presetStocks = [
  {
    name: '银行股',
    codes: ['000001.SZ', '600036.SH', '000002.SZ', '600000.SH']
  },
  {
    name: '科技股',
    codes: ['000858.SZ', '002415.SZ', '000725.SZ']
  },
  {
    name: '白酒股',
    codes: ['000858.SZ', '000596.SZ']
  }
]

const stockCodes = computed(() => {
  return stockCodesText.value
    .split('\n')
    .map(code => code.trim())
    .filter(code => code.length > 0)
    .slice(0, 50) // 限制最多50只
})

const totalValue = computed(() => {
  return results.value.reduce((sum, stock) => sum + stock.price * 100, 0) // 1手=100股
})

async function handleSubmit() {
  if (stockCodes.value.length === 0) return

  await stockStore.fetchBatchStocks(stockCodes.value)

  if (!error.value) {
    results.value = stockCodes.value
      .map(code => stockStore.getStock(code))
      .filter(stock => stock !== undefined) as StockData[]
  }
}

function loadPreset(codes: string[]) {
  stockCodesText.value = codes.join('\n')
}

function clearResults() {
  results.value = []
  stockStore.clearAll()
}

function formatTime(timestamp: number): string {
  return new Date(timestamp).toLocaleTimeString('zh-CN')
}
</script>