/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/IridiumNan/termhelper/internal/storage"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看单词学习状态统计",
	Long: `显示当前单词学习的整体状态概览，包括：

  - 待复习单词数量（已到期未复习）
  - 各熟练度级别的单词分布：
    新词（New）→ 学习中（Learning）→ 较熟悉（Familiar）→ 已掌握（Master）

帮助你了解当前学习进度和复习压力。`,
	Example: `  termhelper status`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")

		status := ""

		db := storage.GetGlobalWordData()

		status += db.GetDueWordsStatus()

		status += db.GetProficiencyStatus()

		fmt.Println(status)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
