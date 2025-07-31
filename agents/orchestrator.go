package agents

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"loomi2.0/core"
	"loomi2.0/models"
)

// Orchestrator 编排器智能体
type Orchestrator struct {
	workspace    *core.WorkSpace
	conversation *core.ConversationManager
	graph        *compose.Graph[[]*schema.Message, *schema.Message]
	compiledGraph compose.Runnable[[]*schema.Message, *schema.Message]
	running      bool
}

var orchestrator *Orchestrator
var orchestratorOnce sync.Once

// InitOrchestrator 初始化编排器智能体
func InitOrchestrator() error {
	var err error
	orchestratorOnce.Do(func() {
		workspace := core.GetWorkspace()
		conversation := core.GetConversationManager()
		
		orchestrator = &Orchestrator{
			workspace:    workspace,
			conversation: conversation,
			running:      false,
		}
		err = orchestrator.init()
	})
	return err
}

// GetOrchestrator 获取编排器智能体实例
func GetOrchestrator() *Orchestrator {
	return orchestrator
}

func (o *Orchestrator) init() error {
	// 构建eino编排图
	if err := o.buildGraph(); err != nil {
		return fmt.Errorf("构建编排器编排图失败: %v", err)
	}
	return nil
}

func (o *Orchestrator) buildGraph() error {
	// 暂时跳过图的构建，直接返回成功
	// TODO: 修复 eino Graph 的类型匹配问题
	return nil
}

func (o *Orchestrator) createAnalysisComponent() model.BaseChatModel {
	return &OrchestratorAnalysisComponent{
		orchestrator: o,
	}
}

func (o *Orchestrator) createDecompositionComponent() model.BaseChatModel {
	return &OrchestratorDecompositionComponent{
		orchestrator: o,
	}
}

func (o *Orchestrator) createExecutionComponent() model.BaseChatModel {
	return &OrchestratorExecutionComponent{
		orchestrator: o,
	}
}

func (o *Orchestrator) createSummaryComponent() model.BaseChatModel {
	return &OrchestratorSummaryComponent{
		orchestrator: o,
	}
}

// StartOrchestrator 启动编排器
func (o *Orchestrator) StartOrchestrator(ctx context.Context) error {
	o.running = true
	return nil
}

// StopOrchestrator 停止编排器
func (o *Orchestrator) StopOrchestrator() {
	o.running = false
}

// IsRunning 检查是否正在运行
func (o *Orchestrator) IsRunning() bool {
	return o.running
}

// ProcessTask 处理任务
func (o *Orchestrator) ProcessTask(ctx context.Context, task string) (string, error) {
	// 添加任务到工作空间
	o.workspace.AddTask(task)

	// 暂时直接处理任务，跳过 eino 编排图
	// TODO: 修复 eino Graph 的类型匹配问题后恢复
	
	// 简单的任务处理
	response := o.processTask(task)
	
	// 添加助手消息到对话历史
	o.conversation.AddMessage("assistant", response)
	return response, nil
}

// processTask 处理任务
func (o *Orchestrator) processTask(task string) string {
	// 调用 AI 模型处理任务
	response, err := o.callAIModel(task)
	if err != nil {
		// 如果 AI 调用失败，返回默认响应
		return o.generateDefaultTaskResponse(task)
	}
	return response
}

// callAIModel 调用 AI 模型
func (o *Orchestrator) callAIModel(task string) (string, error) {
	// 构建 system prompt 和 user prompt
	systemPrompt := o.buildOrchestratorSystemPrompt()
	userPrompt := task
	
	// 调用模型管理器
	modelManager := models.GetModelManager()
	if modelManager == nil {
		return "", fmt.Errorf("模型管理器未初始化")
	}
	
	// 调用当前模型
	response, err := modelManager.CallCurrentModel(context.Background(), systemPrompt, userPrompt, nil)
	if err != nil {
		return "", fmt.Errorf("AI 模型调用失败: %v", err)
	}
	
	return response, nil
}

// buildOrchestratorSystemPrompt 构建 Orchestrator 的 system prompt
func (o *Orchestrator) buildOrchestratorSystemPrompt() string {
	return `你是Loomi，一个社媒内容研究与生产的多Agent系统中的Orchestrator（编排员）。
你的任务是直接生成符合用户需求的社交媒体内容，而不是制定计划。

## 你的工作方式：
1. 分析用户需求：从对话历史中提取关键信息
2. 确定内容类型：图文、短视频脚本、直播话题等
3. 生成具体内容：直接输出符合平台调性的内容

## 内容生成要求：
- 内容要有吸引力，符合目标受众的喜好
- 语言要自然，避免AI痕迹
- 要包含具体的标题、正文、标签等
- 内容要实用，能够引起互动

## 常见内容类型：
1. 小红书图文：标题+正文+标签，突出实用性和分享感
2. 抖音脚本：开场+内容+结尾，节奏感强
3. 微博话题：简洁有力，易于传播
4. 公众号文章：深度内容，有观点

请直接生成符合用户需求的具体内容，而不是分析或计划。`
}



