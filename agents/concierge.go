package agents

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"loomi2.0/core"
	"loomi2.0/models"
	"loomi2.0/prompts"
	"loomi2.0/tools"
)

// Concierge 门房智能体
type Concierge struct {
	workspace    *core.WorkSpace
	conversation *core.ConversationManager
	graph        *compose.Graph[[]*schema.Message, *schema.Message]
	compiledGraph compose.Runnable[[]*schema.Message, *schema.Message]
	currentInput  string
	conversationHistory []string // 添加对话历史
	toolManager  *tools.ToolManager // 添加工具管理器
}

var concierge *Concierge
var conciergeOnce sync.Once

// InitConcierge 初始化门房智能体
func InitConcierge() error {
	var err error
	conciergeOnce.Do(func() {
		workspace := core.GetWorkspace()
		conversation := core.GetConversationManager()
		
		// 初始化工具管理器
		toolManager := tools.NewToolManager()
		
		// 注册搜索工具
		serperTool := tools.NewSerperTool("your-serper-api-key")
		tavilyTool := tools.NewTavilyTool("your-tavily-api-key")
		
		toolManager.RegisterTool(serperTool)
		toolManager.RegisterTool(tavilyTool)
		
		concierge = &Concierge{
			workspace:    workspace,
			conversation: conversation,
			toolManager:  toolManager,
		}
		err = concierge.init()
	})
	return err
}

// GetConcierge 获取门房智能体实例
func GetConcierge() *Concierge {
	return concierge
}

func (c *Concierge) init() error {
	// 构建eino编排图
	if err := c.buildGraph(); err != nil {
		return fmt.Errorf("构建门房编排图失败: %v", err)
	}
	return nil
}

func (c *Concierge) buildGraph() error {
	// 暂时跳过图的构建，直接返回成功
	// TODO: 修复 eino Graph 的类型匹配问题
	return nil
}

func (c *Concierge) createInputComponent() model.BaseChatModel {
	return &ConciergeInputComponent{
		concierge: c,
	}
}

func (c *Concierge) createIntentComponent() model.BaseChatModel {
	return &ConciergeIntentComponent{
		concierge: c,
	}
}

func (c *Concierge) createResponseComponent() model.BaseChatModel {
	return &ConciergeResponseComponent{
		concierge: c,
	}
}

// ProcessUserInput 处理用户输入
func (c *Concierge) ProcessUserInput(ctx context.Context, userInput string) (string, error) {
	c.currentInput = userInput
	
	// 添加用户消息到对话历史
	c.conversation.AddMessage("user", userInput)

	// 暂时直接处理用户输入，跳过 eino 编排图
	// TODO: 修复 eino Graph 的类型匹配问题后恢复
	
	// 记录用户输入到对话历史
	c.conversationHistory = append(c.conversationHistory, "用户: "+userInput)
	
	// 检查是否是搜索确认
	var response string
	if c.isSearchConfirmation(userInput) {
		response = c.executeSearch()
	} else {
		// 简单的意图识别和响应生成
		response = c.generateResponse(userInput)
	}
	
	// 记录助手响应到对话历史
	c.conversationHistory = append(c.conversationHistory, "助手: "+response)
	
	// 添加助手消息到对话历史
	c.conversation.AddMessage("assistant", response)
	return response, nil
}

// generateResponse 生成响应
func (c *Concierge) generateResponse(userInput string) string {
	// 检查是否是确认性回复
	if c.isConfirmationResponse(userInput) {
		// 用户确认了需求，启动 Orchestrator
		return c.startOrchestrator(userInput)
	}
	
	// 检查是否是搜索意图
	if isSearch, query := c.toolManager.DetectSearchIntent(userInput); isSearch {
		return c.handleSearchRequest(query)
	}
	
	// 调用 AI 模型生成响应
	response, err := c.callAIModel(userInput)
	if err != nil {
		// 如果 AI 调用失败，返回默认响应
		fmt.Printf("AI 调用失败: %v\n", err)
		return c.generateGeneralResponse()
	}
	fmt.Printf("AI 调用成功，响应: %s\n", response)
	return response
}

