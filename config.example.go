package main

// 配置示例文件
// 请复制此文件为 config.go 并填入实际的 API 密钥

// API 密钥配置示例
const (
	// DeepSeek API 密钥
	// 获取地址: https://platform.deepseek.com/
	DeepSeekAPIKey = "your-deepseek-api-key"
	
	// 豆包 API 密钥  
	// 获取地址: https://www.doubao.com/
	DoubaoAPIKey = "your-doubao-api-key"
	
	// Gemini API 密钥
	// 获取地址: https://makersuite.google.com/app/apikey
	GeminiAPIKey = "your-gemini-api-key"
	
	// Serper API 密钥
	// 获取地址: https://serper.dev/
	SerperAPIKey = "your-serper-api-key"
	
	// Tavily API 密钥
	// 获取地址: https://tavily.com/
	TavilyAPIKey = "your-tavily-api-key"
)

// 使用说明:
// 1. 复制此文件为 config.go
// 2. 填入实际的 API 密钥
// 3. 在相应的模型文件中引用这些常量
// 4. 确保 config.go 已添加到 .gitignore 中 