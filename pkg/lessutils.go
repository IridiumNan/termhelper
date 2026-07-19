package pkg

import (
	"os"
	"os/exec"
)

func OpenWithLess(file *os.File) (err error) {
	lessCmd := exec.Command("less", "-R", file.Name())

	lessCmd.Stdin = os.Stdin
	lessCmd.Stdout = os.Stdout
	lessCmd.Stderr = os.Stderr

	if err = lessCmd.Run(); err != nil {
		return
	}

	return
}

func OpenWithLessByName(filePath string) (err error) {
	lessCmd := exec.Command("less", "-R", filePath)

	lessCmd.Stdin = os.Stdin
	lessCmd.Stdout = os.Stdout
	lessCmd.Stderr = os.Stderr

	if err = lessCmd.Run(); err != nil {
		return
	}

	return
}
