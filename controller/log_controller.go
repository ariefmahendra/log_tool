package controller

import (
	"Tools/config"
	"Tools/entity"
	"Tools/network"
	"Tools/util"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type LogController interface {
	DownloadLog(dir string) error
	SearchLog(dir, keyword string) error
	ListLogFile(parentDir string) error
	PrintLatestLog(dir string, bufferSize int) error
	CheckEnv()
}

type logControllerImpl struct {
	cfg *config.Config
}

func (l *logControllerImpl) CheckEnv() {
	fmt.Println("SFTP SERVER INFO")
	fmt.Println("HOST : ", l.cfg.FTPConfig.Host)
	fmt.Println("PORT : ", l.cfg.FTPConfig.Port)
	fmt.Println("USERNAME : ", l.cfg.FTPConfig.Username)
}

func (l *logControllerImpl) PrintLatestLog(dir string, bufferSize int) error {
	ssh, sftp, err := network.SetupNetwork(l.cfg)
	if err != nil {
		return err
	}
	defer ssh.Close()
	defer sftp.Close()

	file, err := sftp.Open(dir)
	if err != nil {
		return err
	}
	defer file.Close()

	var size int64 = 1000 * 1024 * (int64)(bufferSize)

	logBuffer := make([]byte, size)
	stat, statErr := file.Stat()
	if statErr != nil {
		return errors.New("file stat error")
	}

	start := stat.Size() - size
	_, err = file.ReadAt(logBuffer, start)
	if err == nil {
		fmt.Println(string(logBuffer))
	}

	return nil
}

func (l *logControllerImpl) ListLogFile(parentDir string) error {
	ssh, sftp, err := network.SetupNetwork(l.cfg)
	if err != nil {
		return err
	}

	defer ssh.Close()
	defer sftp.Close()

	var actualDirectory string
	DirSplited := strings.Split(parentDir, "/")
	for i := 0; i < len(DirSplited)-1; i++ {
		actualDirectory += DirSplited[i] + "/"
	}

	dir, err := sftp.ReadDir(actualDirectory)
	if err != nil {
		return err
	}

	var directories []entity.FileInfo
	var files []entity.FileInfo

	for _, file := range dir {
		if file.IsDir() {
			directories = append(directories, entity.FileInfo{file.Name(), file.ModTime()})
		} else {
			files = append(files, entity.FileInfo{file.Name(), file.ModTime()})
		}
	}

	sort.Slice(directories, func(i, j int) bool {
		return directories[i].ModTime.After(directories[j].ModTime)
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.After(files[j].ModTime)
	})

	if len(directories) != 0 {
		fmt.Println("Directories:")
		util.PrintSeparator(70)
		for _, dir := range directories {
			fmt.Printf("%-50s %s\n", dir.Name, dir.ModTime.Format("2006-01-02 15:04:05"))
		}
		util.PrintSeparator(70)
	}

	if len(files) != 0 {
		fmt.Println("\nFiles:")
		util.PrintSeparator(70)
		for _, file := range files {
			fmt.Printf("%-50s %s\n", file.Name, file.ModTime.Format("2006-01-02 15:04:05"))
		}
		util.PrintSeparator(70)
	}

	if len(directories) == 0 && len(files) == 0 {
		fmt.Println("EMPTY FILE OR DIRECTORY")
	}

	return nil
}

func (l *logControllerImpl) DownloadLog(dir string) error {
	ssh, sftp, err := network.SetupNetwork(l.cfg)
	if err != nil {
		return err
	}
	defer ssh.Close()
	defer sftp.Close()

	sftpFile, err := sftp.Open(dir)
	if err != nil {
		return err
	}
	defer sftpFile.Close()

	filestatus, err := sftpFile.Stat()
	if err != nil {
		return err
	}

	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	fileNameSlice := strings.Split(filestatus.Name(), ".")

	if len(fileNameSlice) == 2 {
		_, err = os.Stat(currentPath + "/logDownloaded/" + fileNameSlice[0] + "." + fileNameSlice[1])
		if os.IsNotExist(err) {
			var file, err = os.Create(currentPath + "/logDownloaded/" + fileNameSlice[0] + "." + fileNameSlice[1])
			if err != nil {
				os.MkdirAll(currentPath+"/logDownloaded", os.ModePerm)
				file, err = os.Create(currentPath + "/logDownloaded/" + fileNameSlice[0] + "." + fileNameSlice[1])
				if err != nil {
					return err
				}
			}
			defer file.Close()

			_, err = sftpFile.WriteTo(file)
			if err != nil {
				return err
			}

			fmt.Println("Successfully Downloaded: ", fileNameSlice[0]+"."+fileNameSlice[1])
		} else {
			fmt.Println("Are you sure to override file [y/n]: ", fileNameSlice[0]+"."+fileNameSlice[1])
			fmt.Print("$ ")
			var input string
			fmt.Scanln(&input)
			if input == "y" {
				var file, err = os.Create(currentPath + "/logDownloaded/" + fileNameSlice[0] + "." + fileNameSlice[1])

				defer file.Close()

				_, err = sftpFile.WriteTo(file)
				if err != nil {
					return err
				}

				fmt.Println("Successfully Downloaded: ", fileNameSlice[0]+"."+fileNameSlice[1])
			} else if input == "n" {
				fmt.Println("Download aborted: ", fileNameSlice[0]+"."+fileNameSlice[1])
			} else {
				fmt.Println("Input is not valid")
			}
		}
	}

	return nil
}

func (l *logControllerImpl) SearchLog(dir, keyword string) error {
	ssh, sftp, err := network.SetupNetwork(l.cfg)
	if err != nil {
		return err
	}

	defer ssh.Close()
	defer sftp.Close()

	file, err := sftp.Open(dir)
	if err != nil {
		return err
	}
	defer file.Close()

	bufferSize, err := strconv.Atoi(l.cfg.BufferSize)
	if err != nil {
		return err
	}

	var size int64 = 1000 * 1024 * (int64)(bufferSize)

	logBuffer := make([]byte, size)
	stat, statErr := file.Stat()
	if statErr != nil {
		return errors.New("file stat error")
	}

	start := stat.Size() - size
	_, err = file.ReadAt(logBuffer, start)
	if err == nil {
		log, err := util.ProcessLog(logBuffer)
		if err != nil {
			return err
		}

		if len(log) == 0 {
			fmt.Println("No log found")
		}

		for _, strLog := range log {
			if strings.Contains(strings.ToLower(strLog), strings.ToLower(keyword)) {
				fmt.Println("====================================================================================")
				payloadInfo, payload := util.ExtractJson(strLog)
				fmt.Print(payloadInfo)
				fmt.Println(payload)
			}
		}

	}

	return nil
}

func NewLogController(cfg *config.Config) LogController {
	return &logControllerImpl{cfg: cfg}
}
