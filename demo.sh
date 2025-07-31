#!/bin/bash

# Loomi 2.0 演示脚本
# 展示基于 eino 框架的 Go 版本功能

echo "🎭 Loomi 2.0 演示脚本"
echo "基于 eino 框架的 Go 版本"
echo "================================"

# 检查可执行文件是否存在
if [ ! -f "./loomi2.0" ]; then
    echo "❌ 可执行文件不存在，请先运行 ./build.sh 构建程序"
    exit 1
fi

echo "✅ 找到可执行文件: ./loomi2.0"

# 显示版本信息
echo ""
echo "📋 版本信息:"
./loomi2.0 version

# 显示帮助信息
echo ""
echo "📚 帮助信息:"
./loomi2.0 --help

echo ""
echo "🎯 演示说明:"
echo "1. 运行 ./loomi2.0 start 启动系统"
echo "2. 选择您要使用的模型"
echo "3. 输入消息与AI对话"
echo "4. 使用特殊命令:"
echo "   - help: 显示帮助"
echo "   - status: 显示状态"
echo "   - clear: 清屏"
echo "   - orchestrator: 启动编排器"
echo "   - quit: 退出系统"

echo ""
echo "🚀 准备启动演示..."
echo "按 Enter 键开始演示，或按 Ctrl+C 取消"
read -r

# 启动演示
echo "🎬 启动 Loomi 2.0 演示..."
./loomi2.0 start 