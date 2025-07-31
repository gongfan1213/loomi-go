package core

import (
	"fmt"
	"sync"
	"time"
)

// Message 消息结构
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ConversationManager 对话管理器
type ConversationManager struct {
	mu       sync.RWMutex
	messages []Message
	session  string
}

var conversation *ConversationManager
var conversationOnce sync.Once

// InitConversationManager 初始化对话管理器
func InitConversationManager() error {
	var err error
	conversationOnce.Do(func() {
		conversation = &ConversationManager{
			messages: make([]Message, 0),
			session:  fmt.Sprintf("session_%d", time.Now().Unix()),
		}
	})
	return err
}

// GetConversationManager 获取对话管理器实例
func GetConversationManager() *ConversationManager {
	return conversation
}

// AddMessage 添加消息
func (c *ConversationManager) AddMessage(role, content string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	message := Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}
	
	c.messages = append(c.messages, message)
}

// GetMessages 获取所有消息
func (c *ConversationManager) GetMessages() []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	messages := make([]Message, len(c.messages))
	copy(messages, c.messages)
	return messages
}

// GetLastMessage 获取最后一条消息
func (c *ConversationManager) GetLastMessage() *Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if len(c.messages) == 0 {
		return nil
	}
	
	lastMessage := c.messages[len(c.messages)-1]
	return &lastMessage
}

// GetMessagesByRole 根据角色获取消息
func (c *ConversationManager) GetMessagesByRole(role string) []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	var filteredMessages []Message
	for _, msg := range c.messages {
		if msg.Role == role {
			filteredMessages = append(filteredMessages, msg)
		}
	}
	return filteredMessages
}

// GetConversationHistory 获取对话历史
func (c *ConversationManager) GetConversationHistory() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	var history string
	for _, msg := range c.messages {
		history += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}
	return history
}

// GetRecentMessages 获取最近的消息
func (c *ConversationManager) GetRecentMessages(count int) []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if count <= 0 || len(c.messages) == 0 {
		return []Message{}
	}
	
	if count > len(c.messages) {
		count = len(c.messages)
	}
	
	start := len(c.messages) - count
	messages := make([]Message, count)
	copy(messages, c.messages[start:])
	return messages
}

// Clear 清空对话
func (c *ConversationManager) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = make([]Message, 0)
}

// GetSessionID 获取会话ID
func (c *ConversationManager) GetSessionID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.session
}

// SetSessionID 设置会话ID
func (c *ConversationManager) SetSessionID(sessionID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.session = sessionID
}

// GetMessageCount 获取消息数量
func (c *ConversationManager) GetMessageCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.messages)
}

// GetConversationSummary 获取对话摘要
func (c *ConversationManager) GetConversationSummary() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	summary := fmt.Sprintf("对话摘要:\n")
	summary += fmt.Sprintf("- 会话ID: %s\n", c.session)
	summary += fmt.Sprintf("- 消息总数: %d\n", len(c.messages))
	
	// 统计各角色消息数量
	roleCount := make(map[string]int)
	for _, msg := range c.messages {
		roleCount[msg.Role]++
	}
	
	for role, count := range roleCount {
		summary += fmt.Sprintf("- %s消息: %d\n", role, count)
	}
	
	return summary
} 