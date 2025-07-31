package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiProvider Gemini模型提供商
type GeminiProvider struct {
	*BaseProvider
	client *genai.Client
	model  *genai.GenerativeModel
	config map[string]interface{}
}

// GeminiConfig Gemini配置
const (
	GeminiAPIKey     = "your-gemini-api-key" // 需要替换为实际的API密钥
	GeminiModelName  = "gemini-1.5-pro"
	
	// 费用计算常量
	GeminiCostInputPer1M  = 0.375
	GeminiCostOutputPer1M = 1.875
)

// NewGeminiProvider 创建Gemini提供商
func NewGeminiProvider() (*GeminiProvider, error) {
	// 创建Gemini客户端
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GeminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("创建Gemini客户端失败: %v", err)
	}

	// 创建模型
	model := client.GenerativeModel(GeminiModelName)

	// 创建基础提供商
	baseProvider := NewBaseProvider("gemini-1.5-pro", "Gemini 1.5 Pro", nil)

	provider := &GeminiProvider{
		BaseProvider: baseProvider,
		client:      client,
		model:       model,
		config:      make(map[string]interface{}),
	}

	// 设置基础提供商的客户端为Gemini提供商
	baseProvider.client = provider

	return provider, nil
}

// CalculateCost 计算费用
func (p *GeminiProvider) CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64 {
	inputCost := float64(inputTokens) / 1_000_000 * GeminiCostInputPer1M
	outputCost := float64(outputTokens+thinkingTokens) / 1_000_000 * GeminiCostOutputPer1M
	return inputCost + outputCost
}

// Generate 实现BaseChatModel接口
func (p *GeminiProvider) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// 转换消息格式
	parts := make([]genai.Part, 0, len(input))
	for _, msg := range input {
		parts = append(parts, genai.Text(msg.Content))
	}

	// 调用Gemini API
	resp, err := p.model.GenerateContent(ctx, parts...)
	if err != nil {
		return nil, fmt.Errorf("Gemini API调用失败: %v", err)
	}

	// 处理响应
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("Gemini API返回空响应")
	}

	content := ""
	if len(resp.Candidates[0].Content.Parts) > 0 {
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			content = string(text)
		}
	}

	// 处理文本
	content = p.ProcessText(content)

	// 更新统计信息（Gemini可能没有详细的token信息）
	inputTokens := len(input[0].Content) / 4  // 估算
	outputTokens := len(content) / 4       // 估算
	thinkingTokens := 0

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
func (p *GeminiProvider) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	// 转换消息格式
	parts := make([]genai.Part, 0, len(input))
	for _, msg := range input {
		parts = append(parts, genai.Text(msg.Content))
	}

	// 调用Gemini API（流式）
	iter := p.model.GenerateContentStream(ctx, parts...)

	// 创建一个适配器来转换流
	streamReader := &GeminiStreamReader{iter: iter}
	
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
func (p *GeminiProvider) GetInputType() string {
	return "message"
}

// GetOutputType 获取输出类型
func (p *GeminiProvider) GetOutputType() string {
	return "message"
}

// GetOptionType 获取选项类型
func (p *GeminiProvider) GetOptionType() string {
	return "config"
}

// SetOption 设置选项
func (p *GeminiProvider) SetOption(option any) error {
	if config, ok := option.(map[string]interface{}); ok {
		p.config = config
		return nil
	}
	return fmt.Errorf("invalid option type")
}

// GetOption 获取选项
func (p *GeminiProvider) GetOption() any {
	return p.config
}

// SetCallbacks 设置回调
func (p *GeminiProvider) SetCallbacks(callbacks callbacks.Handler) error {
	return nil
}

// GetCallbacks 获取回调
func (p *GeminiProvider) GetCallbacks() callbacks.Handler {
	return nil
}

// ProcessText 处理文本
func (p *GeminiProvider) ProcessText(text string) string {
	// 去掉中文双引号
	text = strings.ReplaceAll(text, "\"", "")
	text = strings.ReplaceAll(text, "\"", "")
	return text
}

// CallLLM 调用LLM（兼容原有接口）
func (p *GeminiProvider) CallLLM(ctx context.Context, systemPrompt, userPrompt string, options map[string]interface{}) (string, error) {
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

// GeminiStreamReader Gemini流读取器
type GeminiStreamReader struct {
	iter *genai.GenerateContentResponseIterator
}

func (r *GeminiStreamReader) Recv() (*schema.Message, error) {
	resp, err := r.iter.Next()
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("空响应")
	}

	content := ""
	if len(resp.Candidates[0].Content.Parts) > 0 {
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			content = string(text)
		}
	}

	return &schema.Message{
		Role:    "assistant",
		Content: content,
	}, nil
}

func (r *GeminiStreamReader) Close() {
	// Gemini 迭代器不需要显式关闭
} 