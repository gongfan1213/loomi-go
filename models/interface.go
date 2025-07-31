package models

import (
	"context"
	"fmt"
)

// GetAvailableModels 获取可用模型
func GetAvailableModels() map[string]string {
	manager := GetModelManager()
	if manager == nil {
		return make(map[string]string)
	}
	return manager.GetProviderDisplayNames()
}

// SetCurrentModel 设置当前模型
func SetCurrentModel(modelName string) error {
	manager := GetModelManager()
	if manager == nil {
		return fmt.Errorf("模型管理器未初始化")
	}
	return manager.SetCurrentProvider(modelName)
}

// GetCurrentModelName 获取当前模型名称
func GetCurrentModelName() string {
	manager := GetModelManager()
	if manager == nil {
		return ""
	}
	provider := manager.GetCurrentProvider()
	if provider == nil {
		return ""
	}
	return provider.Name()
}

// GetSessionStats 获取会话统计
func GetSessionStats() SessionStats {
	manager := GetModelManager()
	if manager == nil {
		return SessionStats{}
	}
	return manager.GetStats()
}

// CallLLM 调用LLM（全局函数）
func CallLLM(ctx context.Context, systemPrompt, userPrompt string, options map[string]interface{}) (string, error) {
	manager := GetModelManager()
	if manager == nil {
		return "", fmt.Errorf("模型管理器未初始化")
	}
	return manager.CallCurrentModel(ctx, systemPrompt, userPrompt, options)
}

// Cleanup 清理资源
func Cleanup() {
	manager := GetModelManager()
	if manager != nil {
		manager.Cleanup()
	}
} 