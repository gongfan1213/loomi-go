#!/bin/bash

# Loomi 2.0 构建脚本
# 基于 eino 框架的 Go 版本

set -e

echo "🚀 Loomi 2.0 构建脚本"
echo "基于 eino 框架的 Go 版本"
echo "================================"

# 检查 Go 版本
echo "📋 检查环境..."
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.21+"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go 版本: $GO_VERSION"

# 检查是否在正确的目录
if [ ! -f "go.mod" ]; then
    echo "❌ 请在项目根目录运行此脚本"
    exit 1
fi

# 清理旧的构建文件
echo "🧹 清理旧的构建文件..."
rm -f loomi2.0
rm -rf dist/

# 下载依赖
echo "📦 下载依赖..."
go mod tidy

# 运行测试
echo "🧪 运行测试..."
go test ./...

# 构建程序
echo "🔨 构建程序..."
go build -o loomi2.0 main.go

# 检查构建结果
if [ -f "loomi2.0" ]; then
    echo "✅ 构建成功！"
    echo "📁 可执行文件: ./loomi2.0"
    
    # 显示文件信息
    echo "📊 文件信息:"
    ls -lh loomi2.0
    
    # 显示版本信息
    echo "📋 版本信息:"
    ./loomi2.0 version
    
    echo ""
    echo "🎯 使用方法:"
    echo "  ./loomi2.0 start     # 启动系统"
    echo "  ./loomi2.0 version   # 查看版本"
    echo "  ./loomi2.0 --help    # 查看帮助"
    
else
    echo "❌ 构建失败！"
    exit 1
fi

echo ""
echo "🎉 构建完成！"
echo "现在可以运行 ./loomi2.0 start 启动系统" 