import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { stockApi, type StockPrice } from '@/services/stockApi'

export interface StockData extends StockPrice {
  timestamp: number
}

export const useStockStore = defineStore('stock', () => {
  const stocks = ref<Map<string, StockData>>(new Map())
  const loading = ref(false)
  const error = ref<string | null>(null)

  const stockList = computed(() => Array.from(stocks.value.values()))

  async function fetchSingleStock(code: string) {
    loading.value = true
    error.value = null

    try {
      const stock = await stockApi.getStockPrice(code)
      stocks.value.set(code, {
        ...stock,
        timestamp: Date.now()
      })
    } catch (err) {
      error.value = err instanceof Error ? err.message : '未知错误'
    } finally {
      loading.value = false
    }
  }

  async function fetchBatchStocks(codes: string[]) {
    loading.value = true
    error.value = null

    try {
      const response = await stockApi.getBatchStockPrices(codes)

      Object.entries(response.data).forEach(([code, price]) => {
        stocks.value.set(code, {
          code,
          price,
          timestamp: Date.now()
        })
      })
    } catch (err) {
      error.value = err instanceof Error ? err.message : '未知错误'
    } finally {
      loading.value = false
    }
  }

  function getStock(code: string): StockData | undefined {
    return stocks.value.get(code)
  }

  function removeStock(code: string) {
    stocks.value.delete(code)
  }

  function clearAll() {
    stocks.value.clear()
    error.value = null
  }

  return {
    stocks,
    stockList,
    loading,
    error,
    fetchSingleStock,
    fetchBatchStocks,
    getStock,
    removeStock,
    clearAll
  }
})