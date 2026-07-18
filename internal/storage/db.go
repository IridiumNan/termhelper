package storage

import (
	"fmt"
	"os"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/internal/output"
	"github.com/asdine/storm/v3"
)

var globalWordData *WordData

type WordData struct {
	db *storm.DB
}

func (d *WordData) Close() (err error) {
	err = d.db.Close()
	return
}

func GetGlobalWordData() *WordData {
	var err error
	if globalWordData == nil {
		globalWordData, err = newWordData()
	}

	if err != nil {
		output.Console.Error("error when init global word data", "error", err)
		output.File.Error("error when init global word data", "error", err)
		return nil
	}

	return globalWordData
}

func newWordData() (*WordData, error) {
	err := initDBFile()
	if err != nil {
		return nil, err
	}

	db, err := storm.Open(models.GetDataFilePath())
	if err != nil {
		return nil, err
	}

	return &WordData{
		db: db,
	}, err
}

func initDBFile() error {
	dbPath := models.GetDataFilePath()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		_, err = os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("fail ot create data file -> %s, err -> %w", dbPath, err)
		}
		output.File.Info("data file not found create a new one", "file_path", dbPath)

		db, err := storm.Open(dbPath)
		if err != nil {
			return fmt.Errorf("file to Open data file -> %s, err -> %w", dbPath, err)
		}

		err = db.Init(&models.WordEntry{})
		if err != nil {
			output.File.Error("fail to init the file with struct WordEntry", "err", err)

			return fmt.Errorf("fail to init the data file with struct WordEntry, err -> %w", err)
		}

		err = db.Close()
		if err != nil {
			output.File.Error("file to close the db", "error", err)

			return fmt.Errorf("fail to close the db, err -> %w", err)
		}
	}

	return nil
}
