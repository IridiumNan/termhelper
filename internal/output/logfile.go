package output

import (
	"log/slog"
	"os"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/pkg"
)

var File *slog.Logger

func getLogFile() (*os.File, error) {
	funcName := "GetLogFile"
	logDirPath := models.GetLogDirPath()

	if err := pkg.EnsureDir(logDirPath); err != nil {
		return nil, pkg.ReportErrorInputErr(funcName, err)
	}

	logFilePathh := models.GetLogFilePath()

	file, err := os.OpenFile(logFilePathh, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, pkg.ReportErrorInputErr(funcName, err)
	}

	return file, nil
}

func InitFileLogger(level slog.Level) error {
	funcName := "InitFileLogger"

	file, err := getLogFile()
	if err != nil {
		return pkg.ReportErrorInputErr(funcName, err)
	}

	fileOpts := &slog.HandlerOptions{
		Level: level,
	}

	File = slog.New(slog.NewJSONHandler(file, fileOpts))

	return nil
}
