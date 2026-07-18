package models

import (
	"os"
	"path/filepath"
)

func GetHomeDirPath() string {
	return homeDir
}

func GetDataDirPath() string {
	return filepath.Join(homeDir, dataDirPath, AppName)
}

// return the whole path with Dir/FileName
func GetDataFilePath() string {
	return filepath.Join(GetDataDirPath(), dataFileName)
}

func GetConfigDirPath() string {
	userConfigPath, _ := os.UserConfigDir()
	return filepath.Join(userConfigPath, AppName)
}

func GetLogDirPath() string {
	return filepath.Join(homeDir, logDirPath, AppName)
}

func GetLogFilePath() string {
	return filepath.Join(GetLogDirPath(), logFileName)
}

func GetConfigRawName() string {
	return configFileRawName
}

func GetConfigType() string {
	return configFileType
}

func getConfigFullName() string {
	return configFileRawName + "." + configFileType
}

func GetConfigFilePath() string {
	return filepath.Join(GetConfigDirPath(), getConfigFullName())
}
