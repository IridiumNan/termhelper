package extractor

import (
	"slices"
	"time"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/internal/storage"
	"github.com/asdine/storm/v3"
)

type WordProvider interface {
	provide(word string) (entry *models.WordEntry, found bool)
}

type DisabledProvider struct {
	disabledWords []string
}

func newDisabledProvider(disabledWords []string) *DisabledProvider {
	return &DisabledProvider{
		disabledWords: disabledWords,
	}
}

func (dp *DisabledProvider) provide(word string) (entry *models.WordEntry, found bool) {
	if slices.Contains(dp.disabledWords, word) {
		return &models.WordEntry{
			Word:                word,
			SimpleDefinition:    models.EmptyStr,
			DetailedExplanation: models.EmptyStr,
			Proficiency:         disabledProficiency,
			NextReviewTime:      time.Now().Unix(),
		}, true
	}

	return nil, false
}

type DicProvider struct {
	wordData *storage.WordData
}

func newDicProvider() *DicProvider {
	return &DicProvider{
		wordData: storage.GetGlobalWordData(),
	}
}

func (dp *DicProvider) provide(word string) (entry *models.WordEntry, found bool) {
	// if the wordData init failed, wordData will be nil
	if dp.wordData == nil {
		return nil, false
	}

	// default found, if not found, reset it as false
	found = true
	var err error

	entry, err = dp.wordData.GetWordEntryByStr(word)

	if err == storm.ErrNotFound {
		found = false
	}

	return
}

type defaultProvider struct{}

func newDefaultProvider() *defaultProvider {
	return &defaultProvider{}
}

func (np *defaultProvider) provide(word string) (entry *models.WordEntry, found bool) {
	return &models.WordEntry{
		Word:                word,
		SimpleDefinition:    models.EmptyStr,
		DetailedExplanation: models.EmptyStr,
		Proficiency:         models.InitProficiency,
		NextReviewTime:      time.Now().Unix(),
	}, false
}
