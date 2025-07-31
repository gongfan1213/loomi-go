# Loomi 2.0 Go 版本设置指南

## 🚀 环境要求

### 必需软件
- **Go 1.21+** - Go 编程语言
- **Git** - 版本控制工具

### 推荐软件
- **VS Code** 或 **GoLand** - 代码编辑器
- **Docker** - 容器化部署（可选）

## 📦 安装 Go

### macOS (使用 Homebrew)
```bash
# 安装 Homebrew (如果未安装)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装 Go
brew install go

# 验证安装
go version
```

### macOS (手动安装)
1. 访问 [Go 官网](https://golang.org/dl/)
2. 下载 macOS 版本的 Go
3. 运行安装程序
4. 验证安装：`go version`

### Linux (Ubuntu/Debian)
```bash
# 更新包列表
sudo apt update

# 安装 Go
sudo apt install golang-go

# 验证安装
go version
```

### Windows
1. 访问 [Go 官网](https://golang.org/dl/)
2. 下载 Windows 版本的 Go
3. 运行安装程序
4. 重启命令提示符
5. 验证安装：`go version`

## 🔧 环境配置

### 设置 GOPATH (可选)
```bash
# 添加到 ~/.bashrc 或 ~/.zshrc
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

### 验证 Go 环境
```bash
# 检查 Go 版本
go version

# 检查 Go 环境
go env

# 检查 Go 模块支持
go mod --help
```

## 🏗️ 构建项目

### 1. 进入项目目录
```bash
cd loomi_go
```

### 2. 下载依赖
```bash
go mod tidy
```

### 3. 构建项目
```bash
# 使用构建脚本
./build.sh

# 或手动构建
go build -o loomi2.0 main.go
```

### 4. 运行测试
```bash
go test ./...
```

## 🚀 运行项目

### 启动系统
```bash
./loomi2.0 start
```

### 查看版本
```bash
./loomi2.0 version
```

### 查看帮助
```bash
./loomi2.0 --help
```

## 🔧 配置 API 密钥

在运行项目前，需要配置各模型的 API 密钥：

### 1. 豆包 API
编辑 `models/doubao.go`：
```go
const (
    DoubaoAPIKey = "your-doubao-api-key"
)
```

### 2. DeepSeek API
编辑 `models/deepseek.go`：
```go
const (
    DeepSeekAPIKey = "your-deepseek-api-key"
)
```

### 3. Gemini API
编辑 `models/gemini.go`：
```go
const (
    GeminiAPIKey = "your-gemini-api-key"
)
```

## 🧪 开发环境

### 安装开发工具
```bash
# 安装代码格式化工具
go install golang.org/x/tools/cmd/goimports@latest

# 安装代码检查工具
go install golang.org/x/lint/golint@latest

# 安装测试覆盖率工具
go install golang.org/x/tools/cmd/cover@latest
```

### 代码格式化
```bash
go fmt ./...
goimports -w .
```

### 代码检查
```bash
golint ./...
```

### 测试覆盖率
```bash
go test -cover ./...
```

## 📦 依赖管理

### 添加新依赖
```bash
go get github.com/example/package
```

### 更新依赖
```bash
go get -u github.com/example/package
```

### 清理依赖
```bash
go mod tidy
```

## 🐳 Docker 部署 (可选)

### 创建 Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o loomi2.0 main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/loomi2.0 .
CMD ["./loomi2.0", "start"]
```

### 构建 Docker 镜像
```bash
docker build -t loomi2.0 .
```

### 运行 Docker 容器
```bash
docker run -it loomi2.0
```

## 🔍 故障排除

### 常见问题

1. **go: command not found**
   - 确保 Go 已正确安装
   - 检查 PATH 环境变量

2. **module not found**
   - 运行 `go mod tidy`
   - 检查 `go.mod` 文件

3. **build failed**
   - 检查 Go 版本是否满足要求
   - 确保所有依赖都已下载

4. **API 调用失败**
   - 检查 API 密钥是否正确配置
   - 检查网络连接

### 获取帮助
- [Go 官方文档](https://golang.org/doc/)
- [Go 模块文档](https://golang.org/doc/modules/)
- [eino 框架文档](https://github.com/cloudwego/eino)

---

**Loomi 2.0 Go 版本** - 基于 eino 框架的下一代 AI 助手系统 🚀 