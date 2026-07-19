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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
