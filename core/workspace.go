package core

import (
	"fmt"
	"sync"
)

// WorkSpace 工作空间
type WorkSpace struct {
	mu       sync.RWMutex
	notes    []string
	tasks    []string
	context  map[string]interface{}
}

var workspace *WorkSpace
var workspaceOnce sync.Once

// InitWorkspace 初始化工作空间
func InitWorkspace() error {
	var err error
	workspaceOnce.Do(func() {
		workspace = &WorkSpace{
			notes:   make([]string, 0),
			tasks:   make([]string, 0),
			context: make(map[string]interface{}),
		}
	})
	return err
}

// GetWorkspace 获取工作空间实例
func GetWorkspace() *WorkSpace {
	return workspace
}

// AddNote 添加笔记
func (w *WorkSpace) AddNote(note string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.notes = append(w.notes, note)
}

// GetNotes 获取所有笔记
func (w *WorkSpace) GetNotes() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	notes := make([]string, len(w.notes))
	copy(notes, w.notes)
	return notes
}

// AddTask 添加任务
func (w *WorkSpace) AddTask(task string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.tasks = append(w.tasks, task)
}

// GetTasks 获取所有任务
func (w *WorkSpace) GetTasks() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	tasks := make([]string, len(w.tasks))
	copy(tasks, w.tasks)
	return tasks
}

// SetContext 设置上下文
func (w *WorkSpace) SetContext(key string, value interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.context[key] = value
}

// GetContext 获取上下文
func (w *WorkSpace) GetContext(key string) (interface{}, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	value, exists := w.context[key]
	return value, exists
}

// GetAllContext 获取所有上下文
func (w *WorkSpace) GetAllContext() map[string]interface{} {
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	context := make(map[string]interface{})
	for k, v := range w.context {
		context[k] = v
	}
	return context
}

// Clear 清空工作空间
func (w *WorkSpace) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.notes = make([]string, 0)
	w.tasks = make([]string, 0)
	w.context = make(map[string]interface{})
}

// GetSummary 获取工作空间摘要
func (w *WorkSpace) GetSummary() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	summary := fmt.Sprintf("工作空间状态:\n")
	summary += fmt.Sprintf("- 笔记数量: %d\n", len(w.notes))
	summary += fmt.Sprintf("- 任务数量: %d\n", len(w.tasks))
	summary += fmt.Sprintf("- 上下文键数量: %d\n", len(w.context))
	
	return summary
}

// Cleanup 清理资源
func Cleanup() {
	// 清理工作空间资源
} 