package storage

import (
	"time"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/fatih/color"
)

func (d *WordData) getAllDueWordsCount() (count int, err error) {
	var words []*models.WordEntry
	err = d.db.Range(models.FieldNextReviewTimeStr, 0, time.Now().Unix(), &words)
	if err != nil {
		return -1, err
	}

	count = len(words)

	return
}

func (d *WordData) getCountByProficiency(min float32, max float32) (count int, err error) {
	var words []*models.WordEntry

	err = d.db.Range(models.FieldProficiency, min, max, &words)

	return len(words), err
}

func (d *WordData) getNewProficiencyCount() (count int, err error) {
	return d.getCountByProficiency(-0.1, models.ProficiencyNew)
}

func (d *WordData) getLearningProficiencyCount() (count int, err error) {
	return d.getCountByProficiency(models.ProficiencyNew, models.ProficiencyLearning)
}

func (d *WordData) getFamiliarProficiencyCount() (count int, err error) {
	return d.getCountByProficiency(models.ProficiencyLearning, models.ProficiencyFamiliar)
}

func (d *WordData) getMasterProficiencyCount() (count int, err error) {
	return d.getCountByProficiency(models.ProficiencyFamiliar, 1)
}

func (d *WordData) GetDueWordsStatus() string {
	out := "----------------- due words status ----------------"

	count, err := d.getAllDueWordsCount()

	if err != nil || count != 0 {
		out += color.RedString("\n\twords to be review count: %d\n", count)
	} else {
		out += color.RedString("\n\twords to be review count: %d\n", 0)
	}

	return out
}

// GetProficiencyStatus : return the string with three status word color
func (d *WordData) GetProficiencyStatus() string {
	out := "---------------- word proficiency status ----------------"

	count, err := d.getNewProficiencyCount()
	if err != nil || count >= 0 {
		out += color.GreenString("\n\twords new count: %d\n", count)
	} else {
		out += color.GreenString("\n\twords new count: %d\n", 0)
	}

	count, err = d.getLearningProficiencyCount()
	if err != nil || count != 0 {
		out += color.YellowString("\n\twords learning count: %d\n", count)
	} else {
		out += color.YellowString("\n\twords learning count: %d\n", 0)
	}

	count, err = d.getFamiliarProficiencyCount()
	if err != nil || count != 0 {
		out += color.BlueString("\n\twords familiar count: %d\n", count)
	} else {
		out += color.BlueString("\n\twords familiar count: %d\n", 0)
	}

	count, err = d.getMasterProficiencyCount()
	if err != nil || count != 0 {
		out += color.MagentaString("\n\twords master count: %d\n", count)
	} else {
		out += color.MagentaString("\n\twords master count: %d\n", 0)
	}

	return out
}
