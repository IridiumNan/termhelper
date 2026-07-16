package extractor

import (
	"strings"
	"unicode"

	"github.com/IridiumNan/termhelper/internal/models"
)

const (
	defaultMaxWords    = 5
	defaultMaxChars    = 3000
	invalidateBoundary = -1
)

var defaultFilter = MemoryFilter{
	Words: []string{
		"to",
		"move",
		"the",
	},
}

type Chunk models.LLMRequest

func (chunk *Chunk) ChunkToRequest() *models.LLMRequest {
	return &models.LLMRequest{
		Word:    chunk.Word,
		Context: chunk.Context,
	}
}

type Chunker struct {
	rawText      []rune
	currPos      int // curr_pos
	maxWords     int
	maxChars     int
	isEmpty      bool
	wordInChunk  []string
	contextStart int
	filters      []Filter
}

func (chunker *Chunker) AppendIfPass(newWord string) {
	for i := range chunker.filters {
		if chunker.filters[i].isDuplicate(newWord) {
			return
		}
	}

	for i := range chunker.wordInChunk {
		if strings.Contains(chunker.wordInChunk[i], newWord) {
			return
		}
	}

	chunker.isEmpty = false
	chunker.wordInChunk = append(chunker.wordInChunk, newWord)
}

func isSpace(r rune) bool {
	return r == ' '
}

func isTargetMark(r rune) bool {
	return r == '\n' || r == '\'' || r == '\t' || r == ' ' || r == '.'
}

func striptMark(wordWithMark string) string {
	for isTargetMark(rune(wordWithMark[0])) {
		if len(wordWithMark) <= 1 {
			break
		}
		wordWithMark = wordWithMark[1:]
	}

	for isTargetMark(rune(wordWithMark[len(wordWithMark)-1])) {
		if len(wordWithMark) <= 1 {
			break
		}
		wordWithMark = wordWithMark[:len(wordWithMark)-2]
	}

	return wordWithMark
}

func NewChunker(raw string, maxWords int, maxChars int) *Chunker {
	if maxWords <= 0 {
		maxWords = defaultMaxWords
	}
	if maxChars <= 0 {
		maxChars = defaultMaxChars
	}

	return &Chunker{
		rawText:      []rune(raw),
		currPos:      0,
		maxWords:     maxWords,
		maxChars:     maxChars,
		isEmpty:      true,
		wordInChunk:  []string{},
		contextStart: 0,
		filters:      []Filter{&defaultFilter},
	}
}

func (c *Chunker) NextChunk() *Chunk {
	c.contextStart := c.currPos

	lastBoundaryIdx := c.currPos

	for c.currPos < len(c.rawText)-1 && (len(c.wordInChunk)%c.maxWords == 0 || c.currPos-c.contextStart < c.maxChars) && !c.isEmpty {
		for c.currPos < len(c.rawText) && isSpace(c.rawText[c.currPos]) {
			c.currPos++
		}
		wordStart := c.currPos

		for !isSpace(c.rawText[c.currPos]) {
			c.currPos++

			if isBoundary(c.rawText[c.currPos]) {
				lastBoundaryIdx = c.currPos
			}
		}

		c.AppendIfPass(string(c.rawText[wordStart : c.currPos-1]))

		c.currPos++
	}

	if c.currPos == len(c.rawText)-1 {
		return &Chunk{}
	}
}

func isBoundary(ch rune) bool {
	return ch == '!' || ch == '\n' || ch == '.' || ch == '?'
}

func (chunker *Chunker) findSentenceBoundary() int {
	if chunker.currPos <= chunker.contextStart {
		return -1
	}

	searchStart := chunker.currPos

	if searchStart-chunker.contextStart > chunker.maxChars {
		searchStart = chunker.contextStart + chunker.maxChars
	}

	for i := searchStart - 1; i > chunker.contextStart; i-- {
		ch := chunker.rawText[i]
		if isBoundary(ch) {
			if i+1 < len(chunker.rawText) || (chunker.rawText[i+1] == ' ' || chunker.rawText[i+1] == '\n' || chunker.rawText[i+1] == '\t') {
				return i + 1
			}

			if i+1 == len(chunker.rawText) {
				return i + 1
			}
		}
	}

	for i := searchStart - 1; i > chunker.contextStart; i-- {
		if chunker.rawText[i] == ' ' || chunker.rawText[i] == '\n' || chunker.rawText[i] == '\t' {
			return i + 1
		}
	}

	return invalidateBoundary
}

func (chunker *Chunker) finish() *Chunk {
	if chunker.currPos >= len(chunker.rawText) && len(chunker.wordInChunk) == 0 {
		return nil
	}

	chunkText := string(chunker.rawText[chunker.contextStart:])
	wordsInChunk := chunker.extractWordsFromBuffer(chunker.rawText[chunker.contextStart:])

	chunker.currPos = len(chunker.rawText)

	chunker.textInBuffer = nil

	chunker.wordInChunk = []string{}

	return &Chunk{
		Word:    wordsInChunk,
		Context: chunkText,
	}
}

// return all words with deduplicate
func (chunker *Chunker) extractWordsFromBuffer(text []rune) []string {
	seen := make(map[string]bool)

	var words []string
	i := 0

	for i < len(text) {
		if unicode.IsLetter(text[i]) {
			start := i

			for i < len(text) && unicode.IsLetter(text[i]) {
				i++
			}

			word := strings.ToLower(string(text[start:i]))

			if len(word) > 1 && !seen[word] {
				seen[word] = true
				words = append(words, word)

			} else {
				i++
			}
		}
	}

	return words
}
