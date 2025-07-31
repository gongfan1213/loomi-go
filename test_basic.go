package main

import (
	"testing"
	"loomi2.0/core"
	"loomi2.0/models"
)

// TestWorkspace 测试工作空间
func TestWorkspace(t *testing.T) {
	// 初始化工作空间
	err := core.InitWorkspace()
	if err != nil {
		t.Fatalf("初始化工作空间失败: %v", err)
	}

	workspace := core.GetWorkspace()
	if workspace == nil {
		t.Fatal("获取工作空间失败")
	}

	// 测试添加笔记
	workspace.AddNote("测试笔记")
	notes := workspace.GetNotes()
	if len(notes) != 1 {
		t.Errorf("笔记数量错误，期望 1，实际 %d", len(notes))
	}

	// 测试添加任务
	workspace.AddTask("测试任务")
	tasks := workspace.GetTasks()
	if len(tasks) != 1 {
		t.Errorf("任务数量错误，期望 1，实际 %d", len(tasks))
	}

	// 测试上下文
	workspace.SetContext("test_key", "test_value")
	value, exists := workspace.GetContext("test_key")
	if !exists {
		t.Error("上下文设置失败")
	}
	if value != "test_value" {
		t.Errorf("上下文值错误，期望 test_value，实际 %v", value)
	}
}

// TestConversationManager 测试对话管理器
func TestConversationManager(t *testing.T) {
	// 初始化对话管理器
	err := core.InitConversationManager()
	if err != nil {
		t.Fatalf("初始化对话管理器失败: %v", err)
	}

	conversation := core.GetConversationManager()
	if conversation == nil {
		t.Fatal("获取对话管理器失败")
	}

	// 测试添加消息
	conversation.AddMessage("user", "测试消息")
	messages := conversation.GetMessages()
	if len(messages) != 1 {
		t.Errorf("消息数量错误，期望 1，实际 %d", len(messages))
	}

	// 测试获取最后一条消息
	lastMessage := conversation.GetLastMessage()
	if lastMessage == nil {
		t.Fatal("获取最后一条消息失败")
	}
	if lastMessage.Content != "测试消息" {
		t.Errorf("最后一条消息内容错误，期望 测试消息，实际 %s", lastMessage.Content)
	}
}

// TestModelManager 测试模型管理器
func TestModelManager(t *testing.T) {
	// 初始化模型管理器
	err := models.InitModelManager()
	if err != nil {
		t.Fatalf("初始化模型管理器失败: %v", err)
	}

	manager := models.GetModelManager()
	if manager == nil {
		t.Fatal("获取模型管理器失败")
	}

	// 测试获取可用模型
	availableModels := models.GetAvailableModels()
	if len(availableModels) == 0 {
		t.Error("没有可用的模型")
	}

	// 测试会话统计
	stats := models.GetSessionStats()
	if stats.TotalCalls != 0 {
		t.Errorf("初始调用次数错误，期望 0，实际 %d", stats.TotalCalls)
	}
}

// TestBasicFunctionality 测试基本功能
func TestBasicFunctionality(t *testing.T) {
	t.Run("工作空间测试", TestWorkspace)
	t.Run("对话管理器测试", TestConversationManager)
	t.Run("模型管理器测试", TestModelManager)
} 