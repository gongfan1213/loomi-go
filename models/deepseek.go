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

// DeepSeekProvider DeepSeek模型提供商
type DeepSeekProvider struct {
	*BaseProvider
	client *openai.Client
	config map[string]interface{}
}

// DeepSeekConfig DeepSeek配置
const (
	DeepSeekAPIKey     = "your-deepseek-api-key" // 需要替换为实际的API密钥
	DeepSeekModelName  = "deepseek-chat"
	
	// 费用计算常量
	DeepSeekCostInputPer1K  = 0.00014
	DeepSeekCostOutputPer1K = 0.00028
)

// NewDeepSeekProvider 创建DeepSeek提供商
func NewDeepSeekProvider() (*DeepSeekProvider, error) {
	// 创建OpenAI客户端（DeepSeek兼容OpenAI API）
	config := openai.DefaultConfig(DeepSeekAPIKey)
	config.BaseURL = "https://api.deepseek.com/v1" // DeepSeek API地址
	
	client := openai.NewClientWithConfig(config)

	// 创建基础提供商
	baseProvider := NewBaseProvider("deepseek-chat", "DeepSeek Chat", nil)

	provider := &DeepSeekProvider{
		BaseProvider: baseProvider,
		client:      client,
		config:      make(map[string]interface{}),
	}

	// 设置基础提供商的客户端为DeepSeek提供商
	baseProvider.client = provider

	return provider, nil
}

// CalculateCost 计算费用
func (p *DeepSeekProvider) CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64 {
	inputCost := float64(inputTokens) / 1000 * DeepSeekCostInputPer1K
	outputCost := float64(outputTokens+thinkingTokens) / 1000 * DeepSeekCostOutputPer1K
	return inputCost + outputCost
}

// Generate 实现BaseChatModel接口
func (p *DeepSeekProvider) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// 转换消息格式
	messages := make([]openai.ChatCompletionMessage, 0, len(input))
	for _, msg := range input {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	// 调用DeepSeek API
	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    DeepSeekModelName,
		Messages: messages,
	})
	if err != nil {
		return nil, fmt.Errorf("DeepSeek API调用失败: %v", err)
	}

	// 处理响应
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("DeepSeek API返回空响应")
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
func (p *DeepSeekProvider) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	// 转换消息格式
	messages := make([]openai.ChatCompletionMessage, 0, len(input))
	for _, msg := range input {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	// 调用DeepSeek API（流式）
	stream, err := p.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:    DeepSeekModelName,
		Messages: messages,
	})
	if err != nil {
		return nil, fmt.Errorf("DeepSeek API流式调用失败: %v", err)
	}

	// 创建一个适配器来转换流
	streamReader := &DeepSeekStreamReader{stream: stream}
	
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
func (p *DeepSeekProvider) GetInputType() string {
	return "message"
}

// GetOutputType 获取输出类型
func (p *DeepSeekProvider) GetOutputType() string {
	return "message"
}

// GetOptionType 获取选项类型
func (p *DeepSeekProvider) GetOptionType() string {
	return "config"
}

// SetOption 设置选项
func (p *DeepSeekProvider) SetOption(option any) error {
	if config, ok := option.(map[string]interface{}); ok {
		p.config = config
		return nil
	}
	return fmt.Errorf("invalid option type")
}

// GetOption 获取选项
func (p *DeepSeekProvider) GetOption() any {
	return p.config
}

// SetCallbacks 设置回调
func (p *DeepSeekProvider) SetCallbacks(callbacks callbacks.Handler) error {
	return nil
}

// GetCallbacks 获取回调
func (p *DeepSeekProvider) GetCallbacks() callbacks.Handler {
	return nil
}

// ProcessText 处理文本
func (p *DeepSeekProvider) ProcessText(text string) string {
	// 去掉中文双引号
	text = strings.ReplaceAll(text, "\"", "")
	text = strings.ReplaceAll(text, "\"", "")
	return text
}

// CallLLM 调用LLM（兼容原有接口）
func (p *DeepSeekProvider) CallLLM(ctx context.Context, systemPrompt, userPrompt string, options map[string]interface{}) (string, error) {
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

// DeepSeekStreamReader DeepSeek流读取器
type DeepSeekStreamReader struct {
	stream *openai.ChatCompletionStream
}

func (r *DeepSeekStreamReader) Recv() (*schema.Message, error) {
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

func (r *DeepSeekStreamReader) Close() {
	r.stream.Close()
} 