package output

import "github.com/fatih/color"

var (
	colorNew      = color.New(color.FgGreen)
	colorLearning = color.New(color.FgYellow)
	colorFamiliar = color.New(color.FgCyan)
)

func RenderWordByProficiency(word string, proficiency float32) string {
	switch {
	case proficiency < 0.2:
		return colorNew.Sprint(word)
	case proficiency < 0.5:
		return colorLearning.Sprint(word)
	case proficiency < 0.8:
		return colorFamiliar.Sprint(word)
	}

	return word
}
