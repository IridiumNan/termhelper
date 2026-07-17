package models

import (
	"fmt"
	"math"
	"time"

	"github.com/fatih/color"
)

var (
	colorNew      = color.New(color.FgGreen)
	colorLearning = color.New(color.FgYellow)
	colorFamiliar = color.New(color.FgCyan)
)

const (
	minMinutes = 60.0
	factor     = 2.0
	scale      = 10.0
)

type WordEntry struct {
	Word string `storm:"id"`

	SimpleDefinition    string
	DetailedExplanation string

	Proficiency float32 `storm:"index"`

	NextReviewTime int64 `storm:"index"`
}

type WordResponse struct {
	Word string `json:"word"`

	SimpleDefinition    string `json:"simple_definition"`
	DetailedExplanation string `json:"detailed_explanation"`
}

type LLMExplainResults struct {
	Results []WordResponse `json:"results"`
}

func (word *WordEntry) NextReviewInterval() time.Duration {
	ratio := math.Pow(factor, float64(word.Proficiency)*scale)
	minutes := minMinutes * ratio

	return time.Duration(minutes) * time.Minute
}

func (word *WordEntry) UpdateReviewTime() {
	word.NextReviewTime = word.NextReviewTime + int64(word.NextReviewInterval())
}

func (word *WordEntry) ColorfulPrint() {
	if word == nil {
		return
	}
	var displayColor *color.Color

	switch {
	case word.Proficiency < 0.2:
		displayColor = colorNew
	case word.Proficiency < 0.5:
		displayColor = colorLearning
	case word.Proficiency < 0.8:
		displayColor = colorFamiliar
	}

	fmt.Println("word => ", word.Word)
	fmt.Println("simple_definition => ", word.SimpleDefinition)
	fmt.Println("detailed_explanation => ", word.DetailedExplanation)

	fmt.Println(displayColor.Sprint(word.Word))
	fmt.Println(displayColor.Sprint(word.SimpleDefinition))
	fmt.Println(displayColor.Sprint(word.DetailedExplanation))

	fmt.Println()
}
