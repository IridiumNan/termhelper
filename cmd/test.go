/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/IridiumNan/termhelper/internal/tester"
	"github.com/IridiumNan/termhelper/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const defaultTestWordLimit = 20

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "进行单词测验",
	Long: `从到期单词中随机抽取题目进行四选一交互测验。

每题显示一个英文单词，给出 4 个中文释义选项（A-D），
选择正确答案后显示详细解释。

答对增加熟练度，答错降低熟练度，
结果自动同步到数据库，影响下次复习间隔。

输入 Q 可随时退出测验。`,
	Example: `  termhelper test
  termhelper test 20
  termhelper test 50`,

	// This args not contains the `test` command itself
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")

		limit := pkg.GetLimitIfValid(args, defaultTestWordLimit)

		if limit < 4 {
			fmt.Println("minimal choice 4, current limit ", limit, " use 4 as the limit ")
			limit = 4
		}

		tester, err := tester.NewTester(limit)
		if err != nil {
			fmt.Println(err)
			fmt.Println("maybe you have finish all words review, you can use `termhelper status` command to see it")
			return
		}

		fmt.Println(color.GreenString(" TEST BEGIN, you can enter P for pass the word and it will not disapear again"))

		var hasNext bool
		for {
			hasNext = tester.NextTest()
			if !hasNext {
				return
			}
			fmt.Println("-----------------------------------------------------")
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
