package models

type Chunk struct {
	Words   []string
	Context string
}

type WordStore struct {
	Entries    []*WordEntry
	IndexMap   map[string]int
	IsNewMap   map[string]bool
	NewInDices []int
}
