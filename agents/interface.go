package agents

import (
	"context"
	"fmt"
)

// InitAgents 初始化所有智能体
func InitAgents() error {
	// 初始化门房智能体
	if err := InitConcierge(); err != nil {
		return fmt.Errorf("初始化门房智能体失败: %v", err)
	}

	// 初始化编排器智能体
	if err := InitOrchestrator(); err != nil {
		return fmt.Errorf("初始化编排器智能体失败: %v", err)
	}

	return nil
}

// ProcessUserInput 处理用户输入（全局函数）
func ProcessUserInput(ctx context.Context, userInput string) (string, error) {
	// 首先尝试使用门房智能体处理
	concierge := GetConcierge()
	if concierge != nil {
		response, err := concierge.ProcessUserInput(ctx, userInput)
		if err == nil {
			return response, nil
		}
	}

	// 如果门房处理失败，使用编排器处理
	orchestrator := GetOrchestrator()
	if orchestrator != nil {
		response, err := orchestrator.ProcessTask(ctx, userInput)
		if err == nil {
			return response, nil
		}
	}

	return "", fmt.Errorf("所有智能体都无法处理用户输入")
}

// StartOrchestrator 启动编排器
func StartOrchestrator(ctx context.Context) error {
	orchestrator := GetOrchestrator()
	if orchestrator == nil {
		return fmt.Errorf("编排器未初始化")
	}
	return orchestrator.StartOrchestrator(ctx)
}

// StopOrchestrator 停止编排器
func StopOrchestrator() {
	orchestrator := GetOrchestrator()
	if orchestrator != nil {
		orchestrator.StopOrchestrator()
	}
}

// IsOrchestratorRunning 检查编排器是否正在运行
func IsOrchestratorRunning() bool {
	orchestrator := GetOrchestrator()
	if orchestrator == nil {
		return false
	}
	return orchestrator.IsRunning()
} 