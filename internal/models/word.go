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

	correctStep = 0.1
	wrongStep   = 0.03

	hourSecond = 3600

	FieldWordStr                = "Word"
	FieldSimpleDefinitionStr    = "SimpleDefinition"
	FieldDetailedExplanationStr = "DetailedExplanation"
	FieldNextReviewTimeStr      = "NextReviewTime"
	FieldProficiency            = "Proficiency"

	ProficiencyNew      = 0.2
	ProficiencyLearning = 0.5
	ProficiencyFamiliar = 0.8
)

var LONGTIMELATER = time.Now().Add(100 * 365 * 24 * time.Hour).Unix()

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

func (word *WordEntry) UpdateProficiency(correct bool) {
	fmt.Print("update proficiency ", word.Proficiency)

	if correct {
		word.Proficiency += correctStep
		// word.Proficiency += config.GetGlobalConfig().Test.Proficiency.CorrectIncrement

		if word.Proficiency > 1.0 {
			word.Proficiency = 1.0
		}

		fmt.Print(color.GreenString(" => %f\n\n", word.Proficiency))
		return
	}

	word.Proficiency -= wrongStep
	// word.Proficiency -= config.GetGlobalConfig().Test.Proficiency.WrongDecrement

	if word.Proficiency < 0.0 {
		word.Proficiency = 0.0
	}

	fmt.Print(color.RedString(" => %f\n\n", word.Proficiency))
}

func (word *WordEntry) NextReviewInterval() int64 {
	ratio := math.Pow(factor, float64(word.Proficiency)*scale)
	minutes := minMinutes * ratio

	result := time.Duration(minutes) * time.Minute

	return int64(result.Seconds())
}

// MaskWithLongTime : update the NextReviewTime to an impossible reach time for mask this word
func (word *WordEntry) MaskWithLongTime() {
	word.Proficiency = 1.0
	word.NextReviewTime = LONGTIMELATER
}

func (word *WordEntry) UpdateReviewTime() {
	word.NextReviewTime = word.NextReviewTime + int64(word.NextReviewInterval()) - hourSecond
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
	// Use the proficiency to decide word and SimpleDefinition color
	fmt.Println(displayColor.Sprint(word.Word))
	fmt.Println(displayColor.Sprint(word.SimpleDefinition))
	fmt.Println(word.DetailedExplanation)

	fmt.Printf("word due timestamp; due %d, now %d\n", word.NextReviewTime, time.Now().Unix())

	fmt.Println()
}

func (word *WordEntry) FileFormat() string {
	return fmt.Sprintln() + color.YellowString("%s\n", word.Word) +
		color.HiGreenString("%s\n", word.SimpleDefinition) +
		color.CyanString("%s\n", word.DetailedExplanation) +
		fmt.Sprintln("Proficiency => ", word.Proficiency) +
		fmt.Sprintln("======================================================")
}
