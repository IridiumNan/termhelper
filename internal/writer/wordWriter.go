package writer

import (
	"fmt"
	"os"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/internal/output"
	"github.com/IridiumNan/termhelper/internal/storage"
)

const (
	initSkipCount    = 0
	defaultBatchSize = 20
)

// Split all words database for each block 50 words and write words to a tmp file => run less to check them
type WordWriter struct {
	wordData *storage.WordData

	skipCount int
}

func NewWordWriter() *WordWriter {
	return &WordWriter{
		wordData:  storage.GetGlobalWordData(),
		skipCount: initSkipCount,
	}
}

func batchFormat(words []*models.WordEntry, currPos int, batchSize int) (nextPos int, out string) {
	loopCount := min(len(words)-currPos, batchSize)

	nextPos = currPos + loopCount
	out = ""

	for i := range loopCount {
		out += fmt.Sprintf("\n%d", currPos+1+i)
		out += words[currPos+i].FileFormat()
	}

	return
}

func (w *WordWriter) WriteDueWords(wordLimit int, dst *os.File) (hasNext bool) {
	words, err := w.wordData.GetWordsSortByDue(wordLimit, w.skipCount)

	if len(words) == 0 {
		return false
	}

	if err != nil {
		output.Console.Error("error when get due words", "error", err)
		output.File.Error("error when get due words", "error", err)
		return false
	}

	currPos := 0
	var out string

	for currPos < len(words)-1 {

		currPos, out = batchFormat(words, currPos, defaultBatchSize)

		_, err = dst.WriteString(out)
		if err != nil {
			output.Console.Error("error when write due words batch to file", "file", dst.Name(), "error", err)
			output.File.Error("error when write due words batch to file", "file", dst.Name(), "error", err)
		}
	}

	return
}

func (w *WordWriter) WriteNextDueWord(dst *os.File) (hasNext bool, err error) {
	words, err := w.wordData.GetWordsSortByDue(1, w.skipCount)

	if len(words) == 0 {
		return false, err
	}
	if err != nil {
		output.Console.Error("error when get next due word", "error", err)
		output.File.Error("error when get next due word", "error", err)

		return false, err
	}

	_, err = dst.Write([]byte(words[0].FileFormat()))
	if err != nil {
		output.Console.Error("error when write next due word", "error", err)
	}

	w.skipCount++

	return
}
