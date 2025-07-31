package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/sashabaranov/go-openai"
)

// DoubaoProvider 豆包模型提供商
type DoubaoProvider struct {
	*BaseProvider
	client *openai.Client
	config map[string]interface{}
}

// DoubaoConfig 豆包配置
const (
	DoubaoAPIKey     = "your-doubao-api-key" // 需要替换为实际的API密钥
	DoubaoModelName  = "doubao-pro"
	
	// 费用计算常量
	DoubaoCostInputPer1K  = 0.00012
	DoubaoCostOutputPer1K = 0.00024
)

// NewDoubaoProvider 创建豆包提供商
func NewDoubaoProvider() (*DoubaoProvider, error) {
	// 创建OpenAI客户端（豆包兼容OpenAI API）
	config := openai.DefaultConfig(DoubaoAPIKey)
	config.BaseURL = "https://api.doubao.com/v1" // 豆包 API地址
	
	client := openai.NewClientWithConfig(config)

	// 创建基础提供商
	baseProvider := NewBaseProvider("doubao-pro", "豆包 Pro", nil)

	provider := &DoubaoProvider{
		BaseProvider: baseProvider,
		client:      client,
		config:      make(map[string]interface{}),
	}

	// 设置基础提供商的客户端为豆包提供商
	baseProvider.client = provider

	return provider, nil
}

// CalculateCost 计算费用
func (p *DoubaoProvider) CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64 {
	inputCost := float64(inputTokens) / 1000 * DoubaoCostInputPer1K
	outputCost := float64(outputTokens+thinkingTokens) / 1000 * DoubaoCostOutputPer1K
	return inputCost + outputCost
}

// Generate 实现BaseChatModel接口
func (p *DoubaoProvider) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// 转换消息格式
	messages := make([]openai.ChatCompletionMessage, 0, len(input))
	for _, msg := range input {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	// 调用豆包 API
	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    DoubaoModelName,
		Messages: messages,
	})
	if err != nil {
		return nil, fmt.Errorf("豆包 API调用失败: %v", err)
	}

	// 处理响应
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("豆包 API返回空响应")
	}

	content := resp.Choices[0].Message.Content
	content = p.ProcessText(content)

	// 更新统计信息
	inputTokens := 0
	outputTokens := 0
	thinkingTokens := 0
	
	// 检查Usage字段
	if resp.Usage.PromptTokens > 0 {
		inputTokens = resp.Usage.PromptTokens
		outputTokens = resp.Usage.CompletionTokens
	}

	cost := p.CalculateCost(inputTokens, outputTokens, thinkingTokens)
	
	// 更新全局统计
	modelManager := GetModelManager()
	if modelManager != nil {
		modelManager.UpdateStats(inputTokens, outputTokens, thinkingTokens, cost)
	}

	return &schema.Message{
		Role:    "assistant",
		Content: content,
	}, nil
}

// Stream 实现BaseChatModel接口
func (p *DoubaoProvider) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	// 转换消息格式
	messages := make([]openai.ChatCompletionMessage, 0, len(input))
	for _, msg := range input {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	// 调用豆包 API（流式）
	stream, err := p.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:    DoubaoModelName,
		Messages: messages,
	})
	if err != nil {
		return nil, fmt.Errorf("豆包 API流式调用失败: %v", err)
	}

	// 创建一个适配器来转换流
	streamReader := &DoubaoStreamReader{stream: stream}
	
	// 创建一个简单的流适配器
	reader, writer := schema.Pipe[*schema.Message](5)
	
	go func() {
		defer writer.Close()
		
		for {
			msg, err := streamReader.Recv()
			if err != nil {
				writer.Send(msg, err)
				break
			}
			writer.Send(msg, nil)
		}
	}()
	
	return reader, nil
}

// GetInputType 获取输入类型
func (p *DoubaoProvider) GetInputType() string {
	return "message"
}

// GetOutputType 获取输出类型
func (p *DoubaoProvider) GetOutputType() string {
	return "message"
}

// GetOptionType 获取选项类型
func (p *DoubaoProvider) GetOptionType() string {
	return "config"
}

// SetOption 设置选项
func (p *DoubaoProvider) SetOption(option any) error {
	if config, ok := option.(map[string]interface{}); ok {
		p.config = config
		return nil
	}
	return fmt.Errorf("invalid option type")
}

// GetOption 获取选项
func (p *DoubaoProvider) GetOption() any {
	return p.config
}

// SetCallbacks 设置回调
func (p *DoubaoProvider) SetCallbacks(callbacks callbacks.Handler) error {
	return nil
}

// GetCallbacks 获取回调
func (p *DoubaoProvider) GetCallbacks() callbacks.Handler {
	return nil
}

// ProcessText 处理文本
func (p *DoubaoProvider) ProcessText(text string) string {
	// 去掉中文双引号
	text = strings.ReplaceAll(text, "\"", "")
	text = strings.ReplaceAll(text, "\"", "")
	return text
}

// CallLLM 调用LLM（兼容原有接口）
func (p *DoubaoProvider) CallLLM(ctx context.Context, systemPrompt, userPrompt string, options map[string]interface{}) (string, error) {
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

// DoubaoStreamReader 豆包流读取器
type DoubaoStreamReader struct {
	stream *openai.ChatCompletionStream
}

func (r *DoubaoStreamReader) Recv() (*schema.Message, error) {
	chunk, err := r.stream.Recv()
	if err != nil {
		return nil, err
	}

	if len(chunk.Choices) == 0 {
		return nil, fmt.Errorf("空响应")
	}

	content := chunk.Choices[0].Delta.Content
	return &schema.Message{
		Role:    "assistant",
		Content: content,
	}, nil
}

func (r *DoubaoStreamReader) Close() {
	r.stream.Close()
} 