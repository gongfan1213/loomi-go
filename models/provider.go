package models

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// ModelProvider 模型提供商接口
type ModelProvider interface {
	model.BaseChatModel
	
	// 获取提供商名称
	Name() string
	
	// 获取显示名称
	DisplayName() string
	
	// 计算费用
	CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64
	
	// 调用LLM（兼容原有接口）
	CallLLM(ctx context.Context, systemPrompt, userPrompt string, options map[string]interface{}) (string, error)
}

// BaseProvider 基础提供商实现
type BaseProvider struct {
	name        string
	displayName string
	client      model.BaseChatModel
}

// NewBaseProvider 创建基础提供商
func NewBaseProvider(name, displayName string, client model.BaseChatModel) *BaseProvider {
	return &BaseProvider{
		name:        name,
		displayName: displayName,
		client:      client,
	}
}

// Name 获取提供商名称
func (p *BaseProvider) Name() string {
	return p.name
}

// DisplayName 获取显示名称
func (p *BaseProvider) DisplayName() string {
	return p.displayName
}

// CalculateCost 计算费用（基础实现）
func (p *BaseProvider) CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64 {
	// 基础费用计算，子类可以重写
	return 0.0
}

// Generate 实现BaseChatModel接口
func (p *BaseProvider) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	return p.client.Generate(ctx, input, opts...)
}

// Stream 实现BaseChatModel接口
func (p *BaseProvider) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return p.client.Stream(ctx, input, opts...)
}

// CallLLM 调用LLM（兼容原有接口）
func (p *BaseProvider) CallLLM(ctx context.Context, systemPrompt, userPrompt string, options map[string]interface{}) (string, error) {
	// 构建消息
	messages := []*schema.Message{}
	
	if systemPrompt != "" {
		messages = append(messages, schema.SystemMessage(systemPrompt))
	}
	
	messages = append(messages, schema.UserMessage(userPrompt))

	// 调用模型
	response, err := p.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("调用模型失败: %v", err)
	}

	return response.Content, nil
}

// GetInputType 获取输入类型
func (p *BaseProvider) GetInputType() string {
	return "message"
}

// GetOutputType 获取输出类型
func (p *BaseProvider) GetOutputType() string {
	return "message"
}

// GetOptionType 获取选项类型
func (p *BaseProvider) GetOptionType() string {
	return "config"
}

// SetOption 设置选项
func (p *BaseProvider) SetOption(option any) error {
	return nil
}

// GetOption 获取选项
func (p *BaseProvider) GetOption() any {
	return nil
}

// SetCallbacks 设置回调
func (p *BaseProvider) SetCallbacks(callbacks callbacks.Handler) error {
	return nil
}

// GetCallbacks 获取回调
func (p *BaseProvider) GetCallbacks() callbacks.Handler {
	return nil
} 