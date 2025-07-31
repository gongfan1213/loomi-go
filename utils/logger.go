package utils

import (
	"log"
	"os"
)

// InitLogger 初始化日志
func InitLogger() {
	// 创建日志文件
	logFile, err := os.OpenFile("loomi.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法创建日志文件:", err)
	}

	// 设置日志输出到文件
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// ClearScreen 清屏
func ClearScreen() {
	// 简单的清屏实现
	print("\033[H\033[2J")
} 