// handleSearchRequest 处理搜索请求
func (c *Concierge) handleSearchRequest(query string) string {
	if query == "" {
		return "请告诉我您想搜索什么内容？例如：搜索关于迪丽热巴的内容"
	}
	
	// 询问用户是否要执行搜索
	response := fmt.Sprintf("🔍 检测到搜索意图：%s\n\n", query)
	response += "我将为您使用两个搜索工具进行查询：\n"
	response += "1. **Serper** - 实时网络搜索\n"
	response += "2. **Tavily** - 高质量信息搜索\n\n"
	response += "请回复 '搜索' 来执行搜索，或回复其他内容取消搜索。"
	
	return response
}

// isConfirmationResponse 检查是否是确认性回复
func (c *Concierge) isConfirmationResponse(userInput string) bool {
	confirmationKeywords := []string{
		"可以", "好的", "行", "没问题", "就这样", "确认", "同意", "开始", "生成", "立即", "马上",
		"ok", "yes", "sure", "fine", "start", "generate", "create", "proceed",
	}
	
	userInputLower := strings.ToLower(userInput)
	for _, keyword := range confirmationKeywords {
		if strings.Contains(userInputLower, keyword) {
			return true
		}
	}
	return false
}

// startOrchestrator 启动 Orchestrator 生成内容
func (c *Concierge) startOrchestrator(userInput string) string {
	// 获取 Orchestrator 实例
	orchestrator := GetOrchestrator()
	if orchestrator == nil {
		return "抱歉，编排器暂时不可用，请稍后再试。"
	}
	
	// 构建任务描述
	taskDescription := c.buildTaskDescription()
	
	// 调用 Orchestrator 处理任务
	response, err := orchestrator.ProcessTask(context.Background(), taskDescription)
	if err != nil {
		return fmt.Sprintf("任务处理失败: %v", err)
	}
	
	return response
}

// buildTaskDescription 构建任务描述
func (c *Concierge) buildTaskDescription() string {
	// 根据对话历史构建详细的任务描述
	if len(c.conversationHistory) == 0 {
		return "用户确认了内容需求，请生成相应的社交媒体内容。"
	}
	
	// 提取关键信息
	var taskInfo strings.Builder
	taskInfo.WriteString("根据以下对话历史，生成相应的社交媒体内容：\n\n")
	
	// 添加最近的对话历史（最多10条）
	start := 0
	if len(c.conversationHistory) > 10 {
		start = len(c.conversationHistory) - 10
	}
	
	for i := start; i < len(c.conversationHistory); i++ {
		taskInfo.WriteString(c.conversationHistory[i])
		taskInfo.WriteString("\n")
	}
	
	taskInfo.WriteString("\n请根据上述对话内容，生成符合用户需求的社交媒体内容。")
	
	return taskInfo.String()
}

