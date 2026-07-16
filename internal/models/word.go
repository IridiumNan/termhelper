package models

type Word struct {
	Word       string    `json:"word"`
	Definition string    `json:"definition"`
	Examples   []Example `json:"examples"`
}

type Example struct {
	Sentence string `json:"sentence"`
	Context  string `json:"context"`
}

type ExtractedWord struct {
	Word     string `json:"word"`
	Context  string `json:"context"`
	Sentence string `json:"sentence"`
}
