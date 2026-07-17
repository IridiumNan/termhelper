package models

type WordEntry struct {
	Word string

	SimpleDefinition    string
	DetailedExplanation string

	Proficiency float32
}

type WordResponse struct {
	Word string `json:"word"`

	SimpleDefinition    string `json:"simple_definition"`
	DetailedExplanation string `json:"detail_explain"`
}

type LLMExplainResults struct {
	Results []WordResponse
}
