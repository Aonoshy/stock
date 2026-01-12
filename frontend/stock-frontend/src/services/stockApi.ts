export interface StockPrice {
  code: string
  price: number
}

export interface BatchStockResponse {
  data: Record<string, number>
}

export interface ApiError {
  error: string
}

export interface StockTemplate {
  id: string
  name: string
  description: string
  codes: string[]
  created_at: string
}

export interface TemplateResponse {
  templates: StockTemplate[]
}

export interface TemplateQueryResponse {
  template: StockTemplate
  data: Record<string, number>
}

class StockApiService {
  private baseUrl = 'http://58783abe.r40.cpolar.top/api/v1'

  async getStockPrice(code: string): Promise<StockPrice> {
    const url = `${this.baseUrl}/stock/${code}`
    console.log('发送请求到:', url)

    const response = await fetch(url)
    console.log('响应状态:', response.status, response.statusText)

    if (!response.ok) {
      const error: ApiError = await response.json()
      console.error('API错误:', error)
      throw new Error(error.error || '获取股票价格失败')
    }

    const result = await response.json()
    console.log('API响应:', result)
    return result
  }

  async getBatchStockPrices(codes: string[]): Promise<BatchStockResponse> {
    const response = await fetch(`${this.baseUrl}/stocks/batch`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ codes })
    })

    if (!response.ok) {
      const error: ApiError = await response.json()
      throw new Error(error.error || '获取批量股票价格失败')
    }

    return response.json()
  }

  async getTemplates(): Promise<TemplateResponse> {
    const response = await fetch(`${this.baseUrl}/templates`)

    if (!response.ok) {
      const error: ApiError = await response.json()
      throw new Error(error.error || '获取模版列表失败')
    }

    return response.json()
  }

  async queryTemplate(templateId: string): Promise<TemplateQueryResponse> {
    const response = await fetch(`${this.baseUrl}/templates/${templateId}/query`)

    if (!response.ok) {
      const error: ApiError = await response.json()
      throw new Error(error.error || '查询模版股票失败')
    }

    return response.json()
  }

  async createTemplate(template: {
    id: string
    name: string
    description: string
    codes: string[]
  }): Promise<StockTemplate> {
    const response = await fetch(`${this.baseUrl}/templates`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(template)
    })

    if (!response.ok) {
      const error: ApiError = await response.json()
      throw new Error(error.error || '创建模版失败')
    }

    return response.json()
  }

  async deleteTemplate(templateId: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/templates/${templateId}`, {
      method: 'DELETE'
    })

    if (!response.ok) {
      const error: ApiError = await response.json()
      throw new Error(error.error || '删除模版失败')
    }
  }
}

export const stockApi = new StockApiService()