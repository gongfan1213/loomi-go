package models

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/eino/schema"
)

// SessionStats 会话统计
type SessionStats struct {
	TotalCalls         int     `json:"total_calls"`
	TotalInputTokens   int     `json:"total_input_tokens"`
	TotalOutputTokens  int     `json:"total_output_tokens"`
	TotalThinkingTokens int    `json:"total_thinking_tokens"`
	TotalCost          float64 `json:"total_cost"`
}

// ModelManager 模型管理器
type ModelManager struct {
	providers      map[string]ModelProvider
	currentProvider ModelProvider
	stats          SessionStats
	mu             sync.RWMutex
}

var (
	manager *ModelManager
	once    sync.Once
)

// InitModelManager 初始化模型管理器
func InitModelManager() error {
	var err error
	once.Do(func() {
		manager = &ModelManager{
			providers: make(map[string]ModelProvider),
		}
		err = manager.init()
	})
	return err
}

// GetModelManager 获取模型管理器实例
func GetModelManager() *ModelManager {
	return manager
}

func (m *ModelManager) init() error {
	// 注册默认提供商
	if err := m.registerDefaultProviders(); err != nil {
		return fmt.Errorf("注册默认提供商失败: %v", err)
	}

	return nil
}

func (m *ModelManager) registerDefaultProviders() error {
	// 注册豆包提供商
	doubaoProvider, err := NewDoubaoProvider()
	if err != nil {
		return fmt.Errorf("创建豆包提供商失败: %v", err)
	}
	m.RegisterProvider(doubaoProvider)

	// 注册DeepSeek提供商
	deepseekProvider, err := NewDeepSeekProvider()
	if err != nil {
		return fmt.Errorf("创建DeepSeek提供商失败: %v", err)
	}
	m.RegisterProvider(deepseekProvider)

	// 注册Gemini提供商
	geminiProvider, err := NewGeminiProvider()
	if err != nil {
		return fmt.Errorf("创建Gemini提供商失败: %v", err)
	}
	m.RegisterProvider(geminiProvider)

	return nil
}

// RegisterProvider 注册提供商
func (m *ModelManager) RegisterProvider(provider ModelProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.providers[provider.Name()] = provider
	
	// 如果没有当前提供商，设置为第一个
	if m.currentProvider == nil {
		m.currentProvider = provider
	}
}

// SetCurrentProvider 设置当前提供商
func (m *ModelManager) SetCurrentProvider(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	provider, exists := m.providers[name]
	if !exists {
		return fmt.Errorf("提供商不存在: %s", name)
	}
	
	m.currentProvider = provider
	return nil
}

// GetCurrentProvider 获取当前提供商
func (m *ModelManager) GetCurrentProvider() ModelProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return m.currentProvider
}

// ListProviders 列出所有提供商
func (m *ModelManager) ListProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	names := make([]string, 0, len(m.providers))
	for name := range m.providers {
		names = append(names, name)
	}
	return names
}

// GetProviderDisplayNames 获取提供商显示名称映射
func (m *ModelManager) GetProviderDisplayNames() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	displayNames := make(map[string]string)
	for name, provider := range m.providers {
		displayNames[name] = provider.DisplayName()
	}
	return displayNames
}

// CallCurrentModel 调用当前模型
func (m *ModelManager) CallCurrentModel(ctx context.Context, systemPrompt, userPrompt string, options map[string]interface{}) (string, error) {
	m.mu.RLock()
	provider := m.currentProvider
	m.mu.RUnlock()
	
	if provider == nil {
		return "", fmt.Errorf("没有设置当前模型")
	}

	// 直接调用提供商
	return provider.CallLLM(ctx, systemPrompt, userPrompt, options)
}

// CallLLM 调用LLM（简化版本）
func (m *ModelManager) CallLLM(prompt string) (string, error) {
	ctx := context.Background()
	provider := m.GetCurrentProvider()
	if provider == nil {
		return "", fmt.Errorf("没有设置当前模型")
	}

	// 构建消息
	messages := []*schema.Message{
		schema.UserMessage(prompt),
	}

	// 调用模型
	response, err := provider.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("调用模型失败: %v", err)
	}

	return response.Content, nil
}

// UpdateStats 更新统计信息
func (m *ModelManager) UpdateStats(inputTokens, outputTokens, thinkingTokens int, cost float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.stats.TotalCalls++
	m.stats.TotalInputTokens += inputTokens
	m.stats.TotalOutputTokens += outputTokens
	m.stats.TotalThinkingTokens += thinkingTokens
	m.stats.TotalCost += cost
}

// GetStats 获取统计信息
func (m *ModelManager) GetStats() SessionStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return m.stats
}

// ResetStats 重置统计信息
func (m *ModelManager) ResetStats() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.stats = SessionStats{}
}

// Cleanup 清理资源
func (m *ModelManager) Cleanup() {
	// 清理资源
}

 