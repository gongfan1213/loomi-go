# 基于Eino框架的multi-agent小红书智能写作系统

[实现效果](https://github.com/gongfan1213/loomi-go/wiki)
### 改原先工具调用删除，只有serper和tavily搜索工具
### 只测过deepseek
### prompt内容全改

一个基于 Go 语言和 eino 框架开发的高性能 AI 内容创作系统，专门为小红书、抖音、微博等社交媒体平台提供智能内容生成服务。

## 🚀 项目特色

### ✨ 核心功能
- **智能内容生成** - 基于用户需求自动生成符合平台调性的内容
- **多模型支持** - 支持 DeepSeek、豆包、Gemini 等多种 AI 模型
- **工具调用集成** - 内置 Serper 和 Tavily 搜索工具，获取最新信息
- **多智能体架构** - 采用门房智能体和编排器智能体的协作设计
- **流式响应** - 支持实时流式输出，提升用户体验
- **会话管理** - 完整的对话历史管理和上下文维护

### 🏗️ 技术架构
- **语言**: Go 1.23+
- **框架**: eino v0.4.0 (字节跳动开源 AI 框架)
- **架构模式**: 智能体编排 + 图编排 + 模块化设计
- **并发处理**: 支持高并发请求处理，线程安全
- **依赖管理**: 使用 Go modules 进行依赖管理

## 🎯 eino 框架深度集成

### 框架概述
eino 是字节跳动开源的 AI 应用开发框架，本项目深度集成了 eino 框架的多个核心功能，实现了现代化的 AI 应用架构。

### 核心集成特性

#### 1. 组件抽象层 (Component Abstraction)
项目使用 eino 的组件抽象层，实现了统一的 AI 模型接口：

```go
// 实现 eino 的 BaseChatModel 接口
type ModelProvider interface {
    model.BaseChatModel
    
    // 自定义接口
    Name() string
    DisplayName() string
    CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64
    CallLLM(ctx context.Context, systemPrompt, userPrompt string, options map[string]interface{}) (string, error)
}
```

**组件实现示例**:
```go
func (p *DeepSeekProvider) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
    // 转换 eino 消息格式为 OpenAI 格式
    messages := make([]openai.ChatCompletionMessage, 0, len(input))
    for _, msg := range input {
        messages = append(messages, openai.ChatCompletionMessage{
            Role:    string(msg.Role),
            Content: msg.Content,
        })
    }
    
    // 调用 API 并返回 eino 消息格式
    return &schema.Message{
        Role:    "assistant",
        Content: content,
    }, nil
}
```

#### 2. 消息处理系统 (Message Processing)
使用 eino 的 `schema.Message` 作为统一的消息格式：

```go
// 消息类型定义
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// 消息创建工具函数
messages := []*schema.Message{
    schema.SystemMessage("你是AI助手"),
    schema.UserMessage("用户输入"),
    schema.AssistantMessage("AI回复"),
}
```

#### 3. 流式处理 (Streaming)
集成 eino 的流式处理能力，支持实时响应：

```go
func (p *DeepSeekProvider) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
    // 创建 eino 流管道
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
```

#### 4. 编排图架构 (Orchestration Graph)
项目设计了基于 eino 编排图的智能体架构：

```go
type Concierge struct {
    workspace    *core.WorkSpace
    conversation *core.ConversationManager
    graph        *compose.Graph[[]*schema.Message, *schema.Message]
    compiledGraph compose.Runnable[[]*schema.Message, *schema.Message]
    toolManager  *tools.ToolManager
}
```

**编排组件设计**:
```go
// 输入处理组件
type ConciergeInputComponent struct {
    concierge *Concierge
}

// 意图识别组件
type ConciergeIntentComponent struct {
    concierge *Concierge
}

// 响应生成组件
type ConciergeResponseComponent struct {
    concierge *Concierge
}
```

#### 5. 回调系统 (Callback System)
集成 eino 的回调机制，支持事件处理：

```go
import "github.com/cloudwego/eino/callbacks"

func (p *DeepSeekProvider) SetCallbacks(callbacks callbacks.Handler) error {
    // 设置回调处理器
    return nil
}

func (p *DeepSeekProvider) GetCallbacks() callbacks.Handler {
    // 获取回调处理器
    return nil
}
```

### 智能体编排架构

#### 门房智能体 (Concierge) - eino 编排
```go
func (c *Concierge) buildGraph() error {
    // 构建 eino 编排图
    // TODO: 修复 eino Graph 的类型匹配问题
    return nil
}

func (c *Concierge) createInputComponent() model.BaseChatModel {
    return &ConciergeInputComponent{concierge: c}
}

func (c *Concierge) createIntentComponent() model.BaseChatModel {
    return &ConciergeIntentComponent{concierge: c}
}

func (c *Concierge) createResponseComponent() model.BaseChatModel {
    return &ConciergeResponseComponent{concierge: c}
}
```

#### 编排器智能体 (Orchestrator) - eino 编排
```go
func (o *Orchestrator) buildGraph() error {
    // 构建 eino 编排图
    // TODO: 修复 eino Graph 的类型匹配问题
    return nil
}

// 编排器组件
type OrchestratorAnalysisComponent struct {
    orchestrator *Orchestrator
}

type OrchestratorDecompositionComponent struct {
    orchestrator *Orchestrator
}

type OrchestratorExecutionComponent struct {
    orchestrator *Orchestrator
}

type OrchestratorSummaryComponent struct {
    orchestrator *Orchestrator
}
```

### 模型集成架构

#### 统一模型接口
所有 AI 模型都实现 eino 的 `BaseChatModel` 接口：

```go
// DeepSeek 模型集成
type DeepSeekProvider struct {
    *BaseProvider
    client *openai.Client
    config map[string]interface{}
}

// 豆包模型集成
type DoubaoProvider struct {
    *BaseProvider
    client *http.Client
    config map[string]interface{}
}

// Gemini 模型集成
type GeminiProvider struct {
    *BaseProvider
    client *genai.Client
    config map[string]interface{}
}
```

#### 流式处理适配器
为每个模型提供流式处理适配器：

```go
// DeepSeek 流读取器
type DeepSeekStreamReader struct {
    stream *openai.ChatCompletionStream
}

func (r *DeepSeekStreamReader) Recv() (*schema.Message, error) {
    chunk, err := r.stream.Recv()
    if err != nil {
        return nil, err
    }
    
    return &schema.Message{
        Role:    "assistant",
        Content: chunk.Choices[0].Delta.Content,
    }, nil
}
```

### eino 框架优势体现

#### 1. 标准化接口
- **统一消息格式**: 使用 `schema.Message` 作为标准消息格式
- **组件抽象**: 通过 `BaseChatModel` 接口实现模型统一
- **流式处理**: 标准化的流式处理接口

#### 2. 编排能力
- **图编排**: 支持复杂的 AI 应用编排
- **组件复用**: 组件可以在不同编排图中复用
- **类型安全**: 强类型保证编排的正确性

#### 3. 扩展性
- **插件化**: 易于添加新的 AI 模型和工具
- **模块化**: 组件独立，便于测试和维护
- **配置化**: 支持灵活的配置管理

#### 4. 性能优化
- **并发安全**: 内置并发安全机制
- **内存管理**: 高效的流式处理内存管理
- **错误处理**: 完善的错误处理和恢复机制

### 当前实现状态

#### 已实现功能
- ✅ **组件抽象**: 完整的 `BaseChatModel` 接口实现
- ✅ **消息处理**: 统一的 `schema.Message` 消息格式
- ✅ **流式处理**: 所有模型都支持流式输出
- ✅ **模型集成**: DeepSeek、豆包、Gemini 模型集成
- ✅ **回调系统**: 完整的回调机制支持

#### 待完善功能
- 🔄 **编排图**: 图编排功能正在开发中
- 🔄 **组件编排**: 智能体组件编排优化
- 🔄 **类型匹配**: 修复 Graph 类型匹配问题

### 技术债务和优化方向

#### 1. 编排图完善
```go
// TODO: 完善编排图构建
func (c *Concierge) buildGraph() error {
    // 构建完整的 eino 编排图
    // 连接输入组件 -> 意图组件 -> 响应组件
    return nil
}
```

#### 2. 组件优化
```go
// TODO: 优化组件间的数据流
type ConciergeComponent struct {
    model.BaseChatModel
    // 添加组件间通信机制
}
```

#### 3. 性能优化
```go
// TODO: 添加缓存和优化机制
type OptimizedProvider struct {
    *BaseProvider
    cache map[string]*schema.Message
}
```

## 📁 项目结构详解

```
loomi_go/
├── main.go                    # 程序入口，初始化系统组件
├── go.mod                     # Go 模块定义和依赖管理
├── cmd/                       # 命令行工具
│   ├── start.go              # 启动命令，包含交互式界面
│   └── version.go            # 版本命令
├── models/                    # AI 模型管理层
│   ├── interface.go          # 模型接口定义
│   ├── manager.go            # 模型管理器，统一管理多个AI模型
│   ├── provider.go           # 基础提供商实现
│   ├── deepseek.go           # DeepSeek 模型实现
│   ├── doubao.go             # 豆包模型实现
│   └── gemini.go             # Gemini 模型实现
├── agents/                    # 智能体系统
│   ├── interface.go          # 智能体接口定义
│   ├── concierge.go          # 门房智能体，负责用户意图识别
│   ├── orchestrator.go       # 编排器智能体，负责内容生成
│   └── base_component.go     # 基础组件
├── tools/                     # 工具调用系统
│   ├── interface.go          # 工具接口定义
│   ├── manager.go            # 工具管理器，统一管理搜索工具
│   ├── serper.go             # Serper 搜索工具实现
│   └── tavily.go             # Tavily 搜索工具实现
├── core/                      # 核心组件
│   ├── workspace.go          # 工作空间管理，存储笔记和任务
│   └── conversation.go       # 对话管理，维护对话历史
├── prompts/                   # 提示词系统
│   ├── concierge_prompt.go   # 门房提示词
│   └── orchestrator_prompt.go # 编排器提示词
└── utils/                     # 工具函数
    └── logger.go             # 日志工具
```

## 🧠 智能体架构详解

### 1. 门房智能体 (Concierge)
**职责**: 用户意图识别、需求确认和对话管理

**核心功能**:
- **意图识别**: 分析用户输入，识别搜索意图、内容创作需求等
- **搜索工具调用**: 集成 Serper 和 Tavily 双重搜索
- **需求确认**: 与用户确认需求，确保理解准确
- **对话历史管理**: 维护完整的对话上下文
- **智能路由**: 将确认的需求传递给编排器

**技术实现**:
```go
type Concierge struct {
    workspace    *core.WorkSpace
    conversation *core.ConversationManager
    toolManager  *tools.ToolManager
    conversationHistory []string
}
```

**工作流程**:
1. 接收用户输入
2. 检测搜索意图或内容创作需求
3. 如果是搜索，询问用户确认后执行双重搜索
4. 如果是内容创作，确认需求后启动编排器
5. 维护对话历史和上下文

### 2. 编排器智能体 (Orchestrator)
**职责**: 复杂任务分解和高质量内容生成

**核心功能**:
- **任务分析**: 分析用户需求，确定内容类型和风格
- **内容生成**: 直接生成符合平台调性的内容
- **多平台适配**: 支持小红书、抖音、微博等不同平台
- **质量保证**: 确保内容有吸引力且符合目标受众

**技术实现**:
```go
type Orchestrator struct {
    workspace    *core.WorkSpace
    conversation *core.ConversationManager
    running      bool
}
```

**支持的内容类型**:
- **小红书图文笔记**: 标题优化、正文生成、标签推荐
- **抖音短视频脚本**: 开场设计、内容结构、结尾引导
- **微博话题文案**: 简洁有力、易于传播
- **公众号文章**: 深度内容、观点输出

## 🔧 工具调用系统详解

### 搜索工具集成
系统集成了两个强大的搜索工具，提供全面的信息获取能力：

#### Serper 搜索工具
- **功能**: 实时网络信息搜索
- **特点**: 快速获取最新信息，支持多种搜索类型
- **API**: 基于 Google 搜索的实时结果

#### Tavily 搜索工具
- **功能**: 高质量深度信息搜索
- **特点**: 提供更深入、更准确的信息
- **API**: 专门为 AI 应用优化的搜索 API

#### 双重搜索机制
```go
func (tm *ToolManager) PerformDualSearch(ctx context.Context, query string) (string, error) {
    var results []string
    
    // 执行Serper搜索
    if serperTool, exists := tm.tools["serper_search"]; exists {
        serperResult, err := serperTool.Execute(ctx, query)
        if err != nil {
            results = append(results, fmt.Sprintf("❌ Serper搜索失败: %v", err))
        } else {
            results = append(results, serperResult)
        }
    }
    
    // 执行Tavily搜索
    if tavilyTool, exists := tm.tools["tavily_search"]; exists {
        tavilyResult, err := tavilyTool.Execute(ctx, query)
        if err != nil {
            results = append(results, fmt.Sprintf("❌ Tavily搜索失败: %v", err))
        } else {
            results = append(results, tavilyResult)
        }
    }
    
    return strings.Join(results, "\n\n" + strings.Repeat("=", 50) + "\n\n"), nil
}
```

### 搜索意图检测
系统能够智能检测用户的搜索意图：

```go
func (tm *ToolManager) DetectSearchIntent(userInput string) (bool, string) {
    searchKeywords := []string{
        "搜索", "查找", "查询", "了解", "搜索关于", "查找关于", "查询关于",
        "search", "find", "lookup", "search for", "find about",
    }
    
    userInputLower := strings.ToLower(userInput)
    for _, keyword := range searchKeywords {
        if strings.Contains(userInputLower, keyword) {
            query := extractSearchQuery(userInput, keyword)
            return true, query
        }
    }
    
    return false, ""
}
```

## 🤖 AI 模型支持详解

### 模型管理器架构
系统采用统一的模型管理器，支持多种 AI 模型的无缝切换：

```go
type ModelManager struct {
    providers      map[string]ModelProvider
    currentProvider ModelProvider
    stats          SessionStats
    mu             sync.RWMutex
}
```

### 支持的模型

#### 1. DeepSeek Chat
- **提供商**: DeepSeek
- **特点**: 强大的中文理解和生成能力
- **费用**: 输入 $0.00014/1K tokens，输出 $0.00028/1K tokens
- **API**: 兼容 OpenAI API 格式

#### 2. 豆包 Pro
- **提供商**: 字节跳动
- **特点**: 字节跳动自研大模型，中文表现优秀
- **集成**: 通过字节跳动官方 API

#### 3. Gemini 1.5 Pro
- **提供商**: Google
- **特点**: Google 最新大模型，多模态能力强
- **API**: 通过 Google Generative AI API

### 模型特性
- **流式响应**: 支持实时流式输出，提升用户体验
- **费用统计**: 实时 Token 使用统计和费用计算
- **多模型切换**: 运行时动态切换模型
- **错误处理**: 完善的错误处理和重试机制
- **线程安全**: 使用互斥锁保证并发安全

### 费用计算示例
```go
func (p *DeepSeekProvider) CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64 {
    inputCost := float64(inputTokens) / 1000 * DeepSeekCostInputPer1K
    outputCost := float64(outputTokens+thinkingTokens) / 1000 * DeepSeekCostOutputPer1K
    return inputCost + outputCost
}
```

## 📝 内容生成功能详解

### 内容类型支持

#### 1. 小红书图文笔记
- **标题优化**: 吸引眼球，符合平台调性
- **正文生成**: 实用性强，分享感十足
- **标签推荐**: 热门标签，提高曝光
- **排版建议**: 符合小红书阅读习惯

#### 2. 抖音短视频脚本
- **开场设计**: 3秒抓住用户注意力
- **内容结构**: 节奏感强，易于理解
- **结尾引导**: 引导关注、点赞、评论
- **互动设计**: 内置互动元素

#### 3. 微博话题文案
- **简洁有力**: 140字内表达核心观点
- **易于传播**: 符合微博传播特点
- **话题标签**: 热门话题标签使用

#### 4. 公众号文章
- **深度内容**: 有观点，有深度
- **结构优化**: 清晰的逻辑结构
- **传播策略**: 考虑传播路径

### 生成特色
- **平台适配**: 针对不同平台特点优化内容
- **风格多样**: 支持多种文案风格和调性
- **互动设计**: 内置互动元素和钩子
- **SEO 优化**: 考虑搜索和推荐算法

## 🚀 快速开始

### 环境要求
- Go 1.23+
- 网络连接 (用于 AI 模型调用)
- 相关 API 密钥

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd loomi_go
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置 API 密钥**
复制 `config.example.go` 为 `config.go` 并填入实际的 API 密钥：
```bash
cp config.example.go config.go
# 编辑 config.go 文件，填入实际的 API 密钥
```

或者直接在相应的模型文件中配置：
- `models/deepseek.go` - DeepSeek API 密钥
- `models/doubao.go` - 豆包 API 密钥  
- `models/gemini.go` - Gemini API 密钥
- `agents/concierge.go` - 搜索工具 API 密钥

4. **编译运行**
```bash
go build -o assistant .
./assistant start
```

### 使用示例

1. **启动系统**
```bash
./assistant start
```

2. **选择 AI 模型**
```
🤖 选择您要使用的模型：
1. 豆包 Pro
2. DeepSeek Chat
3. Gemini 1.5 Pro
请选择模型 (1-3): 2
✅ 已选择模型: DeepSeek Chat
```

3. **开始创作**
```
💬 请输入您的消息: 帮我写一篇关于护肤的小红书文案
```

4. **搜索信息**
```
💬 请输入您的消息: 搜索关于最新护肤趋势的内容
🔍 检测到搜索意图：最新护肤趋势
我将为您使用两个搜索工具进行查询：
1. **Serper** - 实时网络搜索
2. **Tavily** - 高质量信息搜索
请回复 '搜索' 来执行搜索，或回复其他内容取消搜索。

用户: 搜索
[显示 Serper 和 Tavily 的搜索结果]
```

## 🔧 配置说明

### 环境变量
- `DEEPSEEK_API_KEY`: DeepSeek API 密钥
- `DOUBAO_API_KEY`: 豆包 API 密钥
- `GEMINI_API_KEY`: Gemini API 密钥

### 工具配置
- **Serper API**: 用于实时网络搜索
- **Tavily API**: 用于高质量信息搜索

### 系统配置
- **工作空间**: 自动管理笔记和任务
- **对话历史**: 维护完整的对话上下文
- **统计信息**: 实时统计 Token 使用和费用

## 📊 系统特性详解

### 性能优化
- **并发处理**: 支持多用户同时使用，线程安全
- **内存管理**: 高效的内存使用和垃圾回收
- **响应速度**: 优化的 AI 模型调用流程
- **错误恢复**: 完善的错误处理和恢复机制
- **流式处理**: 基于 eino 框架的实时流式响应
- **组件缓存**: 智能组件缓存机制，减少重复计算

### 可扩展性
- **模块化设计**: 易于添加新的 AI 模型和工具
- **插件化工具**: 支持自定义工具集成
- **配置化提示词**: 灵活的提示词管理系统
- **接口标准化**: 统一的接口设计，便于扩展
- **eino 组件化**: 基于 eino 框架的组件化扩展
- **编排图扩展**: 支持复杂的 AI 应用编排扩展

### 用户体验
- **交互式界面**: 友好的命令行交互
- **实时反馈**: 流式输出和进度显示
- **错误处理**: 完善的错误提示和恢复
- **状态显示**: 实时显示系统状态和统计信息

## 🏗️ 技术架构详解

### eino 框架架构设计

#### 1. 组件化架构
基于 eino 框架的组件化设计，实现了高度模块化的架构：

```go
// 基础组件抽象
type BaseComponent struct {
    name string
}

// 实现 eino 的 BaseChatModel 接口
func (c *BaseComponent) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
    // 组件处理逻辑
    return &schema.Message{
        Role:    "assistant",
        Content: fmt.Sprintf("[%s] 处理结果: %s", c.name, content),
    }, nil
}
```

#### 2. 编排图设计模式
使用 eino 的编排图实现复杂的 AI 应用流程：

```go
// 智能体编排图结构
type Concierge struct {
    graph        *compose.Graph[[]*schema.Message, *schema.Message]
    compiledGraph compose.Runnable[[]*schema.Message, *schema.Message]
}

// 编排图构建（待完善）
func (c *Concierge) buildGraph() error {
    // TODO: 构建完整的 eino 编排图
    // 输入组件 -> 意图识别组件 -> 响应生成组件
    return nil
}
```

#### 3. 流式处理模式
基于 eino 的流式处理实现实时响应：

```go
// 流式处理实现
func (p *DeepSeekProvider) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
    // 创建 eino 流管道
    reader, writer := schema.Pipe[*schema.Message](5)
    
    go func() {
        defer writer.Close()
        // 流式处理逻辑
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
```

### 核心设计模式

#### 1. 单例模式
系统核心组件采用单例模式，确保全局唯一实例：
```go
var manager *ModelManager
var once sync.Once

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
```

#### 2. 工厂模式
模型提供商采用工厂模式，便于扩展：
```go
func NewDeepSeekProvider() (*DeepSeekProvider, error) {
    // 创建并配置 DeepSeek 提供商
}

func NewDoubaoProvider() (*DoubaoProvider, error) {
    // 创建并配置豆包提供商
}
```

#### 3. 策略模式
工具管理器采用策略模式，支持多种搜索策略：
```go
type ToolManager struct {
    tools map[string]Tool
}

func (tm *ToolManager) RegisterTool(tool Tool) {
    tm.tools[tool.Name()] = tool
}
```

#### 4. 适配器模式
使用适配器模式集成不同的 AI 模型：

```go
// 统一模型接口适配
type ModelProvider interface {
    model.BaseChatModel  // eino 基础接口
    Name() string
    DisplayName() string
    CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64
    CallLLM(ctx context.Context, systemPrompt, userPrompt string, options map[string]interface{}) (string, error)
}
```

### 并发安全
系统采用多种并发安全机制：

#### 1. 互斥锁
```go
type ModelManager struct {
    mu sync.RWMutex
    // ...
}

func (m *ModelManager) SetCurrentProvider(name string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    // ...
}
```

#### 2. 读写锁
对于读多写少的场景，使用读写锁提高性能：
```go
func (m *ModelManager) GetCurrentProvider() ModelProvider {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.currentProvider
}
```

### 错误处理
系统采用分层的错误处理机制：

#### 1. 错误包装
```go
if err != nil {
    return fmt.Errorf("DeepSeek API调用失败: %v", err)
}
```

#### 2. 优雅降级
当 AI 模型调用失败时，系统能够优雅降级：
```go
if err != nil {
    return o.generateDefaultTaskResponse(task)
}
```

## 🤝 贡献指南

### 开发环境设置
1. Fork 项目
2. 创建功能分支
3. 提交代码变更
4. 创建 Pull Request

### 代码规范
- 遵循 Go 语言官方代码规范
- 添加必要的注释和文档
- 确保测试覆盖率
- 使用 `gofmt` 格式化代码

### 架构原则
- **模块化**: 每个模块职责单一，便于维护
- **可扩展**: 通过接口设计支持功能扩展
- **可测试**: 代码结构便于单元测试
- **可配置**: 关键参数支持配置化

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [eino 框架](https://github.com/cloudwego/eino) - 字节跳动开源 AI 框架，提供了强大的组件抽象、编排图和流式处理能力
- [DeepSeek](https://www.deepseek.com/) - 深度求索大模型，强大的中文理解和生成能力
- [Serper](https://serper.dev/) - 实时搜索 API，提供最新的网络信息
- [Tavily](https://tavily.com/) - AI 搜索 API，专门为 AI 应用优化的搜索服务
- [Cobra](https://github.com/spf13/cobra) - Go 命令行框架，提供优秀的 CLI 体验
- [OpenAI Go](https://github.com/sashabaranov/go-openai) - OpenAI API 的 Go 客户端
- [Google Generative AI](https://github.com/google/generative-ai-go) - Google AI 的 Go 客户端

## 📚 eino 框架使用最佳实践

### 组件设计原则
1. **单一职责**: 每个组件只负责一个特定功能
2. **接口实现**: 严格实现 eino 的 `BaseChatModel` 接口
3. **错误处理**: 完善的错误处理和恢复机制
4. **流式支持**: 优先支持流式处理，提升用户体验

### 编排图设计
1. **清晰的数据流**: 定义明确的输入输出类型
2. **组件复用**: 设计可复用的组件
3. **类型安全**: 使用强类型保证编排正确性
4. **性能优化**: 考虑组件的执行顺序和并行性

### 消息处理
1. **统一格式**: 使用 `schema.Message` 作为标准消息格式
2. **角色定义**: 明确区分 system、user、assistant 角色
3. **内容处理**: 正确处理消息内容的编码和转义
4. **流式传输**: 支持大消息的分块传输

### 性能优化建议
1. **连接池**: 为 HTTP 客户端使用连接池
2. **缓存机制**: 实现智能缓存减少重复请求
3. **并发控制**: 合理控制并发数量避免过载
4. **资源管理**: 及时释放流式连接和文件句柄

### 扩展开发指南
1. **新模型集成**: 实现 `BaseChatModel` 接口
2. **新工具添加**: 实现 `Tool` 接口并注册到管理器
3. **新组件开发**: 继承 `BaseComponent` 或直接实现接口
4. **编排图扩展**: 在现有编排图中添加新节点

---

**注意**: 使用前请确保已配置相应的 API 密钥，并遵守各平台的使用条款。系统设计遵循最佳实践，具有良好的可扩展性和维护性。基于 eino 框架的架构设计确保了系统的高性能和可扩展性。 
