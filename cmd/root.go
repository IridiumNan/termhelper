/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/internal/output"
	"github.com/IridiumNan/termhelper/pkg"
	"github.com/spf13/cobra"
)

var versionFlag string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "termhelper",
	Short: "终端英语词汇辅助学习工具",
	Long: `termhelper 帮助计算机专业学生和技术人员在日常工作中学习英文技术词汇。

它从报错信息、技术文档、终端输出等真实文本中提取英文词汇，
借助 LLM（Ollama / DeepSeek / OpenAI）生成中文释义，
并通过间隔重复（Spaced Repetition）机制进行复习巩固。

常用命令：
  termhelper add    解析文本并提取词汇
  termhelper list   列出待复习的到期单词
  termhelper test   进行单词测验
  termhelper status 查看学习状态
  termhelper config 编辑配置文件
  termhelper clear  清理临时文件`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.termhelper.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	err := initConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	initOutput()

	if err := initDir(); err != nil {
		fmt.Println(err)
		return
	}
}

func initOutput() {
	output.InitConsoleLogger(config.GetGlobalConfig().Logging.ConsoleLevel)
	err := output.InitFileLogger(config.GetGlobalConfig().Logging.FileLevel)
	if err != nil {
		fmt.Println(err)
	}
}

func initConfig() error {
	funcName := "initConfig"
	err := config.LoadConfig()
	if err != nil {
		return pkg.ReportErrorInputErr(funcName+" Load config", err)
	}
	return nil
}

func initDir() error {
	funcName := "initDir"

	if err := initDataDir(); err != nil {
		return pkg.ReportErrorInputErr(funcName, err)
	}

	if err := initTmpDir(); err != nil {
		return pkg.ReportErrorInputErr(funcName, err)
	}

	return nil
}

func initDataDir() error {
	dataPath := models.GetDataDirPath()
	if err := pkg.EnsureDir(dataPath); err != nil {
		return pkg.ReportErrorInputErr("initDataPath", err)
	}

	filePath := models.GetDataFilePath()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err = os.Create(filePath)
		if err != nil {
			return pkg.ReportErrorInputErr("initDataDir", err)
		}
		output.File.Info("data file not found create a new one", "path", filePath)
	}

	return nil
}

func initTmpDir() error {
	tmpDir := models.GetTmpDirPath()

	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return err
	}

	return nil
}
