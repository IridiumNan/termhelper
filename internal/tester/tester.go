package tester

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"
	"unicode"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/internal/storage"
	"github.com/fatih/color"
)

type mark rune

const marksCount = 4

var allMarks = []mark{'A', 'B', 'C', 'D'}

type option struct {
	m    mark
	word *models.WordEntry
}

func newOption(m mark, word *models.WordEntry) *option {
	return &option{
		m:    m,
		word: word,
	}
}

func (opt *option) printQuiz() {
	out := fmt.Sprintf(
		"%s) %s",
		string(opt.m),
		opt.word.SimpleDefinition,
	)

	fmt.Println(out)
}

func (opt *option) printDetail(c *color.Color) string {
	return c.Sprintf(
		"%s: %s => %s",
		opt.word.Word,
		opt.word.SimpleDefinition,
		opt.word.DetailedExplanation,
	)
}

type Tester struct {
	wordData *storage.WordData

	wordPool []*models.WordEntry

	candidatesIdx []int

	currPos int

	targetMarkIdx int

	reader bufio.Reader
}

func NewTester(limit int) (tester *Tester, err error) {
	tester = &Tester{
		wordData:      storage.GetGlobalWordData(),
		wordPool:      []*models.WordEntry{},
		candidatesIdx: []int{},
		currPos:       0,
		targetMarkIdx: -1,
		reader:        *bufio.NewReader(os.Stdin),
	}

	dueWords, err := tester.wordData.GetDueWords(limit)
	if err != nil {
		return nil, fmt.Errorf("error when get due words %w", err)
	}
	tester.wordPool = append(tester.wordPool, dueWords...)

	return
}

func (t *Tester) genRandWordFromPool() {
	for len(t.candidatesIdx) < 3 {
		i := rand.Intn(len(t.wordPool))

		if i != t.currPos && !slices.Contains(t.candidatesIdx, i) {
			t.candidatesIdx = append(t.candidatesIdx, i)
		}
	}
}

func (t *Tester) updateTargetMarkIdx() {
	t.targetMarkIdx = rand.Intn(marksCount)
}

func (t *Tester) clearCandidates() {
	t.candidatesIdx = []int{}
}

func (t *Tester) getOptions() (options []*option) {
	t.genRandWordFromPool()

	t.updateTargetMarkIdx()

	options = make([]*option, 0)

	// TODO: append four words

	markIdx := 0

	for i := range 3 {

		if markIdx == t.targetMarkIdx {
			options = append(options, newOption(allMarks[markIdx], t.wordPool[t.currPos]))

			markIdx++
		}

		options = append(options, newOption(allMarks[markIdx], t.wordPool[t.candidatesIdx[i]]))

		markIdx++

	}

	if t.targetMarkIdx == 3 {
		options = append(options, newOption(allMarks[3], t.wordPool[t.currPos]))
	}

	return
}

func (t *Tester) printOptions(options []*option) {
	for _, op := range options {
		op.printQuiz()
	}
}

func (t *Tester) askThenPrintAnswer(opts []*option) (quit bool) {
	var choice mark
	choice, quit, mask := t.askAnswer()

	if mask {
		// mask this word
		t.wordPool[t.currPos].MaskWithLongTime()

		fmt.Println(color.RedString("you have masked this word => %s", t.wordPool[t.currPos].Word))
		return
	}

	if quit {
		return
	}

	t.wordPool[t.currPos].UpdateProficiency(choice == allMarks[t.targetMarkIdx])

	t.wordPool[t.currPos].UpdateReviewTime()

	for i, op := range opts {
		if i == t.targetMarkIdx {
			fmt.Println(op.printDetail(color.New(color.FgGreen)))
			continue
		}

		if allMarks[i] == choice {
			fmt.Println(op.printDetail(color.New(color.FgRed)))

			continue
		}

		fmt.Println(op.printDetail(color.New(color.FgBlue)))
	}

	return
}

func (t *Tester) askAnswer() (choice mark, quit bool, mask bool) {
	for {
		fmt.Print("enter your choice [A-D] (enter q for quit) ->")
		input, _ := t.reader.ReadString('\n')

		input = strings.ReplaceAll(input, "\n", "")

		if len(input) == 0 {
			continue
		}
		choice = mark(unicode.ToUpper(rune(input[0])))

		if choice == mark('Q') {
			quit = true
			return
		}

		if choice == mark('P') {
			mask = true
			return
		}

		if slices.Contains(allMarks, choice) {
			break
		}
	}

	return
}

func (t *Tester) syncToDB() {
	fmt.Println("sync to db")
	t.wordData.UpdateAll(t.wordPool[:t.currPos])
}

func (t *Tester) NextTest() (hasNext bool) {
	if t.currPos >= len(t.wordPool) {
		t.syncToDB()
		return false
	}

	options := t.getOptions()

	fmt.Println("word => ", t.wordPool[t.currPos].Word)

	t.printOptions(options)

	quit := t.askThenPrintAnswer(options)

	if quit {
		t.syncToDB()
		return false
	}
	// clear the candidatesIdx
	t.clearCandidates()
	// Move the currPos for next word
	t.currPos++

	return true
}
