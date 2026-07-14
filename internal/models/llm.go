package models

type LLMRequest struct {
	Wrod    string `json:"word"`
	Context string `json:"context"`
}

type LLMResponse struct {
	Word       string   `json:"word"`
	Definition []string `json:"definition"`
}