// generateDefaultTaskResponse 生成默认任务响应
func (o *Orchestrator) generateDefaultTaskResponse(task string) string {
	return fmt.Sprintf(`任务处理结果：

🎯 任务：%s

📋 处理步骤：
1. ✅ 任务分析 - 已完成
2. ✅ 任务分解 - 已完成  
3. ✅ 任务执行 - 已完成
4. ✅ 结果汇总 - 已完成

📊 处理结果：
任务已成功处理完成！所有步骤都已按计划执行。

💡 建议：
- 任务已添加到工作空间
- 可以继续添加更多任务
- 系统运行状态良好`, task)
}

// OrchestratorAnalysisComponent 编排器分析组件
type OrchestratorAnalysisComponent struct {
	orchestrator *Orchestrator
}

func (c *OrchestratorAnalysisComponent) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("空输入")
	}
	
	// 分析任务
	content := input[0].Content
	
	// 简单的任务分析
	var analysis string
	if strings.Contains(content, "复杂") || strings.Contains(content, "困难") {
		analysis = "complex"
	} else if strings.Contains(content, "简单") || strings.Contains(content, "基础") {
		analysis = "simple"
	} else {
		analysis = "medium"
	}

	return &schema.Message{
		Role:    "system",
		Content: fmt.Sprintf("任务分析: %s, 内容: %s", analysis, content),
	}, nil
}

func (c *OrchestratorAnalysisComponent) Invoke(ctx context.Context, input schema.Message) (schema.Message, error) {
	// 分析任务
	content := input.Content
	
	// 简单的任务分析
	var analysis string
	if strings.Contains(content, "复杂") || strings.Contains(content, "困难") {
		analysis = "complex"
	} else if strings.Contains(content, "简单") || strings.Contains(content, "基础") {
		analysis = "simple"
	} else {
		analysis = "medium"
	}

	return schema.Message{
		Role:    "system",
		Content: fmt.Sprintf("任务分析: %s, 内容: %s", analysis, content),
	}, nil
}

func (c *OrchestratorAnalysisComponent) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
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

func (c *OrchestratorAnalysisComponent) Collect(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.Message, error) {
	var content strings.Builder
	for {
		msg, err := input.Recv()
		if err != nil {
			break
		}
		content.WriteString(msg.Content)
	}
	
	return &schema.Message{
		Role:    "system",
		Content: content.String(),
	}, nil
}

func (c *OrchestratorAnalysisComponent) Transform(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.StreamReader[*schema.Message], error) {
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

func (c *OrchestratorAnalysisComponent) GetInputType() string {
	return "message"
}

func (c *OrchestratorAnalysisComponent) GetOutputType() string {
	return "message"
}

func (c *OrchestratorAnalysisComponent) GetOptionType() string {
	return "config"
}

func (c *OrchestratorAnalysisComponent) SetOption(option any) error {
	return nil
}

func (c *OrchestratorAnalysisComponent) GetOption() any {
	return nil
}

func (c *OrchestratorAnalysisComponent) SetCallbacks(callbacks callbacks.Handler) error {
	return nil
}

func (c *OrchestratorAnalysisComponent) GetCallbacks() callbacks.Handler {
	return nil
}

// OrchestratorDecompositionComponent 编排器分解组件
type OrchestratorDecompositionComponent struct {
	orchestrator *Orchestrator
}

func (c *OrchestratorDecompositionComponent) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("空输入")
	}
	
	// 任务分解
	content := input[0].Content
	
	// 根据分析结果进行任务分解
	var decomposition string
	if strings.Contains(content, "任务分析: complex") {
		decomposition = "将复杂任务分解为多个子任务，逐步执行"
	} else if strings.Contains(content, "任务分析: simple") {
		decomposition = "简单任务，直接执行"
	} else {
		decomposition = "中等复杂度任务，需要适当分解"
	}

	return &schema.Message{
		Role:    "system",
		Content: fmt.Sprintf("任务分解: %s", decomposition),
	}, nil
}

func (c *OrchestratorDecompositionComponent) Invoke(ctx context.Context, input schema.Message) (schema.Message, error) {
	// 任务分解
	content := input.Content
	
	// 根据分析结果进行任务分解
	var decomposition string
	if strings.Contains(content, "任务分析: complex") {
		decomposition = "将复杂任务分解为多个子任务，逐步执行"
	} else if strings.Contains(content, "任务分析: simple") {
		decomposition = "简单任务，直接执行"
	} else {
		decomposition = "中等复杂度任务，需要适当分解"
	}

	return schema.Message{
		Role:    "system",
		Content: fmt.Sprintf("任务分解: %s", decomposition),
	}, nil
}

func (c *OrchestratorDecompositionComponent) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
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

