package pkg

import "fmt"

func ReportErrorInputErr(overview string, detail error) (err error) {
	return fmt.Errorf("error when %s: %w", overview, detail)
}

func ReportErrorInputStr(overview string, detail string) (err error) {
	return fmt.Errorf("error when %s: %s", overview, detail)
}
