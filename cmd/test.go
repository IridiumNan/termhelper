/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/IridiumNan/termhelper/internal/tester"
	"github.com/spf13/cobra"
)

const defaultLimit = 20

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	// This args not contains the `test` command itself
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")

		limit := defaultLimit

		if len(args) > 0 {
			num, err := strconv.Atoi(args[0])
			if err == nil {
				limit = num
				fmt.Println("use the limit => ", num)
			}
		}

		if limit < 4 {
			fmt.Println("minimal choice 4, current limit ", limit, " use 4 as the limit ")
			limit = 4
		}

		tester, err := tester.NewTester(limit)
		if err != nil {
			log.Fatal(err)
		}

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
