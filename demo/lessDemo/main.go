// demo_less.go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// 1. 生成带颜色的内容
	content := generateColoredContent()

	// 2. 创建临时文件
	tmpFile, err := os.CreateTemp("", "demo-*.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建临时文件失败: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tmpFile.Name()) // 程序退出时删除文件

	// 3. 写入内容
	if _, err := tmpFile.WriteString(content); err != nil {
		fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		os.Exit(1)
	}
	tmpFile.Close()

	// 4. 用 less -R 打开文件
	cmd := exec.Command("less", "-R", tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("打开 less 查看彩色文本，按 q 退出...")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "less 执行失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 已退出 less，临时文件已删除。")
}

// generateColoredContent 生成带 ANSI 颜色的示例文本
func generateColoredContent() string {
	const (
		reset  = "\033[0m"
		red    = "\033[31m"
		green  = "\033[32m"
		yellow = "\033[33m"
		cyan   = "\033[36m"
	)

	return fmt.Sprintf(`这是一个彩色输出测试。

%s绿色文本%s
%s黄色文本%s
%s红色文本%s
%s青色文本%s

你可以看到这些颜色在 less -R 中都能正常显示。

如果你看到的是普通文本而不是彩色，说明 less 没有正确解析 ANSI 转义序列。
请确认你使用了 less -R，并且终端支持颜色。
`, green, reset, yellow, reset, red, reset, cyan, reset)
}
