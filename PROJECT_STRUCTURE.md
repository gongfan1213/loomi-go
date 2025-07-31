# Loomi 2.0 项目结构

## 📁 目录结构

```
loomi2.0/
├── main.go                    # 主程序入口
├── go.mod                     # Go 模块定义
├── go.sum                     # 依赖校验文件
├── build.sh                   # 构建脚本
├── demo.sh                    # 演示脚本
├── test_basic.go              # 基础测试
├── README_GO.md               # Go 版本说明文档
├── PROJECT_STRUCTURE.md       # 项目结构文档
├── cmd/                       # 命令行工具
│   ├── start.go              # 启动命令
│   └── version.go            # 版本命令
├── models/                    # 模型管理层
│   ├── provider.go           # 模型提供商接口
│   ├── manager.go            # 模型管理器
│   ├── doubao.go             # 豆包模型实现
│   ├── deepseek.go           # DeepSeek模型实现
│   ├── gemini.go             # Gemini模型实现
│   └── interface.go          # 全局接口函数
├── core/                      # 核心组件
│   ├── workspace.go          # 工作空间管理
│   └── conversation.go       # 对话管理器
├── agents/                    # 智能体系统
│   ├── concierge.go          # 门房智能体
│   ├── orchestrator.go       # 编排器智能体
│   └── interface.go          # 智能体接口
├── prompts/                   # 提示词系统 ⭐
│   ├── action_prompts.go     # 行动提示词
│   ├── concierge_prompt.go   # 门房提示词
│   └── orchestrator_prompt.go # 编排器提示词
└── utils/                     # 工具函数
    └── logger.go             # 日志工具
```

## 🔧 核心组件说明

### 1. 主程序 (main.go)
- 程序入口点
- 初始化日志系统
- 设置命令行工具
- 优雅关闭处理

### 2. 命令行工具 (cmd/)
- **start.go**: 启动命令，处理用户交互
- **version.go**: 版本信息显示

### 3. 模型管理层 (models/)
- **provider.go**: 定义模型提供商接口，继承 eino 的 ChatModel
- **manager.go**: 模型管理器，使用 eino 编排图
- **doubao.go**: 豆包模型实现
- **deepseek.go**: DeepSeek 模型实现
- **gemini.go**: Gemini 模型实现
- **interface.go**: 全局接口函数

### 4. 核心组件 (core/)
- **workspace.go**: 工作空间管理（笔记、任务、上下文）
- **conversation.go**: 对话管理器（消息历史、会话管理）

### 5. 智能体系统 (agents/)
- **concierge.go**: 门房智能体，处理用户意图
- **orchestrator.go**: 编排器智能体，处理复杂任务
- **interface.go**: 智能体接口和全局函数

### 6. 提示词系统 (prompts/)
- **action_prompts.go**: 行动相关提示词（洞察、画像、打点、文体）
- **concierge_prompt.go**: 门房智能体提示词
- **orchestrator_prompt.go**: 编排器智能体提示词

### 7. 工具函数 (utils/)
- **logger.go**: 日志初始化和工具函数

## 🏗️ eino 框架集成

### 组件抽象
```go
type ModelProvider interface {
    components.ChatModel
    
    Name() string
    DisplayName() string
    CalculateCost(inputTokens, outputTokens, thinkingTokens int) float64
}
```

### 图编排
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

### 流处理支持
- **Invoke**: 非流输入 → 非流输出
- **Stream**: 非流输入 → 流输出
- **Collect**: 流输入 → 非流输出
- **Transform**: 流输入 → 流输出

## 🔄 与原 Python 版本对比

### 功能保持
- ✅ 多模型支持（豆包、DeepSeek、Gemini）
- ✅ 智能体系统（门房、编排器）
- ✅ 工作空间管理（笔记、任务、上下文）
- ✅ 对话历史管理
- ✅ 费用统计和监控
- ✅ 流式响应支持
- ✅ CLI 交互界面

### 技术优势
- 🚀 **性能提升**: Go 语言的高并发性能
- 🔒 **类型安全**: 编译时错误检查
- 💾 **内存效率**: 更低的内存占用
- 📦 **部署简单**: 单一二进制文件
- 🏗️ **框架优势**: eino 的编排能力

## 🎯 使用流程

1. **构建程序**
   ```bash
   ./build.sh
   ```

2. **启动系统**
   ```bash
   ./loomi2.0 start
   ```

3. **选择模型**
   - 豆包 (Doubao-Seed-1.6-thinking)
   - DeepSeek Chat
   - Gemini 1.5 Pro

4. **开始对话**
   - 直接输入消息
   - 使用特殊命令（help、status、clear、orchestrator、quit）

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
每个模型提供商都实现了费用计算功能，支持实时统计和追踪。

## 🔧 配置说明

### API 密钥配置
编辑对应模型文件中的 API 密钥：
- `models/doubao.go`: 豆包 API 密钥
- `models/deepseek.go`: DeepSeek API 密钥
- `models/gemini.go`: Gemini API 密钥

### 费用配置
每个模型文件都包含费用计算常量：
```go
const (
    CostInputPer1M  = 0.12
    CostOutputPer1M = 1.2
)
```

## 🧪 测试

运行基础测试：
```bash
go test ./...
```

或运行特定测试：
```bash
go test -v -run TestBasicFunctionality
```

## 📝 开发说明

### 添加新模型
1. 在 `models/` 目录下创建新的模型文件
2. 实现 `ModelProvider` 接口
3. 在 `models/manager.go` 中注册新模型
4. 更新费用计算常量

### 添加新智能体
1. 在 `agents/` 目录下创建新的智能体文件
2. 实现 eino 组件接口
3. 在 `agents/interface.go` 中添加全局函数
4. 更新初始化流程

### 扩展功能
- 使用 eino 的图编排能力添加新的处理节点
- 利用流处理支持实现实时响应
- 通过切面机制添加日志和监控

---

**Loomi 2.0** - 基于 eino 框架的下一代 AI 助手系统 🚀 