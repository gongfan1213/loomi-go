package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"loomi2.0/cmd"
	"loomi2.0/core"
	"loomi2.0/models"
	"loomi2.0/utils"
)

func main() {
	// 初始化日志
	utils.InitLogger()

	// 创建根命令
	rootCmd := &cobra.Command{
		Use:   "assistant",
		Short: "小红书文案写作助手 - AI文案创作系统",
		Long: `小红书文案写作助手是一个基于eino框架的AI文案创作系统
支持多种AI模型，提供智能对话和任务编排功能`,
	}

	// 添加子命令
	rootCmd.AddCommand(cmd.StartCmd())
	rootCmd.AddCommand(cmd.VersionCmd())

	// 设置默认命令
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

// 初始化系统
func initSystem() error {
	// 初始化模型管理器
	if err := models.InitModelManager(); err != nil {
		return fmt.Errorf("初始化模型管理器失败: %v", err)
	}

	// 初始化工作空间
	if err := core.InitWorkspace(); err != nil {
		return fmt.Errorf("初始化工作空间失败: %v", err)
	}

	// 初始化对话管理器
	if err := core.InitConversationManager(); err != nil {
		return fmt.Errorf("初始化对话管理器失败: %v", err)
	}

	return nil
}

// 优雅关闭
func gracefulShutdown(ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("正在关闭系统...")
	
	// 清理资源
	core.Cleanup()
	models.Cleanup()
	
	os.Exit(0)
} 