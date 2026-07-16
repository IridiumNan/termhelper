package models

const (
	EmptyPrompt = ""
)

type LLMRequest struct {
	Word    []string `json:"word"`
	Context string   `json:"context"`
}

type LLMExplainResponse struct {
	Word                string `json:"word"`
	SimpleDefinition    string `json:"simple_definition"`
	DetailedExplanation string `json:"detailed_explanation"`
}

type LLMExplainResults struct {
	Results []LLMExplainResponse `json:"results"`
}
