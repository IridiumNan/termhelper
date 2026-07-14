package models

import (
	"os"
	"path/filepath"
)

type ModelName string

const (
	AppName                     = "termhelper"
	logFileName                 = "log.txt"
	configFileRawName           = "config"
	configFileType              = "toml"
	dataFileName                = "words.db"
	OllamaClientStr   ModelName = "ollama"
	OpenAIClientStr   ModelName = "openai"
	DeepSeekClientStr ModelName = "deepseek"
)

var (
	dataDirPath = filepath.Join(".local", "share")
	homeDir, _  = os.UserHomeDir()
	logDirPath  = filepath.Join(".local", "state")
)
