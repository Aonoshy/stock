<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
    <!-- 模版管理区域 -->
    <div class="mb-6">
      <div class="flex justify-between items-center mb-4 px-4">
        <h3 class="text-lg font-medium text-gray-700">我的查询模版</h3>
        <button
          @click="showCreateTemplate = true"
          class="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          + 新建模版
        </button>
      </div>

      <!-- Tab栏 -->
      <div v-if="templates.length > 0" class="border-b border-gray-200">
        <div class="flex space-x-1 px-2">
          <button
            v-for="template in templates"
            :key="template.id"
            @click="selectTemplate(template)"
            :class="[
              'px-4 py-3 text-sm font-medium rounded-t-lg transition-colors',
              activeTemplate?.id === template.id
                ? 'bg-blue-600 text-white border-b-2 border-blue-600'
                : 'text-gray-600 hover:text-blue-600 hover:bg-gray-50'
            ]"
          >
            {{ template.name }}
          </button>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="templates.length === 0" class="text-center py-8 text-gray-500">
        <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        <p class="text-sm">还没有保存的模版</p>
        <p class="text-xs mt-1">点击"新建模版"创建第一个股票组合</p>
      </div>

      <!-- 创建模版表单 -->
      <div v-if="showCreateTemplate" class="mt-4 p-4 bg-white border border-gray-200 rounded-md">
        <h4 class="text-sm font-medium text-gray-700 mb-3">创建新模版</h4>

        <div class="space-y-3">
          <div>
            <label class="block text-xs text-gray-600 mb-1">模版名称</label>
            <input
              v-model="newTemplate.name"
              type="text"
              placeholder="例如：我的自选股"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:border-gray-400"
            />
          </div>

          <div>
            <label class="block text-xs text-gray-600 mb-1">模版描述</label>
            <input
              v-model="newTemplate.description"
              type="text"
              placeholder="例如：长期关注的优质股票"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:border-gray-400"
            />
          </div>

          <div>
            <label class="block text-xs text-gray-600 mb-1">当前查询的股票将保存到此模版</label>
            <div class="text-xs text-gray-500">
              {{ getValidQueries().map(q => `${q.code}.${q.exchange}`).join(', ') || '请先填写股票代码' }}
            </div>
          </div>

          <div class="flex gap-2 pt-2">
            <button
              @click="saveTemplate"
              :disabled="!canSaveTemplate"
              class="flex-1 py-2 bg-gray-600 text-white rounded text-sm hover:bg-gray-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              保存模版
            </button>
            <button
              @click="cancelCreateTemplate"
              class="flex-1 py-2 bg-gray-300 text-gray-700 rounded text-sm hover:bg-gray-400 transition-colors"
            >
              取消
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 模版查询预览区域 -->
    <div v-if="showCreateTemplate" class="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
      <h3 class="text-sm font-medium text-blue-800 mb-3">准备创建模版的股票</h3>

      <!-- 动态查询框列表 -->
      <div class="space-y-3">
        <div
          v-for="(query, index) in queries"
          :key="query.id"
          class="flex items-center gap-3 p-3 bg-white rounded-lg border"
        >
          <span class="text-sm text-gray-500 w-4">{{ index + 1 }}</span>

          <!-- 股票代码输入 -->
          <input
            v-model="query.code"
            type="text"
            placeholder="股票代码"
            class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-400 focus:border-blue-400"
          />

          <!-- 交易所选择 -->
          <select
            v-model="query.exchange"
            class="w-20 px-2 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-400 focus:border-blue-400"
          >
            <option value="SZ">SZ</option>
            <option value="SH">SH</option>
          </select>

          <!-- 删除按钮 -->
          <button
            v-if="queries.length > 1"
            @click="removeQuery(index)"
            class="p-2 text-gray-400 hover:text-red-500 transition-colors"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>

          <!-- 占位符，保持对齐 -->
          <div v-else class="w-8"></div>
        </div>

        <!-- 添加股票按钮 -->
        <button
          v-if="queries.length < 5"
          @click="addQuery"
          class="w-full flex items-center justify-center gap-2 py-3 border-2 border-dashed border-blue-300 rounded-lg text-blue-600 hover:border-blue-400 hover:bg-blue-50 transition-colors"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
          添加股票 (最多5只)
        </button>
      </div>
    </div>

    <!-- 错误信息 -->
    <div v-if="error" class="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
      <p class="text-sm text-red-700">{{ error }}</p>
    </div>

    <!-- Tab内容区域 -->
    <div v-if="activeTemplate && !showCreateTemplate" class="mt-6">
      <!-- 模版详情头部 -->
      <div class="bg-gray-50 rounded-lg p-4 mb-4">
        <div class="flex justify-between items-start">
          <div class="flex-1">
            <h4 class="text-lg font-semibold text-gray-800 mb-1">{{ activeTemplate.name }}</h4>
            <p class="text-sm text-gray-600 mb-2">{{ activeTemplate.description }}</p>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="code in activeTemplate.codes"
                :key="code"
                class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-md"
              >
                {{ code }}
              </span>
            </div>
          </div>
          <div class="flex gap-2 ml-4">
            <button
              @click="queryTemplate(activeTemplate.id)"
              :disabled="loading"
              class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 text-sm"
            >
              <span v-if="loading && currentQueryingTemplate === activeTemplate.id" class="flex items-center gap-2">
                <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                查询中...
              </span>
              <span v-else>🔍 查询股价</span>
            </button>
            <button
              @click="deleteTemplate(activeTemplate.id)"
              :disabled="loading"
              class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 text-sm"
            >
              🗑️ 删除
            </button>
          </div>
        </div>
      </div>

      <!-- 查询结果 -->
      <div v-if="results.length > 0 && lastQueriedTemplate === activeTemplate.id" class="space-y-3">
        <h5 class="text-base font-medium text-gray-700 pb-2 border-b border-gray-200">
          {{ activeTemplate.name }} - 查询结果
        </h5>

        <div class="grid gap-3">
          <div
            v-for="stock in results"
            :key="stock.code"
            class="flex justify-between items-center p-4 bg-gradient-to-r from-green-50 to-blue-50 rounded-lg border border-green-200"
          >
            <div>
              <div class="font-medium text-gray-800">{{ stock.code }}</div>
              <div class="text-xs text-gray-500 mt-1">{{ formatTime(stock.timestamp) }}</div>
            </div>
            <div class="text-right">
              <div class="text-xl font-bold text-green-600">¥{{ stock.price.toFixed(2) }}</div>
            </div>
          </div>
        </div>

        <div class="mt-4 flex justify-between items-center pt-3 border-t border-gray-200">
          <div class="text-sm text-gray-600">
            📊 共 {{ results.length }} 只股票 · 总市值约 ¥{{ totalValue.toLocaleString() }} (按1手计算)
          </div>
          <button
            @click="clearResults"
            class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
          >
            清空结果
          </button>
        </div>
      </div>

      <!-- 无结果提示 -->
      <div v-else-if="!loading" class="text-center py-8 text-gray-500">
        <svg class="mx-auto h-8 w-8 text-gray-400 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <p class="text-sm">点击"查询股价"获取最新价格</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useStockStore, type StockData } from '@/stores/stock'
