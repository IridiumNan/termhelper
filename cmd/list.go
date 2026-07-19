/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/internal/output"
	"github.com/IridiumNan/termhelper/internal/writer"
	"github.com/IridiumNan/termhelper/pkg"
	"github.com/spf13/cobra"
)

const defaultWordListLimit = 30

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出待复习的到期单词",
	Long: `从数据库中查询所有到期（到复习时间）的单词，
按下次复习时间排序输出到临时文件，并通过 less 分页查看。

单词按熟练度着色显示：
  - 绿色  → 新词（New）
  - 黄色  → 学习中（Learning）
  - 青色  → 较熟悉（Familiar）

不传参数时默认列出 30 个单词。`,
	Example: `  termhelper list
  termhelper list 50
  termhelper list 100`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		// db := storage.GetGlobalWordData()
		//
		// words, err := db.GetAllWordEntries()
		// if err != nil {
		// 	log.Fatal(err)
		// 	return
		// }
		//
		// for _, word := range words {
		// 	word.ColorfulPrint()
		// }
		//
		// defer db.Close()
		limit := pkg.GetLimitIfValid(args, defaultWordListLimit)

		writer := writer.NewWordWriter()

		file, err := os.CreateTemp(models.GetTmpDirPath(), "list-*.tmp")

		defer os.Remove(file.Name())

		if err != nil {
			output.Console.Error("fail to create temp file", "error", err)
			output.File.Error("fail to create temp file", "error", err)
			return
		}

		hasNext := writer.WriteDueWords(limit, file)
		if err != nil {
			fmt.Println(err)

			return
		}

		// This may not work => list [num] just list the num limit words
		// but is reasonable
		for hasNext {
			hasNext = writer.WriteDueWords(limit, file)
		}

		file.Close()

		err = pkg.OpenWithLess(file)
		if err != nil {
			output.Console.Error("fail to use less command to open file", "file", file.Name(), "error", err)
		}

		fmt.Println("clear the tmp file => ", file.Name())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
