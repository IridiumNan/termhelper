package extractor

import (
	"fmt"
	"strings"
	"time"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
)

const (
	defaultMaxWords    = 6
	defaultMaxChars    = 3000
	defaultCacheCount  = 3
	invalidateBoundary = -1

	disabledProficiency = 1.0

	notFoundIdx = -1
)

type Chunker struct {
	rawText string

	maxWordCount int
	maxCharCount int

	allWords models.WordStore

	currPos      int
	contextStart int

	providers []WordProvider

	currChunkPos int
	chunkCache   []*models.Chunk
}

func (c *Chunker) UpdateWords(response *models.LLMExplainResults) {
	for i := range response.Results {
		c.updateSingleWord(&response.Results[i])
	}
}

func (c *Chunker) updateSingleWord(word *models.WordResponse) {
	idx := c.allWords.IndexMap[word.Word]

	if idx >= len(c.allWords.Entries) {
		return
	}

	c.allWords.Entries[idx].SimpleDefinition = word.SimpleDefinition
	c.allWords.Entries[idx].DetailedExplanation = word.DetailedExplanation

	c.allWords.Entries[idx].NextReviewTime = time.Now().Unix()
	c.allWords.Entries[idx].UpdateReviewTime()
}

func (c *Chunker) AllWordEntries() []*models.WordEntry {
	var wordEntries []*models.WordEntry

	for _, word := range c.allWords.Entries {
		if word.Proficiency > 0.9 {
			continue
		}

		wordEntries = append(wordEntries, word)
	}

	return wordEntries
}

func (c *Chunker) getCurrNewWords() (words []string) {
	words = make([]string, 0)

	for _, i := range c.allWords.NewInDices {
		words = append(words, c.allWords.Entries[i].Word)
	}

	return
}

func (c *Chunker) pushNewWord(word *models.WordEntry) {
	if _, ok := c.allWords.IndexMap[word.Word]; ok {
		return
	}

	c.allWords.Entries = append(c.allWords.Entries, word)
	c.allWords.IndexMap[word.Word] = len(c.allWords.Entries) - 1
	c.allWords.IsNewMap[word.Word] = true
	c.allWords.NewInDices = append(c.allWords.NewInDices, c.allWords.IndexMap[word.Word])
}

func (c *Chunker) pushExistWord(word *models.WordEntry) {
	if _, ok := c.allWords.IndexMap[word.Word]; ok {
		return
	}

	c.allWords.Entries = append(c.allWords.Entries, word)
	c.allWords.IndexMap[word.Word] = len(c.allWords.Entries) - 1
	c.allWords.IsNewMap[word.Word] = false
}

func (c *Chunker) pushWord(word string) {
	cleanWord := strings.Trim(strings.ToLower(word), ".,!?;:()\"'")
	// 如果去除后为空，则忽略（如全是标点）
	if cleanWord == "" || len(cleanWord) < 2 {
		return
	}

	var wordEntry *models.WordEntry
	var found bool

	for idx := range c.providers {
		wordEntry, found = c.providers[idx].provide(cleanWord)

		if found {
			break
		}
	}

	if !found {
		c.pushNewWord(wordEntry)
		return
	}

	c.pushExistWord(wordEntry)
}

func NewChunker(text string) *Chunker {
	fmt.Println("disabled words:")
	fmt.Println(config.GetGlobalConfig().WordProvider.DisabledWords)

	fmt.Println()
	return &Chunker{
		rawText: text,

		maxWordCount: defaultMaxWords,
		maxCharCount: defaultMaxChars,

		allWords: models.WordStore{
			Entries:    []*models.WordEntry{},
			IndexMap:   map[string]int{},
			IsNewMap:   map[string]bool{},
			NewInDices: []int{},
		},

		currPos:      0,
		contextStart: 0,

		providers: []WordProvider{
			newDisabledProvider(config.GetGlobalConfig().WordProvider.DisabledWords),
			// TODO: dicProvider from database
			newDefaultProvider(),
		},
		currChunkPos: 0,
		chunkCache:   []*models.Chunk{},
	}
}

