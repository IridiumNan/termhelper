package models

import (
	"os"
	"path/filepath"
)

type ModelName string

const (
	EmptyStr          string    = ""
	AppName           string    = "termhelper"
	logFileName       string    = "log.txt"
	configFileRawName string    = "config"
	configFileType    string    = "toml"
	dataFileName      string    = "words.db"
	OllamaClientStr   ModelName = "ollama"
	OpenAIClientStr   ModelName = "openai"
	DeepSeekClientStr ModelName = "deepseek"

	InitProficiency float32 = 0.1
)

var (
	dataDirPath = filepath.Join(".local", "share")
	homeDir, _  = os.UserHomeDir()
	logDirPath  = filepath.Join(".local", "state")
)
