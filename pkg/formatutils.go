package pkg

import (
	"fmt"
	"strconv"
)

func ReportErrorInputErr(overview string, detail error) (err error) {
	return fmt.Errorf("error when %s: %w", overview, detail)
}

func ReportErrorInputStr(overview string, detail string) (err error) {
	return fmt.Errorf("error when %s: %s", overview, detail)
}

func GetLimitIfValid(args []string, defaultLimit int) (limit int) {
	limit = defaultLimit

	if len(args) > 0 {
		num, err := strconv.Atoi(args[0])
		if err == nil {
			limit = num
			fmt.Println("use the limit => ", num)
		}
	}

	return
}
