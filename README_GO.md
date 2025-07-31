# multi-agent小红书智能写作助手- 基于 eino 框架的 AI 助手系统

## 🚀 项目简介

multi-agent小红书智能写作助手是一个基于字节跳动 [eino 框架](https://github.com/cloudwego/eino) 构建的高性能 AI 助手系统。使用 Go 语言开发，支持多种 AI 模型，提供智能对话和任务编排功能。

## ✨ 主要特性

### 🤖 多模型支持
- **豆包 (Doubao-Seed-1.6-thinking)** - 字节跳动自研模型
- **DeepSeek** - 深度求索大模型
- **Gemini 1.5 Pro** - Google 最新模型
- 支持流式响应和费用统计

### 🏗️ eino 框架优势
- **图编排** - 使用有向图控制数据流
- **类型安全** - 编译时类型检查
- **流处理** - 支持 4 种流处理范式
- **组件抽象** - 丰富的组件生态
- **切面机制** - 支持回调和处理

### 🧠 智能体系统
- **门房智能体 (Concierge)** - 处理用户意图和基础对话
- **编排器智能体 (Orchestrator)** - 复杂任务分解和执行
- **工作空间管理** - 笔记、任务、上下文管理

### 📊 监控和统计
- 实时 Token 使用统计
- 费用计算和追踪
- 会话历史管理
- 系统状态监控

## 🛠️ 技术架构

```
Loomi 2.0
├── main.go                 # 主程序入口
├── cmd/                    # 命令行工具
│   ├── start.go           # 启动命令
│   └── version.go         # 版本命令
├── models/                 # 模型管理层
│   ├── provider.go        # 模型提供商接口
│   ├── manager.go         # 模型管理器
│   ├── doubao.go          # 豆包模型
│   ├── deepseek.go        # DeepSeek模型
│   ├── gemini.go          # Gemini模型
│   └── interface.go       # 全局接口
├── core/                   # 核心组件
│   ├── workspace.go       # 工作空间
│   └── conversation.go    # 对话管理
├── agents/                 # 智能体系统
│   ├── concierge.go       # 门房智能体
│   ├── orchestrator.go    # 编排器智能体
│   └── interface.go       # 智能体接口
└── utils/                  # 工具函数
    └── logger.go          # 日志工具
```

## 🚀 快速开始

### 环境要求
- Go 1.21+
- 各模型 API 密钥

### 安装依赖
```bash
go mod tidy
```

### 配置 API 密钥
编辑 `models/doubao.go`、`models/deepseek.go`、`models/gemini.go` 中的 API 密钥：

```go
const (
    DoubaoAPIKey = "your-doubao-api-key"
    DeepSeekAPIKey = "your-deepseek-api-key"
    GeminiAPIKey = "your-gemini-api-key"
)
```

### 运行程序
```bash
# 启动系统
go run main.go start

# 查看版本
go run main.go version

# 查看帮助
go run main.go --help
```

## 🎯 使用指南

### 基本命令
```
help, h          - 显示帮助信息
status           - 显示系统状态
clear            - 清屏
orchestrator     - 启动任务编排器
quit, exit, q    - 退出系统
```

### 模型选择
启动时会提示选择模型：
```
🤖 选择您要使用的模型：

1. Doubao-Seed-1.6-thinking
2. DeepSeek Chat
3. Gemini 1.5 Pro

请选择模型 (1-3):
```

### 对话示例
```
💬 请输入您的消息: 你好

🤖 multi-agent小红书智能写作助手: 我是multi-agent小红书智能写作助手的门房，可以为您提供以下服务：

1. 帮助信息 - 输入"帮助"或"help"
2. 系统状态 - 输入"状态"或"status"  
3. 笔记管理 - 输入"笔记"或"note"
4. 任务管理 - 输入"任务"或"task"
5. 一般对话 - 直接输入您的问题

请告诉我您需要什么帮助？
```

## 🔧 eino 框架集成

### 组件抽象
使用 eino 的组件抽象层：

```go
type ModelProvider interface {
    components.ChatModel
    
    Name() string
    DisplayName() string
    CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64
}
```

### 图编排
构建数据流编排图：

```go
func (m *ModelManager) buildGraph() error {
    m.graph = flow.NewGraph()
    
    // 添加节点
    modelNode := m.graph.AddNode("model", m.createModelComponent())
    templateNode := m.graph.AddNode("template", m.createTemplateComponent())
    toolsNode := m.graph.AddNode("tools", m.createToolsComponent())
    
    // 连接节点
    m.graph.AddEdge("template", "model")
    m.graph.AddEdge("model", "tools")
    
    // 编译图
    compiledGraph, err := m.graph.Compile(context.Background())
    if err != nil {
        return err
    }
    
    m.compiledGraph = compiledGraph
    return nil
}
```

### 流处理
支持 4 种流处理范式：

| 范式 | 说明 |
|------|------|
| Invoke | 非流输入 → 非流输出 |
| Stream | 非流输入 → 流输出 |
| Collect | 流输入 → 非流输出 |
| Transform | 流输入 → 流输出 |

## 📊 监控和统计

### 会话统计
```go
type SessionStats struct {
    TotalCalls         int     `json:"total_calls"`
    TotalInputTokens   int     `json:"total_input_tokens"`
    TotalOutputTokens  int     `json:"total_output_tokens"`
    TotalThinkingTokens int    `json:"total_thinking_tokens"`
    TotalCost          float64 `json:"total_cost"`
}
```

### 费用计算
每个模型提供商都实现了费用计算：

```go
func (p *DoubaoProvider) CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64 {
    inputCost := float64(inputTokens) / 1_000_000 * CostInputPer1M
    outputCost := float64(outputTokens+thinkingTokens) / 1_000_000 * CostOutputPer1M
    return inputCost + outputCost
}
```

## 🔄 与原 Python 版本的对比

### 优势
- **性能提升** - Go 语言的高并发性能
- **类型安全** - 编译时错误检查
- **内存效率** - 更低的内存占用
- **部署简单** - 单一二进制文件
- **框架优势** - eino 的编排能力

### 功能保持
- ✅ 多模型支持
- ✅ 智能体系统
- ✅ 工作空间管理
- ✅ 对话历史
- ✅ 费用统计
- ✅ 流式响应

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 📄 许可证

本项目采用 Apache-2.0 许可证。

## 🙏 致谢

- [eino 框架](https://github.com/cloudwego/eino) - 字节跳动开源
- [CloudWeGo](https://github.com/cloudwego) - 字节跳动微服务生态
- 原 Python 版本的贡献者

---

**Loomi 2.0** - 基于 eino 框架的下一代 AI 助手系统 🚀 