func (c *Chunker) HasNext() bool {
	return c.currPos < len(c.rawText)-1
}

func (c *Chunker) hasCache() bool {
	return c.CacheCount() > 0
}

func (c *Chunker) CacheCount() int {
	return len(c.chunkCache) - c.currChunkPos
}

func (c *Chunker) NextChunk() (chunk *models.Chunk, hasNext bool) {
	if !c.hasCache() {
		c.makeCache()
	}

	if c.CacheCount() == 0 {
		return nil, false
	}

	hasNext = c.hasCache()

	fmt.Println("cache length: ", len(c.chunkCache))
	fmt.Println("cache count: ", c.CacheCount())
	chunk = c.chunkCache[c.currChunkPos]
	c.currChunkPos++

	return
}

func isGap(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t' || r == '.' || r == '!' || r == ';' || r == ':' || r == '\'' || r == ','
}

func isEnd(r rune) bool {
	return r == '.' || r == '\n' || r == '?' || r == '!' || r == ';'
}

// getNextWord : return the idx of wordStart and wordEnd which is last Letter + 1
//
//	func getNextWordWithMark(text string, start int, end int) (wordStart int, wordEnd int, found bool) {
//		curr := start
//
//		for curr < end && isGap(rune(text[curr])) {
//			curr++
//		}
//
//		if curr >= end {
//			return notFoundIdx, notFoundIdx, false
//		}
//
//		wordStart = curr
//
//		for curr < end && !isGap(rune(text[curr])) {
//			curr++
//		}
//
//		wordEnd = curr
//		found = true
//
//		return
//	}
func getNextWordWithMark(text string, start int, end int) (wordStart int, wordEnd int, found bool) {
	curr := start

	for curr < end && isGap(rune(text[curr])) {
		curr++
	}

	if curr >= end {
		return notFoundIdx, notFoundIdx, false
	}

	wordStart = curr

	for curr < end && !isGap(rune(text[curr])) {
		curr++
	}

	wordEnd = curr
	found = true

	return
}

func (c *Chunker) makeCache() {
	c.contextStart = c.currPos

	for c.HasNext() && c.CacheCount() < defaultCacheCount {
		for c.HasNext() && len(c.allWords.NewInDices) < (c.maxWordCount/2) && c.currPos-c.contextStart < c.maxCharCount {
			wordStart, wordEnd, found := getNextWordWithMark(c.rawText, c.currPos, len(c.rawText))

			if !found {
				c.currPos++
				// break
			}

			c.pushWord(c.rawText[wordStart:wordEnd])
			c.currPos = wordEnd
		}

		for c.HasNext() && !isEnd(rune(c.rawText[c.currPos])) {
			wordStart, wordEnd, found := getNextWordWithMark(c.rawText, c.currPos, len(c.rawText))
			if !found {
				c.currPos++
				// break
			}
			c.pushWord(c.rawText[wordStart:wordEnd])
			c.currPos = wordEnd
		}

		if c.currPos < len(c.rawText)-1 && isEnd(rune(c.rawText[c.currPos])) {
			c.currPos++
		}

		c.pushCache(c.rawText[c.contextStart:c.currPos], c.getCurrNewWords())

		c.contextStart = c.currPos

		c.allWords.NewInDices = []int{}

	}
}

func (c *Chunker) pushCache(text string, words []string) {
	if len(words) > c.maxWordCount {
		c.pushCache(text, words[:len(words)/2])
		c.pushCache(text, words[len(words)/2:])

		return
	}
	// fmt.Println("chunkCache:")
	// fmt.Println("words: ", color.RedString("%v", words))
	// fmt.Println("text: ", color.GreenString("%s", text))

	c.chunkCache = append(c.chunkCache, &models.Chunk{
		Words:   words,
		Context: text,
	})
}
