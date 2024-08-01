package service

import (
	"Tools/config"
	"Tools/entity"
	"Tools/network"
	"Tools/util"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type LogService interface {
	ProcessLog(log []byte) ([]string, error)
	ExtractJson(log string) (string, string)
	CheckENV() map[string]string
	PrintLatestLog(dir string, bufferSize int) (string, error)
	ListLogFile(parentDir string) error
	DownloadLog(dir string) error
	SearchLog(dir, keyword string) error
}

type LogServiceImpl struct {
	cfg *config.Config
}

func (l *LogServiceImpl) SearchLog(dir, keyword string) error {
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
	if start < 0 {
		start = 0
	}

	var bytesRead int
	if start == 0 {
		bytesRead, err = file.Read(logBuffer)
	} else {
		bytesRead, err = file.ReadAt(logBuffer, start)
	}

	if err != nil && err != io.EOF {
		return errors.New("file read error")
	}

	log, err := l.ProcessLog(logBuffer[:bytesRead])
	if err != nil {
		return err
	}

	if len(log) == 0 {
		fmt.Println("No log found")
	}

	for _, strLog := range log {
		if strings.Contains(strings.ToLower(strLog), strings.ToLower(keyword)) {
			fmt.Println("====================================================================================")
			payloadInfo, payload := l.ExtractJson(strLog)
			fmt.Println(payloadInfo)
			fmt.Println(payload)
		}
	}
	return nil
}

func (l *LogServiceImpl) DownloadLog(dir string) error {
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

func (l *LogServiceImpl) ListLogFile(parentDir string) error {
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

func (l *LogServiceImpl) PrintLatestLog(dir string, bufferSize int) (string, error) {
	ssh, sftp, err := network.SetupNetwork(l.cfg)
	if err != nil {
		return "", err
	}
	defer ssh.Close()
	defer sftp.Close()

	file, err := sftp.Open(dir)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var size = 1000 * 1024 * (int64)(bufferSize)

	logBuffer := make([]byte, size)
	stat, statErr := file.Stat()
	if statErr != nil {
		return "", errors.New("file stat error")
	}

	start := stat.Size() - size
	if start < 0 {
		start = 0
	}

	var bytesRead int
	if start == 0 {
		bytesRead, err = file.Read(logBuffer)
	} else {
		bytesRead, err = file.ReadAt(logBuffer, start)
	}
	if err != nil && err != io.EOF {
		return "", errors.New("file read error")
	}

	return string(logBuffer[:bytesRead]), nil
}

func (l *LogServiceImpl) CheckENV() map[string]string {
	env := make(map[string]string)
	env["Buffer Size"] = l.cfg.BufferSize
	env["Default Path"] = l.cfg.DefaultFolder
	env["SSH User"] = l.cfg.Username
	env["SSH Password"] = l.cfg.Password
	env["SSH Host"] = l.cfg.Host
	env["SSH Port"] = l.cfg.Port
	return env
}

func (l *LogServiceImpl) ProcessLog(log []byte) ([]string, error) {
	var filteredLogs []string
	var currentLog string
	var capturing bool

	lines := strings.Split(string(log), "\n")
	for _, line := range lines {
		if capturing {
			if bytes.HasPrefix([]byte(line), []byte("ERROR")) || bytes.HasPrefix([]byte(line), []byte("DEBUG")) {
				currentLog += "🕒 TIME : " + strings.Split(line, " ")[1] + " " + strings.Split(line, " ")[2] + " 🕒\n"
				filteredLogs = append(filteredLogs, currentLog)
				currentLog = ""
				capturing = false
			} else {
				currentLog += line + "\n"
			}
		} else if bytes.HasPrefix([]byte(line), []byte("Type")) {
			currentLog += line + "\n"
			capturing = true
		}
	}

	if len(filteredLogs) == 0 {
		return nil, fmt.Errorf("no logs captured")
	}

	return filteredLogs, nil
}

func (l *LogServiceImpl) ExtractJson(log string) (string, string) {
	var payloadBuilder strings.Builder
	var payloadInfoBuilder strings.Builder
	var payloadInfoTimeBuilder strings.Builder

	isJson := false
	isTime := false

	for _, char := range log {
		if char == '{' {
			isJson = true
		}

		if char == '🕒' {
			payloadInfoBuilder.WriteString("\n")
			isJson = false
			isTime = true
		}

		if isJson {
			payloadBuilder.WriteRune(char)
		} else if !isTime {
			payloadInfoBuilder.WriteRune(char)
		} else {
			payloadInfoTimeBuilder.WriteRune(char)
		}

	}

	jsonString := payloadBuilder.String()
	PayloadInfoString := payloadInfoTimeBuilder.String() + payloadInfoBuilder.String()

	var payloadInfoTrim string
	if strings.HasSuffix(PayloadInfoString, "\n\n") {
		payloadInfoTrim = strings.TrimSuffix(PayloadInfoString, "\n\n")
	}

	if strings.HasSuffix(PayloadInfoString, "\n\n\n") {
		payloadInfoTrim = strings.TrimSuffix(PayloadInfoString, "\n\n\n")
	}
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, []byte(jsonString), "", "  ")
	if err != nil {
		return PayloadInfoString, jsonString
	}

	return payloadInfoTrim, prettyJson.String()
}

func NewLogService(cfg *config.Config) LogService {
	return &LogServiceImpl{cfg: cfg}
}
