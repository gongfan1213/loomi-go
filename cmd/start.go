package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"loomi2.0/agents"
	"loomi2.0/core"
	"loomi2.0/models"
	"loomi2.0/utils"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动文案写作助手",
	Long:  "启动小红书文案写作助手，进入交互模式",
	Run:   runStart,
}

func StartCmd() *cobra.Command {
	return startCmd
}

func runStart(cmd *cobra.Command, args []string) {
	// 显示启动信息
	showStartupInfo()

	// 初始化系统
	if err := initSystem(); err != nil {
		color.Red("❌ 系统初始化失败: %v", err)
		os.Exit(1)
	}

	// 选择模型
	if err := selectModel(); err != nil {
		color.Red("❌ 模型选择失败: %v", err)
		os.Exit(1)
	}

	// 启动交互循环
	startInteractiveLoop()
}

func showStartupInfo() {
	color.Cyan("🚀 小红书文案写作助手启动中...")
	color.Cyan("基于 eino 框架构建")
	color.Cyan(strings.Repeat("=", 50))
}

func initSystem() error {
	color.Green("🔧 初始化系统组件...")
	
	// 初始化模型管理器
	if err := models.InitModelManager(); err != nil {
		return fmt.Errorf("模型管理器初始化失败: %v", err)
	}
	color.Green("✅ 模型管理器初始化完成")

	// 初始化工作空间
	if err := core.InitWorkspace(); err != nil {
		return fmt.Errorf("工作空间初始化失败: %v", err)
	}
	color.Green("✅ 工作空间初始化完成")

	// 初始化对话管理器
	if err := core.InitConversationManager(); err != nil {
		return fmt.Errorf("对话管理器初始化失败: %v", err)
	}
	color.Green("✅ 对话管理器初始化完成")

	// 初始化智能体
	if err := agents.InitAgents(); err != nil {
		return fmt.Errorf("智能体初始化失败: %v", err)
	}
	color.Green("✅ 智能体初始化完成")

	return nil
}

func selectModel() error {
	availableModels := models.GetAvailableModels()
	if len(availableModels) == 0 {
		return fmt.Errorf("没有可用的模型")
	}

	color.Cyan("\n🤖 选择您要使用的模型：")
	color.Cyan("")

	// 显示可用模型
	modelNames := make([]string, 0, len(availableModels))
	for name, displayName := range availableModels {
		modelNames = append(modelNames, name)
		color.Yellow("%d. %s", len(modelNames), displayName)
	}

	color.Cyan("")
	
	// 读取用户选择
	reader := bufio.NewReader(os.Stdin)
	color.Green("请选择模型 (1-%d): ", len(modelNames))
	
	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("读取用户输入失败: %v", err)
	}

	choice = strings.TrimSpace(choice)
	if choice == "" {
		choice = "1" // 默认选择第一个
	}

	// 验证选择
	var selectedIndex int
	if _, err := fmt.Sscanf(choice, "%d", &selectedIndex); err != nil {
		return fmt.Errorf("无效的选择: %s", choice)
	}

	if selectedIndex < 1 || selectedIndex > len(modelNames) {
		return fmt.Errorf("选择超出范围: %d", selectedIndex)
	}

	selectedModel := modelNames[selectedIndex-1]
	selectedDisplayName := availableModels[selectedModel]

	// 设置当前模型
	if err := models.SetCurrentModel(selectedModel); err != nil {
		return fmt.Errorf("设置模型失败: %v", err)
	}

	color.Green("✅ 已选择模型: %s", selectedDisplayName)
	return nil
}

func startInteractiveLoop() {
	color.Cyan("\n🎯 系统已启动，输入 'help' 查看帮助，输入 'quit' 退出")
	color.Cyan(strings.Repeat("=", 60))

	reader := bufio.NewReader(os.Stdin)
	
	for {
		color.Cyan("\n💬 请输入您的消息: ")
		
		input, err := reader.ReadString('\n')
		if err != nil {
			color.Red("❌ 读取输入失败: %v", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// 处理特殊命令
		if handleSpecialCommands(input) {
			continue
		}

		// 处理用户输入
		if err := handleUserInput(input); err != nil {
			color.Red("❌ 处理用户输入失败: %v", err)
		}
	}
}

func handleSpecialCommands(input string) bool {
	switch strings.ToLower(input) {
	case "quit", "exit", "q":
		color.Yellow("👋 再见！")
		os.Exit(0)
		return true
	case "help", "h":
		showHelp()
		return true
	case "status":
		showStatus()
		return true
	case "clear":
		utils.ClearScreen()
		return true
	case "orchestrator", "orch":
		startOrchestrator()
		return true
	}
	return false
}

func handleUserInput(input string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用eino框架处理用户输入
	response, err := agents.ProcessUserInput(ctx, input)
	if err != nil {
		return fmt.Errorf("处理用户输入失败: %v", err)
	}

	// 显示响应
	color.Green("\n🤖 Loomi: %s", response)
	
	return nil
}

func showHelp() {
	helpText := `
📚 可用命令:
  help, h          - 显示此帮助信息
  status           - 显示系统状态
  clear            - 清屏
  orchestrator     - 启动任务编排器
  quit, exit, q    - 退出系统

💡 提示:
  - 直接输入消息与AI对话
  - 支持多轮对话
  - 系统会自动保存对话历史
`
	color.Cyan(helpText)
}

func showStatus() {
	stats := models.GetSessionStats()
	color.Cyan("\n📊 系统状态:")
	color.Cyan("  总调用次数: %d", stats.TotalCalls)
	color.Cyan("  总输入Token: %d", stats.TotalInputTokens)
	color.Cyan("  总输出Token: %d", stats.TotalOutputTokens)
	color.Cyan("  总思考Token: %d", stats.TotalThinkingTokens)
	color.Cyan("  总费用: $%.4f", stats.TotalCost)
	
	currentModel := models.GetCurrentModelName()
	color.Cyan("  当前模型: %s", currentModel)
}

func startOrchestrator() {
	color.Yellow("🚀 启动任务编排器...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := agents.StartOrchestrator(ctx); err != nil {
		color.Red("❌ 启动编排器失败: %v", err)
		return
	}

	color.Green("✅ 编排器已启动，可以开始任务编排")
} 