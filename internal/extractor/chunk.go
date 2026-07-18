package extractor

import (
	"fmt"
	"strings"
	"time"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/pkg"
	"github.com/fatih/color"
)

const (
	defaultMaxWords    = 6
	defaultMaxChars    = 3000
	defaultCacheCount  = 3
	invalidateBoundary = -1

	disabledProficiency = 1.0

	noDisplayProficiency = 0.9

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

// UpdateWords : use the LLM Explainations to update new words attributes
func (c *Chunker) UpdateWords(response *models.LLMExplainResults) {
	for i := range response.Results {

		// if word exist, skip it
		if !c.allWords.IsNewMap[response.Results[i].Word] {
			continue
		}

		// just update new words
		c.updateSingleWord(&response.Results[i])
	}
}

// appendWordResponse: this func is for LLM response don't use the original word but split it
func (c *Chunker) appendWordResponse(word *models.WordResponse) {
	wordEntry := models.WordEntry{
		Word:                word.Word,
		SimpleDefinition:    word.SimpleDefinition,
		DetailedExplanation: word.DetailedExplanation,
		Proficiency:         models.InitProficiency,
		NextReviewTime:      time.Now().Unix(),
	}

	wordEntry.UpdateReviewTime()

	c.allWords.Entries = append(c.allWords.Entries, &wordEntry)
	c.allWords.IndexMap[word.Word] = len(c.allWords.Entries) - 1

	c.allWords.IsNewMap[word.Word] = true
}

// updateSingleWord: update the word which is new, their Word attribute is not empty and other is unset
func (c *Chunker) updateSingleWord(word *models.WordResponse) {
	idx, ok := c.allWords.IndexMap[word.Word]

	// for preserve data, you shound always use the Guard Clauses
	if !ok {
		// if the LLM's response is not match the original word, append it as the new word Instead then update the entry

		c.appendWordResponse(word)
		// We handle this word then return
		return
	}

	if idx >= len(c.allWords.Entries) {
		fmt.Println(color.RedString("out of range word entry: %v", word))
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
		if word.Proficiency > noDisplayProficiency || pkg.IsEmptyStr(word.SimpleDefinition) || pkg.IsEmptyStr(word.DetailedExplanation) {
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

	// push to the NewInDices
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
	// if it's not word, return
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
			// disabled provider return the word entry with hight proficiency to avoid displaying them
			newDisabledProvider(config.GetGlobalConfig().WordProvider.DisabledWords),

			// dic provider which Get words from dicts
			newDicProvider(),

			// default provider for new words construct
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
			}

			c.pushWord(c.rawText[wordStart:wordEnd])
			c.currPos = wordEnd
		}

		for c.HasNext() && !isEnd(rune(c.rawText[c.currPos])) {
			wordStart, wordEnd, found := getNextWordWithMark(c.rawText, c.currPos, len(c.rawText))
			if !found {
				c.currPos++
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
		// clear current NewInDices slice

	}
}

func (c *Chunker) pushCache(text string, words []string) {
	if len(words) > c.maxWordCount {
		c.pushCache(text, words[:len(words)/2])
		c.pushCache(text, words[len(words)/2:])

		return
	}

	c.chunkCache = append(c.chunkCache, &models.Chunk{
		Words:   words,
		Context: text,
	})
}
