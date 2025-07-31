package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TavilyTool Tavily搜索工具
type TavilyTool struct {
	apiKey string
	client *http.Client
}

// NewTavilyTool 创建Tavily工具实例
func NewTavilyTool(apiKey string) *TavilyTool {
	return &TavilyTool{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name 工具名称
func (t *TavilyTool) Name() string {
	return "tavily_search"
}

// Description 工具描述
func (t *TavilyTool) Description() string {
	return "使用Tavily API进行网络搜索，获取高质量的网络信息"
}

// Execute 执行搜索
func (t *TavilyTool) Execute(ctx context.Context, query string) (string, error) {
	// 构建请求
	requestBody := map[string]interface{}{
		"query": query,
		"search_depth": "basic",
		"include_answer": false,
		"include_raw_content": false,
		"max_results": 10,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}
	
	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.tavily.com/search", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", t.apiKey)
	
	// 发送请求
	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}
	
	// 解析响应
	var tavilyResponse struct {
		Results []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Content     string `json:"content"`
			PublishedAt string `json:"published_at,omitempty"`
		} `json:"results"`
	}
	
	if err := json.Unmarshal(body, &tavilyResponse); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	
	// 构建搜索结果
	var results []SearchResult
	for _, result := range tavilyResponse.Results {
		results = append(results, SearchResult{
			Title:       result.Title,
			URL:         result.URL,
			Snippet:     result.Content,
			Source:      "Tavily",
			PublishedAt: result.PublishedAt,
		})
	}
	
	// 格式化输出
	output := fmt.Sprintf("🔍 Tavily搜索结果 - 查询: %s\n\n", query)
	for i, result := range results {
		output += fmt.Sprintf("%d. **%s**\n", i+1, result.Title)
		output += fmt.Sprintf("   %s\n", result.Snippet)
		output += fmt.Sprintf("   链接: %s\n", result.URL)
		if result.PublishedAt != "" {
			output += fmt.Sprintf("   发布时间: %s\n", result.PublishedAt)
		}
		output += "\n"
	}
	
	return output, nil
} 