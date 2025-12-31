package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type StockService struct {
	client    *http.Client
	apiKey    string
	baseURL   string
}

func NewStockService(apiKey string) *StockService {
	return &StockService{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: "https://quote.alltick.co/quote-stock-b-api",
	}
}

// 获取股票价格的接口
func (s *StockService) GetStockPrice(code string) (float64, error) {
	if s.apiKey == "" {
		return 0, fmt.Errorf("API密钥未配置")
	}

	// 使用Alltick API获取A股最新交易价格
	// 构建查询JSON
	queryRequest := AlltickRequest{
		Trace: "golang-api-" + fmt.Sprintf("%d", time.Now().UnixNano()),
		Data: AlltickRequestData{
			SymbolList: []AlltickSymbol{
				{Code: code},
			},
		},
	}

	// 序列化为JSON
	queryJSON, err := json.Marshal(queryRequest)
	if err != nil {
		return 0, fmt.Errorf("序列化查询参数失败: %v", err)
	}

	// URL编码JSON查询参数
	encodedQuery := url.QueryEscape(string(queryJSON))

	// 构建完整的URL
	apiURL := fmt.Sprintf("%s/trade-tick?token=%s&query=%s", s.baseURL, s.apiKey, encodedQuery)

	// 创建GET请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return 0, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加必要的请求头
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析Alltick API响应
	var result AlltickTradeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("解析JSON失败: %v, 响应内容: %s", err, string(body[:200]))
	}

	// 检查响应状态
	if result.Ret != 200 {
		return 0, fmt.Errorf("API错误: %s (状态码: %d)", result.Msg, result.Ret)
	}

	// 检查数据是否存在
	if result.Data == nil || len(result.Data.TickList) == 0 {
		return 0, fmt.Errorf("未找到股票代码 '%s' 的交易数据", code)
	}

	// 获取最新价格并转换为浮点数
	tick := result.Data.TickList[0]
	price, err := strconv.ParseFloat(tick.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("价格转换失败: %v", err)
	}

	return price, nil
}

// 批量获取股票价格的接口
func (s *StockService) GetBatchStockPrices(codes []string) (map[string]float64, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("API密钥未配置")
	}

	if len(codes) == 0 {
		return make(map[string]float64), nil
	}

	// 构建批量查询请求
	symbols := make([]AlltickSymbol, len(codes))
	for i, code := range codes {
		symbols[i] = AlltickSymbol{Code: code}
	}

	queryRequest := AlltickRequest{
		Trace: "golang-batch-api-" + fmt.Sprintf("%d", time.Now().UnixNano()),
		Data: AlltickRequestData{
			SymbolList: symbols,
		},
	}

	// 序列化为JSON
	queryJSON, err := json.Marshal(queryRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化查询参数失败: %v", err)
	}

	// URL编码JSON查询参数
	encodedQuery := url.QueryEscape(string(queryJSON))

	// 构建完整的URL
	apiURL := fmt.Sprintf("%s/trade-tick?token=%s&query=%s", s.baseURL, s.apiKey, encodedQuery)

	// 创建GET请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加必要的请求头
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析Alltick API响应
	var result AlltickTradeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v, 响应内容: %s", err, string(body[:200]))
	}

	// 检查响应状态
	if result.Ret != 200 {
		return nil, fmt.Errorf("API错误: %s (状态码: %d)", result.Msg, result.Ret)
	}

	// 检查数据是否存在
	if result.Data == nil {
		return make(map[string]float64), nil
	}

	// 构建结果映射
	prices := make(map[string]float64)
	for _, tick := range result.Data.TickList {
		price, err := strconv.ParseFloat(tick.Price, 64)
		if err != nil {
			// 如果某个价格转换失败，记录错误但继续处理其他股票
			fmt.Printf("价格转换失败 %s: %v\n", tick.Code, err)
			continue
		}
		prices[tick.Code] = price
	}

	return prices, nil
}

// Alltick API请求结构体
type AlltickRequest struct {
	Trace string              `json:"trace"`
	Data  AlltickRequestData `json:"data"`
}

type AlltickRequestData struct {
	SymbolList []AlltickSymbol `json:"symbol_list"`
}

type AlltickSymbol struct {
	Code string `json:"code"`
}

// Alltick API响应结构体
type AlltickTradeResponse struct {
	Ret   int    `json:"ret"`
	Msg   string `json:"msg"`
	Trace string `json:"trace"`
	Data  *struct {
		TickList []struct {
			Code           string `json:"code"`
			Seq            string `json:"seq"`
			TickTime       string `json:"tick_time"`
			Price          string `json:"price"`          // API返回字符串，需要转换
			Volume         string `json:"volume"`
			Turnover       string `json:"turnover"`
			TradeDirection int    `json:"trade_direction"`
		} `json:"tick_list"`
	} `json:"data"`
}