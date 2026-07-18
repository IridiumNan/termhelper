package storage

import (
	"time"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/asdine/storm/v3"
)

const maxWordLimit = 100

func (d *WordData) GetDueWords(wordLimit int) (words []*models.WordEntry, err error) {
	if wordLimit < 0 {
		wordLimit = maxWordLimit
	}
	err = d.db.Range(models.FieldNextReviewTimeStr, 0, time.Now().Unix(), &words, storm.Limit(wordLimit))
	if err != nil {
		return nil, err
	}

	return
}
