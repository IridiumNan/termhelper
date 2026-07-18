package storage

import (
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/internal/output"
)

func (d *WordData) PushAllWordEntriesIfNew(words []*models.WordEntry) {
	for _, word := range words {
		d.PushSingleWordEntryIfNew(word)
	}
}

func (d *WordData) PushSingleWordEntryIfNew(word *models.WordEntry) {
	if word.Proficiency != models.InitProficiency {
		return
	}

	err := d.db.Save(word)
	if err != nil {
		output.File.Error("error when save word", "word", word.Word, "error", err)
	}
}

// GetWordEntryByStr : query for a word by wordStr
func (d *WordData) GetWordEntryByStr(wordStr string) (*models.WordEntry, error) {
	var word models.WordEntry
	err := d.db.One(models.FieldWordStr, wordStr, &word)
	if err != nil {
		return nil, err
	}

	return &word, nil
}

func (d *WordData) GetAllWordEntries() ([]*models.WordEntry, error) {
	var words []*models.WordEntry

	err := d.db.All(&words)

	return words, err
}

func (d *WordData) UpdateAll(words []*models.WordEntry) {
	for _, word := range words {
		err := d.db.Update(word)
		if err != nil {
			output.File.Error("error when update word entry on database", "error", err)
		}
	}
}
