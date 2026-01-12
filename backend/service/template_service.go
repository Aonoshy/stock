package service

import (
	"fmt"
	"sync"
	"time"
)

// StockTemplate 股票模版结构
type StockTemplate struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Codes       []string  `json:"codes"`
	CreatedAt   time.Time `json:"created_at"`
}

// TemplateService 模版服务
type TemplateService struct {
	templates map[string]*StockTemplate
	mutex     sync.RWMutex
}

// NewTemplateService 创建模版服务
func NewTemplateService() *TemplateService {
	return &TemplateService{
		templates: make(map[string]*StockTemplate),
	}
}

// GetAllTemplates 获取所有模版
func (t *TemplateService) GetAllTemplates() []*StockTemplate {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	templates := make([]*StockTemplate, 0, len(t.templates))
	for _, template := range t.templates {
		templates = append(templates, template)
	}
	return templates
}

// GetTemplate 根据ID获取模版
func (t *TemplateService) GetTemplate(id string) (*StockTemplate, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	template, exists := t.templates[id]
	if !exists {
		return nil, fmt.Errorf("模版 %s 不存在", id)
	}
	return template, nil
}

// CreateTemplate 创建新模版
func (t *TemplateService) CreateTemplate(template *StockTemplate) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if template.ID == "" {
		return fmt.Errorf("模版ID不能为空")
	}

	if template.Name == "" {
		return fmt.Errorf("模版名称不能为空")
	}

	if len(template.Codes) == 0 {
		return fmt.Errorf("股票代码列表不能为空")
	}

	// 检查ID是否已存在
	if _, exists := t.templates[template.ID]; exists {
		return fmt.Errorf("模版ID %s 已存在", template.ID)
	}

	template.CreatedAt = time.Now()
	t.templates[template.ID] = template
	return nil
}

// DeleteTemplate 删除模版
func (t *TemplateService) DeleteTemplate(id string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exists := t.templates[id]; !exists {
		return fmt.Errorf("模版 %s 不存在", id)
	}

	delete(t.templates, id)
	return nil
}

// QueryTemplateStocks 根据模版查询股票价格
func (t *TemplateService) QueryTemplateStocks(id string, stockService *StockService) (map[string]float64, error) {
	template, err := t.GetTemplate(id)
	if err != nil {
		return nil, err
	}

	return stockService.GetBatchStockPrices(template.Codes)
}