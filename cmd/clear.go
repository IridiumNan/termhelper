/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "清理临时文件",
	Long: `删除 /tmp/termhelper/ 目录下的所有临时文件。

在执行 add 命令时，工具会在 /tmp/termhelper/ 下生成临时输出文件，
这些文件可通过本命令集中清理。

删除前会交互式确认，输入 y 确认删除，n 取消。`,
	Example: `  termhelper clear`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clear called")

		fmt.Println("clear the data on", models.GetTmpDirPath())

		if isRemove := askConfirm(); isRemove {
			if err := os.RemoveAll(models.GetTmpDirPath()); err != nil {
				fmt.Printf("error when clear dir: %s, err: %v\n", models.GetTmpDirPath(), err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func askConfirm() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("remove the dir %s [y/N] ->", models.GetTmpDirPath())
		input, _ := reader.ReadString('\n')
		input = strings.ReplaceAll(input, "\n", "")

		choice := unicode.ToLower(rune(input[0]))

		if choice == 'y' {
			return true
		}

		if choice == 'n' {
			return false
		}

	}
}
