package models

type Chunk struct {
	Words   []string
	Context string
}

type WordStore struct {
	Entries []*WordEntry

	// map cache for O(1) to get word entry
	IndexMap map[string]int

	// If this word is new (need to be push to the boltdb), value will be true else false
	IsNewMap map[string]bool

	// just store for current new word index for Next Chunk, clear when this chunk is poped
	NewInDices []int
}
