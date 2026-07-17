package extractor

import (
	"fmt"
	"strings"

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

type Chunk struct {
	Words   []string
	Context string
}

type WordStore struct {
	entries   []*models.WordEntry
	indexMap  map[string]int
	isNewMap  map[string]bool
	newIndics []int
}

type Chunker struct {
	rawText string

	maxWordCount int
	maxCharCount int

	allWords WordStore

	currPos      int
	contextStart int

	providers []WordProvider

	currChunkPos int
	chunkCache   []*Chunk
}

func (c *Chunker) getCurrNewWords() (words []string) {
	words = make([]string, 0)

	for _, i := range c.allWords.newIndics {
		words = append(words, c.allWords.entries[i].Word)
	}

	return
}

func (c *Chunker) pushNewWord(word *models.WordEntry) {
	if _, ok := c.allWords.indexMap[word.Word]; ok {
		return
	}

	c.allWords.entries = append(c.allWords.entries, word)
	c.allWords.indexMap[word.Word] = len(c.allWords.entries) - 1
	c.allWords.isNewMap[word.Word] = true
	c.allWords.newIndics = append(c.allWords.newIndics, c.allWords.indexMap[word.Word])
}

func (c *Chunker) pushExistWord(word *models.WordEntry) {
	if _, ok := c.allWords.indexMap[word.Word]; ok {
		return
	}

	c.allWords.entries = append(c.allWords.entries, word)
	c.allWords.indexMap[word.Word] = len(c.allWords.entries) - 1
	c.allWords.isNewMap[word.Word] = false
}

func (c *Chunker) pushWord(word string) {
	cleanWord := strings.Trim(strings.ToLower(word), ".,!?;:()\"'")
	// 如果去除后为空，则忽略（如全是标点）
	if cleanWord == "" {
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

		allWords: WordStore{
			entries:   []*models.WordEntry{},
			indexMap:  map[string]int{},
			isNewMap:  map[string]bool{},
			newIndics: []int{},
		},

		currPos:      0,
		contextStart: 0,

		providers: []WordProvider{
			newDisabledProvider(config.GetGlobalConfig().WordProvider.DisabledWords),
			// TODO: dicProvider from database
			newDefaultProvider(),
		},
		currChunkPos: 0,
		chunkCache:   []*Chunk{},
	}
}

func (c *Chunker) HasNext() bool {
	return c.currPos < len(c.rawText)-1 || c.cacheCount() > 0
}

func (c *Chunker) hasCache() bool {
	return c.cacheCount() > 0
}

func (c *Chunker) cacheCount() int {
	return len(c.chunkCache) - c.currChunkPos - 1
}

func (c *Chunker) NextChunk() (chunk *Chunk, hasNext bool) {
	if !c.HasNext() {
		return nil, false
	}

	if !c.hasCache() {
		c.makeCache()
	}
	c.currChunkPos++
	return c.chunkCache[c.currChunkPos-1], c.HasNext()
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

//	func (c *Chunker) makeCache() {
//		c.contextStart = c.currPos
//		fmt.Printf("[MAKE] start: currPos=%d, contextStart=%d, textLen=%d\n", c.currPos, c.contextStart, len(c.rawText))
//
//		for c.HasNext() && c.cacheCount() < defaultCacheCount {
//			fmt.Printf("[MAKE] outer loop: currPos=%d, cacheCount=%d, newIndics len=%d\n",
//				c.currPos, c.cacheCount(), len(c.allWords.newIndics))
//
//			// ---- 第一个内层循环 ----
//			loop1 := 0
//			for c.HasNext() && len(c.allWords.newIndics) < (c.maxWordCount/2) && c.currPos-c.contextStart < c.maxCharCount {
//				loop1++
//				fmt.Printf("[LOOP1] iter=%d, currPos=%d, newIndics=%d\n", loop1, c.currPos, len(c.allWords.newIndics))
//				wordStart, wordEnd, found := getNextWordWithMark(c.rawText, c.currPos, len(c.rawText))
//				fmt.Printf("[LOOP1] getNextWord -> start=%d end=%d found=%v\n", wordStart, wordEnd, found)
//				if !found {
//					fmt.Printf("[LOOP1] !found, currPos %d -> %d, break\n", c.currPos, c.currPos+1)
//					if c.currPos < len(c.rawText)-1 {
//						c.currPos++
//					} else {
//						return
//					}
//					break
//				}
//				if c.currPos >= len(c.rawText) {
//					return
//				}
//				word := c.rawText[wordStart:wordEnd]
//				fmt.Printf("[LOOP1] pushing word %q\n", word)
//				c.pushWord(word)
//				c.currPos = wordEnd
//				fmt.Printf("[LOOP1] new currPos=%d\n", c.currPos)
//			}
//			fmt.Printf("[MAKE] after LOOP1: currPos=%d, newIndics=%d\n", c.currPos, len(c.allWords.newIndics))
//
//			// ---- 第二个内层循环 ----
//			loop2 := 0
//			for c.HasNext() && !isEnd(rune(c.rawText[c.currPos])) {
//				loop2++
//				fmt.Printf("[LOOP2] iter=%d, currPos=%d, char=%q\n", loop2, c.currPos, c.rawText[c.currPos])
//				wordStart, wordEnd, found := getNextWordWithMark(c.rawText, c.currPos, len(c.rawText))
//				fmt.Printf("[LOOP2] getNextWord -> start=%d end=%d found=%v\n", wordStart, wordEnd, found)
//				if !found {
//					fmt.Printf("[LOOP2] !found, currPos %d -> %d, break\n", c.currPos, c.currPos+1)
//					if c.currPos < len(c.rawText)-1 {
//						c.currPos++
//					} else {
//						return
//					}
//					c.currPos++
//					break
//				}
//				word := c.rawText[wordStart:wordEnd]
//				fmt.Printf("[LOOP2] pushing word %q\n", word)
//				c.pushWord(word)
//				c.currPos = wordEnd
//				fmt.Printf("[LOOP2] new currPos=%d\n", c.currPos)
//			}
//			fmt.Printf("[MAKE] after LOOP2: currPos=%d, isEnd? %v\n", c.currPos,
//				c.HasNext() && isEnd(rune(c.rawText[c.currPos])))
//
//			// ---- 生成缓存 ----
//			newWords := c.getCurrNewWords()
//			context := c.rawText[c.contextStart:c.currPos]
//			fmt.Printf("[MAKE] pushCache: contextLen=%d, newWords=%v\n", len(context), newWords)
//			c.pushCache(context, newWords)
//			c.allWords.newIndics = []int{}
//			fmt.Printf("[MAKE] after pushCache: chunkCache len=%d, currChunkPos=%d\n",
//				len(c.chunkCache), c.currChunkPos)
//		}
//		fmt.Printf("[MAKE] exit: currPos=%d, cacheCount=%d\n", c.currPos, c.cacheCount())
//	}
func (c *Chunker) makeCache() {
	c.contextStart = c.currPos

	for c.HasNext() && c.cacheCount() < defaultCacheCount {
		for c.HasNext() && len(c.allWords.newIndics) < (c.maxWordCount/2) && c.currPos-c.contextStart < c.maxCharCount {
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

		if isEnd(rune(c.rawText[c.currPos])) {
			c.currPos++
		}

		c.pushCache(c.rawText[c.contextStart:c.currPos], c.getCurrNewWords())

		c.contextStart = c.currPos

		c.allWords.newIndics = []int{}

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

	c.chunkCache = append(c.chunkCache, &Chunk{
		Words:   words,
		Context: text,
	})
}
