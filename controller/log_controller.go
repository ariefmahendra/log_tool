package controller

import (
	"Tools/config"
	"Tools/service"
	"errors"
	"fmt"
)

type LogController interface {
	DownloadLog(dir string) error
	SearchLog(dir, keyword string) error
	ListLogFile(parentDir string) error
	PrintLatestLog(dir string, bufferSize int) error
	CheckEnv()
}

type logControllerImpl struct {
	cfg        *config.Config
	logService service.LogService
}

func (l *logControllerImpl) CheckEnv() {
	envs := l.logService.CheckENV()
	for key, value := range envs {
		fmt.Printf("%s = %s\n", key, value)
	}
}

func (l *logControllerImpl) PrintLatestLog(dir string, bufferSize int) error {
	log, err := l.logService.PrintLatestLog(dir, bufferSize)
	if err != nil {
		return err
	}
	fmt.Println(log)
	return nil
}

func (l *logControllerImpl) ListLogFile(parentDir string) error {
	if err := l.logService.ListLogFile(parentDir); err != nil {
		return err
	}
	return nil
}

func (l *logControllerImpl) DownloadLog(dir string) error {
	err := l.logService.DownloadLog(dir)
	if err != nil {
		return errors.New("Failed to download log because " + err.Error())
	}
	return nil
}

func (l *logControllerImpl) SearchLog(dir, keyword string) error {
	if err := l.logService.SearchLog(dir, keyword); err != nil {
		return errors.New("Failed to search log because " + err.Error())
	}
	return nil
}

func NewLogController(cfg *config.Config, logService service.LogService) LogController {
	return &logControllerImpl{cfg: cfg, logService: logService}
}
