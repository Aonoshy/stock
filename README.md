# 股票价格查询系统

全栈股票价格查询应用，支持A股实时价格查询。

## 项目架构

- **后端**: Go + Gin + Alltick API
- **前端**: Vue 3 + Vite + TypeScript + Tailwind CSS + Headless UI

## 项目结构

```
stock/
├── backend/           # Go 后端服务
│   ├── main.go       # 主程序入口
│   ├── service/      # 股票数据服务
│   ├── config/       # 配置管理
│   ├── go.mod        # Go 依赖
│   └── Dockerfile    # 后端 Docker 配置
├── frontend/         # Vue 前端应用
│   └── stock-frontend/
│       ├── src/      # 源代码
│       ├── public/   # 静态资源
│       └── package.json
└── README.md
```

## API 接口

### 单个股票查询
```bash
GET /api/v1/stock/{code}
# 例如: curl http://localhost:8080/api/v1/stock/000001.SZ
```

### 批量股票查询
```bash
POST /api/v1/stocks/batch
Content-Type: application/json

{
  "codes": ["000001.SZ", "600036.SH"]
}
```

## 股票代码格式

- **深圳股票**: `000001.SZ` (平安银行)
- **上海股票**: `600036.SH` (招商银行)

## 开发指南

### 后端开发
```bash
cd backend
export ALLTICK_API_KEY="your-api-key"
go run main.go
```

### 前端开发
```bash
cd frontend/stock-frontend
npm install
npm run dev
```

### Docker 部署
```bash
cd backend
docker build -t stock-api .
docker run -p 8080:8080 -e ALLTICK_API_KEY="your-api-key" stock-api
```

## 环境变量

- `ALLTICK_API_KEY`: Alltick API 密钥
- `PORT`: 服务端口 (默认: 8080)