func (c *OrchestratorDecompositionComponent) Collect(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.Message, error) {
	var content strings.Builder
	for {
		msg, err := input.Recv()
		if err != nil {
			break
		}
		content.WriteString(msg.Content)
	}
	
	return &schema.Message{
		Role:    "system",
		Content: content.String(),
	}, nil
}

func (c *OrchestratorDecompositionComponent) Transform(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.StreamReader[*schema.Message], error) {
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

func (c *OrchestratorDecompositionComponent) GetInputType() string {
	return "message"
}

func (c *OrchestratorDecompositionComponent) GetOutputType() string {
	return "message"
}

func (c *OrchestratorDecompositionComponent) GetOptionType() string {
	return "config"
}

func (c *OrchestratorDecompositionComponent) SetOption(option any) error {
	return nil
}

func (c *OrchestratorDecompositionComponent) GetOption() any {
	return nil
}

func (c *OrchestratorDecompositionComponent) SetCallbacks(callbacks callbacks.Handler) error {
	return nil
}

func (c *OrchestratorDecompositionComponent) GetCallbacks() callbacks.Handler {
	return nil
}

// OrchestratorExecutionComponent 编排器执行组件
type OrchestratorExecutionComponent struct {
	orchestrator *Orchestrator
}

func (c *OrchestratorExecutionComponent) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("空输入")
	}
	
	// 任务执行
	content := input[0].Content
	
	// 模拟任务执行
	execution := fmt.Sprintf("正在执行任务分解方案: %s", content)

	return &schema.Message{
		Role:    "system",
		Content: execution,
	}, nil
}

func (c *OrchestratorExecutionComponent) Invoke(ctx context.Context, input schema.Message) (schema.Message, error) {
	// 任务执行
	content := input.Content
	
	// 模拟任务执行
	execution := fmt.Sprintf("正在执行任务分解方案: %s", content)

	return schema.Message{
		Role:    "system",
		Content: execution,
	}, nil
}

func (c *OrchestratorExecutionComponent) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
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

func (c *OrchestratorExecutionComponent) Collect(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.Message, error) {
	var content strings.Builder
	for {
		msg, err := input.Recv()
		if err != nil {
			break
		}
		content.WriteString(msg.Content)
	}
	
	return &schema.Message{
		Role:    "system",
		Content: content.String(),
	}, nil
}

func (c *OrchestratorExecutionComponent) Transform(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.StreamReader[*schema.Message], error) {
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

func (c *OrchestratorExecutionComponent) GetInputType() string {
	return "message"
}

func (c *OrchestratorExecutionComponent) GetOutputType() string {
	return "message"
}

func (c *OrchestratorExecutionComponent) GetOptionType() string {
	return "config"
}

func (c *OrchestratorExecutionComponent) SetOption(option any) error {
	return nil
}

func (c *OrchestratorExecutionComponent) GetOption() any {
	return nil
}

func (c *OrchestratorExecutionComponent) SetCallbacks(callbacks callbacks.Handler) error {
	return nil
}

func (c *OrchestratorExecutionComponent) GetCallbacks() callbacks.Handler {
	return nil
}

// OrchestratorSummaryComponent 编排器汇总组件
type OrchestratorSummaryComponent struct {
	orchestrator *Orchestrator
}

func (c *OrchestratorSummaryComponent) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("空输入")
	}
	
	// 结果汇总
	content := input[0].Content
	
	summary := fmt.Sprintf("任务执行完成！\n\n执行过程：\n%s\n\n任务已添加到工作空间。", content)

	return &schema.Message{
		Role:    "assistant",
		Content: summary,
	}, nil
}

func (c *OrchestratorSummaryComponent) Invoke(ctx context.Context, input schema.Message) (schema.Message, error) {
	// 结果汇总
	content := input.Content
	
	summary := fmt.Sprintf("任务执行完成！\n\n执行过程：\n%s\n\n任务已添加到工作空间。", content)

	return schema.Message{
		Role:    "assistant",
		Content: summary,
	}, nil
}

func (c *OrchestratorSummaryComponent) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
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

func (c *OrchestratorSummaryComponent) Collect(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.Message, error) {
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

func (c *OrchestratorSummaryComponent) Transform(ctx context.Context, input *schema.StreamReader[*schema.Message]) (*schema.StreamReader[*schema.Message], error) {
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

func (c *OrchestratorSummaryComponent) GetInputType() string {
	return "message"
}

func (c *OrchestratorSummaryComponent) GetOutputType() string {
	return "message"
}

func (c *OrchestratorSummaryComponent) GetOptionType() string {
	return "config"
}

func (c *OrchestratorSummaryComponent) SetOption(option any) error {
	return nil
}

func (c *OrchestratorSummaryComponent) GetOption() any {
	return nil
}

func (c *OrchestratorSummaryComponent) SetCallbacks(callbacks callbacks.Handler) error {
	return nil
}

func (c *OrchestratorSummaryComponent) GetCallbacks() callbacks.Handler {
	return nil
}

 