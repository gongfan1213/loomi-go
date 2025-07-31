package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  "显示Loomi系统的版本信息",
	Run:   runVersion,
}

func VersionCmd() *cobra.Command {
	return versionCmd
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println("Loomi 2.0")
	fmt.Println("基于 eino 框架构建")
	fmt.Println("版本: v2.0.0")
	fmt.Println("GitHub: https://github.com/cloudwego/eino")
} 