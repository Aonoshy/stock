package main

import (
	"net/http"
	"stock-api/config"
	"stock-api/service"

	"github.com/gin-gonic/gin"
)

type StockResponse struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

type BatchStockRequest struct {
	Codes []string `json:"codes"`
}

type BatchStockResponse struct {
	Data map[string]float64 `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TemplateResponse struct {
	Templates []*service.StockTemplate `json:"templates"`
}

type CreateTemplateRequest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Codes       []string `json:"codes"`
}

type TemplateQueryResponse struct {
	Template *service.StockTemplate `json:"template"`
	Data     map[string]float64     `json:"data"`
}

var stockService *service.StockService
var templateService *service.TemplateService

func main() {
	// 加载配置
	cfg := config.NewConfig()

	// 初始化股票服务
	stockService = service.NewStockService(cfg.GetAlltickAPIKey())

	// 初始化模版服务
	templateService = service.NewTemplateService()

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API路由组
	api := router.Group("/api/v1")
	{
		api.GET("/stock/:code", getStockPrice)
		api.POST("/stocks/batch", getBatchStockPrices)

		// 模版相关API
		api.GET("/templates", getTemplates)
		api.POST("/templates", createTemplate)
		api.DELETE("/templates/:id", deleteTemplate)
		api.GET("/templates/:id", getTemplate)
		api.GET("/templates/:id/query", queryTemplateStocks)
	}

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 启动服务器
	router.Run(":" + cfg.GetPort())
}

// 获取股票价格的处理函数
func getStockPrice(c *gin.Context) {
	stockCode := c.Param("code")

	if stockCode == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "股票代码不能为空",
		})
		return
	}

	// 调用股票服务获取真实价格
	price, err := stockService.GetStockPrice(stockCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, StockResponse{
		Code:  stockCode,
		Price: price,
	})
}

// 批量获取股票价格的处理函数
func getBatchStockPrices(c *gin.Context) {
	var request BatchStockRequest

	// 解析JSON请求体
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "请求格式错误: " + err.Error(),
		})
		return
	}

	// 验证股票代码列表
	if len(request.Codes) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "股票代码列表不能为空",
		})
		return
	}

	// 限制单次查询数量（避免URL过长）
	if len(request.Codes) > 50 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "单次查询股票数量不能超过50只",
		})
		return
	}

	// 调用股票服务获取批量价格
	prices, err := stockService.GetBatchStockPrices(request.Codes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, BatchStockResponse{
		Data: prices,
	})
}

// 获取所有模版
func getTemplates(c *gin.Context) {
	templates := templateService.GetAllTemplates()
	c.JSON(http.StatusOK, TemplateResponse{
		Templates: templates,
	})
}

// 获取指定模版
func getTemplate(c *gin.Context) {
	templateID := c.Param("id")

	template, err := templateService.GetTemplate(templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// 创建新模版
func createTemplate(c *gin.Context) {
	var request CreateTemplateRequest

	// 解析JSON请求体
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "请求格式错误: " + err.Error(),
		})
		return
	}

	template := &service.StockTemplate{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
		Codes:       request.Codes,
	}

	if err := templateService.CreateTemplate(template); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// 删除模版
func deleteTemplate(c *gin.Context) {
	templateID := c.Param("id")

	if err := templateService.DeleteTemplate(templateID); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "模版已删除"})
}

// 根据模版查询股票价格
func queryTemplateStocks(c *gin.Context) {
	templateID := c.Param("id")

	template, err := templateService.GetTemplate(templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	prices, err := templateService.QueryTemplateStocks(templateID, stockService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TemplateQueryResponse{
		Template: template,
		Data:     prices,
	})
}