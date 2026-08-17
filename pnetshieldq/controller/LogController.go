package controller

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"pnetshieldq/entity"
	"strings"
)

type LogController struct {
	LogFilePath string
}

func NewLogController(logFilePath string) *LogController {
	return &LogController{LogFilePath: logFilePath}
}

func ParseLog(file string) ([]entity.LogEntry, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	entries := make([]entity.LogEntry, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		entry := entity.LogEntry{}
		timeStart := strings.Index(line, `time="`) + len(`time="`)
		timeEnd := strings.Index(line[timeStart:], `"`) + timeStart
		if timeStart == -1 || timeEnd == -1 {
			continue
		}
		entry.Time = line[timeStart:timeEnd]

		levelStart := strings.Index(line, `level=`) + len(`level=`)
		levelEnd := strings.Index(line[levelStart:], " ")
		if levelStart == -1 || levelEnd == -1 {
			continue
		}
		entry.Level = line[levelStart : levelStart+levelEnd]

		msgStart := strings.Index(line, `msg="`) + len(`msg="`)
		msgEnd := strings.Index(line[msgStart:], `"`) + msgStart
		if msgStart == -1 || msgEnd == -1 {
			continue
		}
		entry.Message = line[msgStart:msgEnd]

		entries = append(entries, entry)
	}

	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	return entries, nil
}
func (lc *LogController) HandleLog(w http.ResponseWriter, r *http.Request) {
	entries, err := ParseLog("./log/app.log")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(entries)
}

func CountLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}
func (lc *LogController) HandleLogCount(w http.ResponseWriter, r *http.Request) {
	logFile := "./log/app.log"

	lineCount, err := CountLines(logFile)
	if err != nil {
		http.Error(w, "Failed to count lines in log file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = ParseLog(logFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		LineCount int `json:"line_count"`
	}{
		LineCount: lineCount,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
