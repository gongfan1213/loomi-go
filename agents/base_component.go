package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// BaseComponent 基础组件，提供通用的 BaseChatModel 实现
type BaseComponent struct {
	name string
}

// NewBaseComponent 创建基础组件
func NewBaseComponent(name string) *BaseComponent {
	return &BaseComponent{name: name}
}

// Generate 实现 BaseChatModel 接口
func (c *BaseComponent) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("空输入")
	}
	
	// 默认实现：返回第一个消息的内容
	content := input[0].Content
	return &schema.Message{
		Role:    "assistant",
		Content: fmt.Sprintf("[%s] 处理结果: %s", c.name, content),
	}, nil
}

// Stream 实现 BaseChatModel 接口
func (c *BaseComponent) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	response, err := c.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	
	// 创建一个简单的流适配器
	reader, writer := schema.Pipe[*schema.Message](5)
	
	go func() {
		defer writer.Close()
		writer.Send(response, nil)
	}()
	
	return reader, nil
}

// Collect 实现 Collectable 接口
func (c *BaseComponent) Collect(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.Message, error) {
	var content strings.Builder
	for {
		msg, err := input.Recv()
		if err != nil {
			break
		}
		content.WriteString(msg.Content)
	}
	
	return &schema.Message{
		Role:    "assistant",
		Content: content.String(),
	}, nil
}

// Transform 实现 Transformable 接口
func (c *BaseComponent) Transform(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.StreamReader[*schema.Message], error) {
	// 创建一个简单的流适配器
	reader, writer := schema.Pipe[*schema.Message](5)
	
	go func() {
		defer writer.Close()
		
		for {
			msg, err := input.Recv()
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
func (c *BaseComponent) GetInputType() string {
	return "message"
}

// GetOutputType 获取输出类型
func (c *BaseComponent) GetOutputType() string {
	return "message"
}

// GetOptionType 获取选项类型
func (c *BaseComponent) GetOptionType() string {
	return "config"
}

// SetOption 设置选项
func (c *BaseComponent) SetOption(option any) error {
	return nil
}

// GetOption 获取选项
func (c *BaseComponent) GetOption() any {
	return nil
}

// SetCallbacks 设置回调
func (c *BaseComponent) SetCallbacks(callbacks callbacks.Handler) error {
	return nil
}

// GetCallbacks 获取回调
func (c *BaseComponent) GetCallbacks() callbacks.Handler {
	return nil
} 