package models

import (
	"os"
	"path/filepath"
)

const (
	AppName           = "termhelper"
	logFileName       = "log.txt"
	configFileRawName = "config"
	configFileType    = "toml"
)

var (
	dataDirPath = filepath.Join(".local", "share")
	homeDir, _  = os.UserHomeDir()
	logDirPath  = filepath.Join(".local", "state")
)
