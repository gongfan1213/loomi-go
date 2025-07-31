package tools

import (
	"context"
	"fmt"
	"strings"
)

// ToolManager 工具管理器
type ToolManager struct {
	tools map[string]Tool
}

// NewToolManager 创建工具管理器
func NewToolManager() *ToolManager {
	return &ToolManager{
		tools: make(map[string]Tool),
	}
}

// RegisterTool 注册工具
func (tm *ToolManager) RegisterTool(tool Tool) {
	tm.tools[tool.Name()] = tool
}

// GetTool 获取工具
func (tm *ToolManager) GetTool(name string) (Tool, bool) {
	tool, exists := tm.tools[name]
	return tool, exists
}

// ListTools 列出所有工具
func (tm *ToolManager) ListTools() []Tool {
	var tools []Tool
	for _, tool := range tm.tools {
		tools = append(tools, tool)
	}
	return tools
}

// ExecuteTool 执行工具
func (tm *ToolManager) ExecuteTool(ctx context.Context, toolName, input string) (string, error) {
	tool, exists := tm.tools[toolName]
	if !exists {
		return "", fmt.Errorf("工具 %s 不存在", toolName)
	}
	
	return tool.Execute(ctx, input)
}

// DetectSearchIntent 检测搜索意图
func (tm *ToolManager) DetectSearchIntent(userInput string) (bool, string) {
	searchKeywords := []string{
		"搜索", "查找", "查询", "了解", "搜索关于", "查找关于", "查询关于",
		"search", "find", "lookup", "search for", "find about",
	}
	
	userInputLower := strings.ToLower(userInput)
	for _, keyword := range searchKeywords {
		if strings.Contains(userInputLower, keyword) {
			// 提取搜索查询
			query := extractSearchQuery(userInput, keyword)
			return true, query
		}
	}
	
	return false, ""
}

// extractSearchQuery 提取搜索查询
func extractSearchQuery(userInput, keyword string) string {
	// 简单的查询提取逻辑
	userInputLower := strings.ToLower(userInput)
	keywordIndex := strings.Index(userInputLower, strings.ToLower(keyword))
	
	if keywordIndex == -1 {
		return userInput
	}
	
	// 提取关键词后面的内容
	startIndex := keywordIndex + len(keyword)
	if startIndex >= len(userInput) {
		return userInput
	}
	
	query := strings.TrimSpace(userInput[startIndex:])
	
	// 移除常见的后缀词
	suffixes := []string{"的内容", "的信息", "的资料", "的新闻", "的资讯"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(query, suffix) {
			query = strings.TrimSuffix(query, suffix)
			break
		}
	}
	
	return query
}

// PerformDualSearch 执行双重搜索
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
	
	if len(results) == 0 {
		return "❌ 没有可用的搜索工具", nil
	}
	
	return strings.Join(results, "\n\n" + strings.Repeat("=", 50) + "\n\n"), nil
} 