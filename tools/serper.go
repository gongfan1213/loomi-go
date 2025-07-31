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

// SerperTool Serper搜索工具
type SerperTool struct {
	apiKey string
	client *http.Client
}

// NewSerperTool 创建Serper工具实例
func NewSerperTool(apiKey string) *SerperTool {
	return &SerperTool{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name 工具名称
func (s *SerperTool) Name() string {
	return "serper_search"
}

// Description 工具描述
func (s *SerperTool) Description() string {
	return "使用Serper API进行网络搜索，获取最新的网络信息"
}

// Execute 执行搜索
func (s *SerperTool) Execute(ctx context.Context, query string) (string, error) {
	// 构建请求
	requestBody := map[string]interface{}{
		"q": query,
		"num": 10, // 返回10个结果
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}
	
	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", "https://google.serper.dev/search", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", s.apiKey)
	
	// 发送请求
	resp, err := s.client.Do(req)
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
	var serperResponse struct {
		Organic []struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
		} `json:"organic"`
	}
	
	if err := json.Unmarshal(body, &serperResponse); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	
	// 构建搜索结果
	var results []SearchResult
	for _, result := range serperResponse.Organic {
		results = append(results, SearchResult{
			Title:   result.Title,
			URL:     result.Link,
			Snippet: result.Snippet,
			Source:  "Serper",
		})
	}
	
	// 格式化输出
	output := fmt.Sprintf("🔍 Serper搜索结果 - 查询: %s\n\n", query)
	for i, result := range results {
		output += fmt.Sprintf("%d. **%s**\n", i+1, result.Title)
		output += fmt.Sprintf("   %s\n", result.Snippet)
		output += fmt.Sprintf("   链接: %s\n\n", result.URL)
	}
	
	return output, nil
} 