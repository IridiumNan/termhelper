# Extractor

## 依赖的核心结构体

```go
type WordStore struct {
    entries []WordEntry
    indexMap map[string]int
    isNewMap map[string]bool
    newIndices []int
}

type Chunker struct {
    rawText string

    maxWordCount int
    maxCharCount int

    allWords WordStore

    currPos int
    contextStart int

    providers []WordProvider

    currChunkPos int
    chunkCache []*Chunk
}

type WordProvider interface {
    Provide (word string) (entry *WordEntry, found bool)
}

type DisabledProvider struct {
    disabled []string 
}

type BoltProvider struct {
    repo WordRepository
}

type NewProvider struct {}


type Chunk struct {
    Words []string
    Context string
}
```

## 核心功能

- NewChunker() -> 返回一个迭代器， 用于将长文本切成适合发送的片段

- NextChunk() -> 返回下一个能够直接操作的片段(chunkCache[currChunkPos]) 如果已经没有任何缓存则自动调用makeCache

- makeCache() -> 调用3次makeChunk, 接受递归拆分， 可能扩展出更多的Chunk, 等待大模型回复每一次Chunk的时候使用goroutine调用该方法

- makeChunk() -> 移动currPos， 将所有的单词去除首位的特殊字符后经过 Provider处理加入 allWords， 如果found = false, 将这个词加入到 当前的Chunk当中，当达到了 指定的单词个数之后， 自动寻找下一个语义分隔符，然后检查 Chunk当中的单词个数是否超过最大限制， 如果超过则折半递归处理
