/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/spf13/cobra"
)

var currEditorIdx = 0

var editorOptions = []string{
	"vim",
	"vi",
	"nano",
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "编辑配置文件",
	Long: `使用系统默认编辑器（$EDITOR）打开 TOML 格式的配置文件进行编辑。

按优先级依次尝试以下编辑器：vim → vi → nano

配置文件位于 ~/.config/termhelper/config.toml，
支持配置 LLM 提供商、模型参数、复习策略、日志级别等。

也可通过环境变量覆盖配置：
  DEEPSEEK_API_KEY、OPENAI_API_KEY、TERMHELPER_*`,
	Example: `  termhelper config`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")

		editor := getValidEditro()

		configFile := models.GetConfigFilePath()

		editCommand := exec.Command(editor, configFile)

		editCommand.Stdin = os.Stdin
		editCommand.Stdout = os.Stdout
		editCommand.Stderr = os.Stderr

		if err := editCommand.Run(); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("you can also edit your config here => ", models.GetConfigFilePath())
	},
}

func getValidEditro() string {
	editor := getEnvEditor()

	exist := checkIfExist(editor)

	for !exist {

		currEditorIdx++

		exist = checkIfExist(editorOptions[currEditorIdx])
	}
	return editor
}

func getEnvEditor() string {
	return os.Getenv("EDITOR")
}

func checkIfExist(editor string) bool {
	findCmd := exec.Command("which", editor)

	if err := findCmd.Run(); err != nil {
		return false
	}

	return true
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
