package pkg

import (
	"os"
)

func IsEmptyStr(str string) bool {
	return str == ""
}

func EnsureDir(path string) (err error) {
	if IsEmptyStr(path) {
		return ReportErrorInputStr("EnsurePath", "empty path")
	}
	err = os.MkdirAll(path, 0o755)
	return
}