import { stockApi, type StockTemplate } from '@/services/stockApi'

interface StockQuery {
  id: number
  code: string
  exchange: 'SZ' | 'SH'
}

const stockStore = useStockStore()
const loading = computed(() => stockStore.loading)
const error = computed(() => stockStore.error)

let nextId = 1

// 查询列表，默认一个查询框
const queries = ref<StockQuery[]>([
  { id: nextId++, code: '', exchange: 'SZ' }
])

const results = ref<StockData[]>([])
const templates = ref<StockTemplate[]>([])
const activeTemplate = ref<StockTemplate | null>(null)
const showCreateTemplate = ref(false)
const currentQueryingTemplate = ref<string | null>(null)
const lastQueriedTemplate = ref<string | null>(null)
const newTemplate = ref({
  name: '',
  description: ''
})

// 检查是否有有效的查询
const hasValidQueries = computed(() => {
  return queries.value.some(query => query.code.trim().length > 0)
})

// 获取有效的查询
const getValidQueries = () => {
  return queries.value.filter(query => query.code.trim().length > 0)
}

// 检查是否可以保存模版
const canSaveTemplate = computed(() => {
  return newTemplate.value.name.trim().length > 0 && getValidQueries().length > 0
})

// 计算总市值
const totalValue = computed(() => {
  return results.value.reduce((sum, stock) => sum + stock.price * 100, 0) // 1手=100股
})

