package extractor

import (
	"slices"
	"time"

	"github.com/IridiumNan/termhelper/internal/models"
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
	// TODO:
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
