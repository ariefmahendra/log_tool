package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func ProcessLog(log []byte) ([]string, error) {
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

func ExtractJson(log string) (string, string) {
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
