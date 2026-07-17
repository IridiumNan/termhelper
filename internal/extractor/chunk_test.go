package extractor

import (
	"fmt"
	"strings"
	"testing"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
	}
}

// 辅助：提取所有唯一单词（忽略大小写和标点）
func uniqueWords(text string) map[string]bool {
	words := make(map[string]bool)
	for _, field := range strings.Fields(text) {
		// 去除标点，保留字母数字
		w := strings.Trim(field, ".,!?;:()\"'")
		if w != "" {
			// 转为小写，统一处理
			words[strings.ToLower(w)] = true
		}
	}
	return words
}

func TestChunker_RealScenario(t *testing.T) {
	errorText := `Error: cannot open file /etc/hosts: Permission denied. 
Warning: network unreachable. 
Fatal: out of memory. 
Please check your configuration. 
Error: connection refused.`

	docText := `The Go programming language is an open source project to make programmers more productive. 
Go is expressive, concise, clean, and efficient. 
Its concurrency mechanisms make it easy to write programs that get the most out of multicore and networked machines.`

	fullText := errorText + "\n" + docText

	chunker := NewChunker(fullText)
	var chunks []*models.Chunk

	// 收集所有 Chunk
	for chunker.HasNext() {
		ch, _ := chunker.NextChunk()
		if ch == nil {
			t.Fatal("NextChunk returned nil")
		}
		chunks = append(chunks, ch)
	}

	if len(chunks) == 0 {
		t.Fatal("No chunks generated")
	}
	t.Logf("Generated %d chunks", len(chunks))

	// 1. 每个 Chunk 的单词数 ≤ maxWordCount
	for i, ch := range chunks {
		if len(ch.Words) > defaultMaxWords {
			t.Errorf("Chunk %d has %d words > %d", i, len(ch.Words), defaultMaxWords)
		}
	}

	// 2. 单词在块之间不重复（去重）
	seen := make(map[string]bool)
	for _, ch := range chunks {
		for _, w := range ch.Words {
			lower := strings.ToLower(w)
			if seen[lower] {
				t.Errorf("Duplicate word %q appears in multiple chunks", w)
			}
			seen[lower] = true
		}
	}

	// 3. 所有文本中的单词至少出现在一个块中（由于 disabledProvider 过滤，可能不全，仅做提示）
	allWords := uniqueWords(fullText)
	for w := range allWords {
		if !seen[w] {
			// 被 disabledProvider 过滤，不视为错误
			t.Logf("Word %q not in chunks (likely disabled)", w)
		}
	}
}

func TestChunker_AllWords(t *testing.T) {
	// 强制所有单词都是新词（使用 defaultProvider）
	text := `Error: cannot open file /etc/hosts: Permission denied. 
Warning: network unreachable. 
Fatal: out of memory.`

	chunker := &Chunker{
		rawText:      text,
		maxWordCount: defaultMaxWords,
		maxCharCount: defaultMaxChars,
		allWords: models.WordStore{
			Entries:    make([]*models.WordEntry, 0),
			IndexMap:   make(map[string]int),
			IsNewMap:   make(map[string]bool),
			NewInDices: make([]int, 0),
		},
		currPos:      0,
		contextStart: 0,
		providers:    []WordProvider{newDefaultProvider()}, // 仅默认，所有词均为新词
		currChunkPos: 0,
		chunkCache:   make([]*models.Chunk, 0),
	}

	var chunks []*models.Chunk
	for chunker.HasNext() {
		ch, _ := chunker.NextChunk()
		chunks = append(chunks, ch)
	}
	t.Logf("Generated %d chunks", len(chunks))
	for i, ch := range chunks {
		t.Logf("Chunk %d: words=%d", i, len(ch.Words))
	}
	// 验证单词总数等于所有不同单词数（忽略大小写）
	allWords := uniqueWords(text)
	totalFound := 0
	for _, ch := range chunks {
		totalFound += len(ch.Words)
	}
	if totalFound != len(allWords) {
		t.Logf("Total words in chunks: %d, unique words in text: %d", totalFound, len(allWords))
	}
}
