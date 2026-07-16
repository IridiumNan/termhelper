# 系统主要数据流思路

## 数据流示意图

用户输入文本或者文件 -> 根据单词和片段长度切分 -> 将经过 若干个 filter 过滤后的单词加入到 chuner.WordInChunk (如果没有经过过滤则标记为 true, 表示这个单词在单词库当中已经存在， 否则标记为 false, 每次增加一个false就把 newWordCount + 1) 中 -> 三种状况触发文本截断 -> 将截断之后的内容和对应的单词数组发送给大模型 (Chunk <-> LLMRequest) -> 大模型返回 一个数组， 其中的每一个元素都是 word -> simple_definition -> detail_explain -> 大模型的返回的内容加入到单词词库当中， 包含的元数据有 word -> simple_definition -> detail_explain -> proficiency -> 把整个单词列表中的词按照顺序拿出来， 如果熟练度超过 0.9 则直接跳过， 最终得到用户正在学习的词和完全没有见过的词 (当然在这里我们可以让用户自己维护一个配置文件， 直接跳过标记在配置文件里面的单词， 这样子用户就不需要每次都看到一堆冗长的输出了)， 并且这里应该做一个熟练度衰减机制， 比如说过十天熟练度自动衰减 0.1, 方便用户复习内容

## 关键数据结构

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
}

type WordProvider interface {
    Provide (word string) (entry *WordEntry, found bool)
}

type DisabledProvider struct {
    disabled []string 
}

type DicProvider struct {
    repo WordRepository
}

type NewProvider struct {}


type Chunk struct {
    Words []string
    Context string
}

type OllamaClient struct {
    URL        string
    Model      string
    OllamaOpts config.OllamaOptions
    http *http.Client
}
type Prompt struct {
    Role    string   `json:"role"`
    Task    string   `json:"task"`
    Format  string   `json:"format"`
    Example string   `json:"example"`
    Words   []string `json:"words"`
    Context string   `json:"context"`
}

type WordResponse struct {
    Word                string `json:"word"`
    SimpleDefinition    string `json:"simple_definition"`
    DetailedExplanation string `json:"detailed_explanation"`
}

type LLMExplainResults struct {
    Results []WordEntry `json:"results"`
}

type WordEntry struct {
    Word string
    SimpleDefinition string
    DetailedExplanation string

    Proficienty float32
}
```

完整数据流

rawText 原始文本 ->
经过chunker和 providers 进行切割， 获取到Chunk ->
发送Chunk给大模型 ->
大模型返回内容 ->
更新Chunker当中的allWords ->
所有内容处理完毕， 遍历allWords ->
使用将输出保存在 output_path 目录的新文件下， 并自动使用less打开->
遍历allWords之后将所有的数据刷入的词库当中
