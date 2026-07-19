package storage

import (
	"time"

	"github.com/IridiumNan/termhelper/internal/models"
)

func (d *WordData) GetAllDueWordsCount() (count int, err error) {
	var words []*models.WordEntry
	err = d.db.Range(models.FieldNextReviewTimeStr, 0, time.Now().Unix(), &words)
	if err != nil {
		return -1, err
	}

	count = len(words)

	return
}