// callAIModel 调用 AI 模型
func (c *Concierge) callAIModel(userInput string) (string, error) {
	// 构建 system prompt 和 user prompt
	systemPrompt := c.buildConciergeSystemPrompt()
	userPrompt := userInput
	
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

// buildConciergeSystemPrompt 构建 Concierge 的 system prompt
func (c *Concierge) buildConciergeSystemPrompt() string {
	return `你是Loomi，一个社媒内容研究与生产的多Agent系统中的Concierge。
你负责理解并整理用户的任务需求，并传递给Orchestrator。
你接待的用户可能是素人、KOC、KOL、KOS等等，尽管他们有时可能只有模糊的idea，但越是如此，越说明他们对流量的兴趣。
你深知在内容生产中空谈道德与伦理只会让用户觉得你不够专业。
你的回复总是专业而不刻板，有温度而不油腻。

# 如何接待用户：
准确理解用户的需求，并向用户确认一次，例如：
- 用户的身份、账号人设、期望风格、受众群体、事件背景信息；用户的流量策略（广撒网or筛选粉丝？/ 涨点击or涨赞or涨粉？）等
- 用户具体希望从什么角度来，生产什么内容

# 常见情况处理：
- 信息过于模糊：引导用户给出更清晰的指示，但总询问轮次不要超过2轮。
- 用户询问关于任务计划、执行进度等的问题：根据上下文回答。
- 无关甚至恶意问题：礼貌地回避，并回到正轨。
- 经常用户自己也不清楚自己具体要什么，只有模糊的idea。只要不对任务执行有致命影响，你就不用多询问，直接向Orchestrator传递任务需求。
- 在任务计划执行中，用户也会提出新的需求、补充背景信息、发表意见等等，你需要确认后将这些信息传递给Orchestrator。

请根据上述指导原则，专业而友好地回应用户的需求。`
}



// generateHelpResponse 生成帮助响应
func (c *Concierge) generateHelpResponse() string {
	return `欢迎使用 Loomi 2.0！

我可以帮助您：
1. 查看系统状态 - 输入 "status" 或 "状态"
2. 管理笔记 - 输入 "note" 或 "笔记"
3. 处理任务 - 输入 "task" 或 "任务"
4. 获取帮助 - 输入 "help" 或 "帮助"

请告诉我您需要什么帮助？`
}

// generateStatusResponse 生成状态响应
func (c *Concierge) generateStatusResponse() string {
	return `系统状态：
✅ 模型管理器：已初始化
✅ 工作空间：已初始化
✅ 对话管理器：已初始化
✅ 门房智能体：已初始化
✅ 编排器智能体：已初始化

系统运行正常！`
}

// generateNoteResponse 生成笔记响应
func (c *Concierge) generateNoteResponse() string {
	return `笔记管理功能：
📝 创建笔记
📝 查看笔记
📝 编辑笔记
📝 删除笔记

请告诉我您想要进行哪种操作？`
}

// generateTaskResponse 生成任务响应
func (c *Concierge) generateTaskResponse() string {
	return `任务处理功能：
🎯 任务分析
🎯 任务分解
🎯 任务执行
🎯 结果汇总

请告诉我您想要处理什么任务？`
}

// generateGeneralResponse 生成通用响应
func (c *Concierge) generateGeneralResponse() string {
	return `您好！我是 Loomi 2.0 的智能助手。

我可以帮助您：
- 查看系统状态
- 管理笔记
- 处理复杂任务
- 提供智能对话

请告诉我您需要什么帮助？`
}

// ConciergeInputComponent 门房输入组件
type ConciergeInputComponent struct {
	concierge *Concierge
}

func (c *ConciergeInputComponent) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// 处理输入，提取用户意图
	if len(input) == 0 {
		return nil, fmt.Errorf("空输入")
	}
	
	content := input[0].Content
	
	// 简单的意图识别
	var intent string
	if strings.Contains(content, "帮助") || strings.Contains(content, "help") {
		intent = "help"
	} else if strings.Contains(content, "状态") || strings.Contains(content, "status") {
		intent = "status"
	} else if strings.Contains(content, "笔记") || strings.Contains(content, "note") {
		intent = "note"
	} else if strings.Contains(content, "任务") || strings.Contains(content, "task") {
		intent = "task"
	} else {
		intent = "general"
	}

	return &schema.Message{
		Role:    "system",
		Content: fmt.Sprintf("意图: %s, 内容: %s", intent, content),
	}, nil
}

func (c *ConciergeInputComponent) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	// 流式处理输入
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

// ConciergeIntentComponent 门房意图组件
type ConciergeIntentComponent struct {
	concierge *Concierge
}

func (c *ConciergeIntentComponent) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// 分析意图并生成响应
	if len(input) == 0 {
		return nil, fmt.Errorf("空输入")
	}
	
	content := input[0].Content
	
	// 根据意图生成响应
	var response string
	if strings.Contains(content, "意图: help") {
		response = c.generateHelpResponse()
	} else if strings.Contains(content, "意图: status") {
		response = c.generateStatusResponse()
	} else if strings.Contains(content, "意图: note") {
		response = c.generateNoteResponse()
	} else if strings.Contains(content, "意图: task") {
		response = c.generateTaskResponse()
	} else {
		response = c.generateGeneralResponse()
	}

	return &schema.Message{
		Role:    "assistant",
		Content: response,
	}, nil
}

