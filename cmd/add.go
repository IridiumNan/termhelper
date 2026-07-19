/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/IridiumNan/termhelper/internal/extractor"
	"github.com/IridiumNan/termhelper/internal/llm"
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/internal/output"
	"github.com/IridiumNan/termhelper/internal/storage"
	"github.com/IridiumNan/termhelper/internal/writer"
	"github.com/IridiumNan/termhelper/pkg"
	"github.com/spf13/cobra"
)

const emptyStr = ""

var fileFlag string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readInput(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fail to read input: %v\n", err)
		}

		if pkg.IsEmptyStr(input) {
			fmt.Fprintln(os.Stderr, "error: there is no any text")
			// TODO: print the usage
			os.Exit(1)
		}

		// fmt.Printf("receive the input: %s\n", input)
		chunker := extractor.NewChunker(input)

		client, err := llm.NextClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)

		defer cancel()

		for {
			chunk, hasNext := chunker.NextChunk()

			// print this for debug
			// fmt.Println(chunk)
			// fmt.Println(hasNext)

			if !hasNext {
				break
			}
			results, err := client.Explain(ctx, chunk)
			if err != nil {
				fmt.Println(err)
				continue
			}

			chunker.UpdateWords(results)

			// fmt.Println(color.YellowString("%v", results))

		}
		// for _, word := range chunker.AllWordEntries() {
		// 	// if word == nil {
		// 	// 	fmt.Println(color.RedString("*models.WordEntry is nil pointer"))
		// 	// }
		// 	word.ColorfulPrint()
		// }

		db := storage.GetGlobalWordData()
		if db == nil {
			return
		}

		db.PushAllWordEntriesIfNew(chunker.AllWordEntries())

		err = db.Close()
		if err != nil {
			fmt.Println("error when close db", err)
		}

		fileName := writeRawWithOutput(input, chunker.AllWordEntries())

		err = pkg.OpenWithLessByName(fileName)
		if err != nil {
			output.Console.Error("fail to open file with less command ", "error", err)
		}
		fmt.Println("you can look result later => less -R ", fileName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "从文本中读取文件")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func readInput(args []string) (string, error) {
	fmt.Println("fileFlag -> ", fileFlag)
	if !pkg.IsEmptyStr(fileFlag) {
		data, err := os.ReadFile(fileFlag)
		if err != nil {
			return "", fmt.Errorf("fail to read file: %w", err)
		}

		return string(data), nil

	}

	if len(args) > 0 {
		return strings.Join(args, ""), nil
	}

	return readStdin()
}

func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return emptyStr, err
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return emptyStr, nil
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return emptyStr, fmt.Errorf("fail to read std in: %w", err)
	}

	return string(data), nil
}

func ReplaceRaw(rawText string, words []*models.WordEntry) string {
	var colorfulText string
	for _, word := range words {
		colorfulText = strings.ReplaceAll(rawText, word.Word, output.RenderWordByProficiency(word.Word, word.Proficiency))
	}

	return colorfulText
}

func writeRawWithOutput(rawText string, words []*models.WordEntry) (tmpFilePath string) {
	file, err := os.CreateTemp(models.GetTmpDirPath(), "add-output-*.tmp")
	if err != nil {
		output.Console.Error("error when create temp file", "file", file.Name(), "error", err)
		output.File.Error("error when create temp file", "file", file.Name(), "error", err)
		return
	}
	//
	// _, err = file.WriteString(ReplaceRaw(rawText, words))
	// if err != nil {
	// 	output.Console.Error("error when write rawText to temp file", "file", file.Name())
	// 	output.File.Error("error when write rawText to temp file", "file", file.Name())
	// }

	writer := writer.NewWordWriter()

	_ = writer.WriteWords(words, file)

	return file.Name()
}