// 选择模版
function selectTemplate(template: StockTemplate) {
  activeTemplate.value = template
  showCreateTemplate.value = false
  // 清空之前的查询结果，除非是当前模版的结果
  if (lastQueriedTemplate.value !== template.id) {
    results.value = []
  }
}

// 加载模版列表
async function loadTemplates() {
  try {
    const response = await stockApi.getTemplates()
    templates.value = response.templates

    // 如果有模版且当前没有选中的模版，自动选择第一个
    if (templates.value.length > 0 && !activeTemplate.value) {
      activeTemplate.value = templates.value[0]
    }

    console.log('加载模版列表:', templates.value)
  } catch (err) {
    console.error('加载模版失败:', err)
  }
}

// 查询模版股票
async function queryTemplate(templateId: string) {
  try {
    currentQueryingTemplate.value = templateId
    console.log('查询模版:', templateId)

    const response = await stockApi.queryTemplate(templateId)
    console.log('模版查询结果:', response)

    // 将模版查询结果转换为StockData格式
    const templateResults: StockData[] = Object.entries(response.data).map(([code, price]) => ({
      code,
      price,
      timestamp: Date.now()
    }))

    results.value = templateResults
    lastQueriedTemplate.value = templateId
    console.log('设置模版结果:', results.value)
  } catch (err) {
    console.error('模版查询失败:', err)
  } finally {
    currentQueryingTemplate.value = null
  }
}

// 添加查询框
function addQuery() {
  if (queries.value.length < 5) {
    queries.value.push({
      id: nextId++,
      code: '',
      exchange: 'SZ'
    })
  }
}

// 删除查询框
function removeQuery(index: number) {
  if (queries.value.length > 1) {
    queries.value.splice(index, 1)
  }
}

// 清空结果并重置查询表单
function resetQueries() {
  queries.value = [{ id: nextId++, code: '', exchange: 'SZ' }]
  results.value = []
  stockStore.clearAll()
}

// 清空结果
function clearResults() {
  results.value = []
  stockStore.clearAll()
}

// 格式化时间
function formatTime(timestamp: number): string {
  return new Date(timestamp).toLocaleTimeString('zh-CN')
}

// 保存模版
async function saveTemplate() {
  const validQueries = getValidQueries()
  if (!newTemplate.value.name.trim() || validQueries.length === 0) {
    return
  }

  try {
    // 生成模版ID（基于名称和时间戳）
    const templateId = newTemplate.value.name.toLowerCase().replace(/[^a-z0-9]/g, '_') + '_' + Date.now()

    const stockCodes = validQueries.map(query => `${query.code.trim()}.${query.exchange}`)

    await stockApi.createTemplate({
      id: templateId,
      name: newTemplate.value.name.trim(),
      description: newTemplate.value.description.trim(),
      codes: stockCodes
    })

    // 重新加载模版列表
    await loadTemplates()

    // 重置表单和查询框
    cancelCreateTemplate()
    resetQueries()

    console.log('模版保存成功')
  } catch (err) {
    console.error('保存模版失败:', err)
    // 这里可以显示错误提示
  }
}

// 取消创建模版
function cancelCreateTemplate() {
  showCreateTemplate.value = false
  newTemplate.value.name = ''
  newTemplate.value.description = ''
}

// 删除模版
async function deleteTemplate(templateId: string) {
  try {
    await stockApi.deleteTemplate(templateId)

    // 如果删除的是当前选中的模版，清空选中状态
    if (activeTemplate.value?.id === templateId) {
      activeTemplate.value = null
      results.value = []
      lastQueriedTemplate.value = null
    }

    // 重新加载模版列表
    await loadTemplates()

    console.log('模版删除成功')
  } catch (err) {
    console.error('删除模版失败:', err)
  }
}

// 组件挂载时加载模版列表
onMounted(() => {
  loadTemplates()
})
</script>