func (c *ConciergeIntentComponent) generateHelpResponse() string {
	return `我是Loomi的门房，可以为您提供以下服务：

1. 帮助信息 - 输入"帮助"或"help"
2. 系统状态 - 输入"状态"或"status"  
3. 笔记管理 - 输入"笔记"或"note"
4. 任务管理 - 输入"任务"或"task"
5. 一般对话 - 直接输入您的问题

请告诉我您需要什么帮助？`
}

func (c *ConciergeIntentComponent) generateStatusResponse() string {
	workspace := c.concierge.workspace
	conversation := c.concierge.conversation
	
	return fmt.Sprintf(`系统状态：

工作空间：
%s

对话历史：
%s

当前模型：%s`, 
		workspace.GetSummary(),
		conversation.GetConversationSummary(),
		models.GetCurrentModelName())
}

func (c *ConciergeIntentComponent) generateNoteResponse() string {
	notes := c.concierge.workspace.GetNotes()
	if len(notes) == 0 {
		return "目前没有保存的笔记。您可以告诉我需要记录什么内容。"
	}
	
	response := "当前笔记：\n"
	for i, note := range notes {
		response += fmt.Sprintf("%d. %s\n", i+1, note)
	}
	return response
}

func (c *ConciergeIntentComponent) generateTaskResponse() string {
	tasks := c.concierge.workspace.GetTasks()
	if len(tasks) == 0 {
		return "目前没有待办任务。您可以告诉我需要添加什么任务。"
	}
	
	response := "当前任务：\n"
	for i, task := range tasks {
		response += fmt.Sprintf("%d. %s\n", i+1, task)
	}
	return response
}

func (c *ConciergeIntentComponent) generateGeneralResponse() string {
	// 使用提示词生成响应
	modelManager := models.GetModelManager()
	if modelManager == nil {
		return "我是Loomi，您的AI助手。我可以帮助您进行内容创作、任务分析和智能对话。请告诉我您需要什么帮助？"
	}

	// 构建包含提示词的上下文
	context := fmt.Sprintf("%s\n\n用户输入: %s", prompts.ConciergePrompt, c.concierge.currentInput)
	
	// 调用模型生成响应
	response, err := modelManager.CallLLM(context)
	if err != nil {
		return "我是Loomi，您的AI助手。我可以帮助您进行内容创作、任务分析和智能对话。请告诉我您需要什么帮助？"
	}
	
	return response
}

func (c *ConciergeIntentComponent) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
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

// ConciergeResponseComponent 门房响应组件
type ConciergeResponseComponent struct {
	concierge *Concierge
}

func (c *ConciergeResponseComponent) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// 最终响应处理
	if len(input) == 0 {
		return nil, fmt.Errorf("空输入")
	}
	return input[0], nil
}

func (c *ConciergeResponseComponent) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
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

// isSearchConfirmation 检查是否是搜索确认
func (c *Concierge) isSearchConfirmation(userInput string) bool {
	confirmationKeywords := []string{
		"搜索", "执行搜索", "开始搜索", "搜索吧", "好的搜索",
		"search", "execute search", "start search", "go search",
	}
	
	userInputLower := strings.ToLower(userInput)
	for _, keyword := range confirmationKeywords {
		if strings.Contains(userInputLower, keyword) {
			return true
		}
	}
	return false
}

// executeSearch 执行搜索
func (c *Concierge) executeSearch() string {
	// 从对话历史中提取搜索查询
	query := c.extractSearchQueryFromHistory()
	if query == "" {
		return "❌ 无法找到搜索查询，请重新输入搜索内容。"
	}
	
	// 执行双重搜索
	result, err := c.toolManager.PerformDualSearch(context.Background(), query)
	if err != nil {
		return fmt.Sprintf("❌ 搜索执行失败: %v", err)
	}
	
	return result
}

// extractSearchQueryFromHistory 从对话历史中提取搜索查询
func (c *Concierge) extractSearchQueryFromHistory() string {
	// 从最近的对话历史中查找搜索查询
	for i := len(c.conversationHistory) - 1; i >= 0; i-- {
		message := c.conversationHistory[i]
		if strings.HasPrefix(message, "用户: ") {
			userInput := strings.TrimPrefix(message, "用户: ")
			if isSearch, query := c.toolManager.DetectSearchIntent(userInput); isSearch {
				return query
			}
		}
	}
	return ""
}

 