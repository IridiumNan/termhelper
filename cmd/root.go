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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "termhelper",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

func initConfig() error {
	funcName := "initConfig"

	output.InitConsoleLogger(config.GetGlobalConfig().Logging.ConsoleLevel)
	err := output.InitFileLogger(config.GetGlobalConfig().Logging.FileLevel)
	if err != nil {
		return pkg.ReportErrorInputErr(funcName+" init file logger", err)
	}

	err = config.LoadConfig()
	if err != nil {
		return pkg.ReportErrorInputErr(funcName+" Load config", err)
	}

	if err := initDataDir(); err != nil {
		return pkg.ReportErrorInputErr(funcName, err)
	}

	return nil
}
