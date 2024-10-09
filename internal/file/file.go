package file

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/osquery/osquery-go"
)

type FileInfo struct {
	Path         string `json:"path"`
	ModifiedTime string `json:"mtime"`
	Size         int64  `json:"size"`
}

type File struct {
	OsqueryInstance   *osquery.ExtensionManagerServer
	OsquerySocketPath string
	MonitorDirectory  string
	Mutex             sync.Mutex
	Logger            *log.Logger
}

func (a *File) GetFileModificationStats() (string, error) {
	if a.OsqueryInstance == nil {
		return "", fmt.Errorf("osquery instance not initialized")
	}

	client, err := osquery.NewClient(a.OsquerySocketPath, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to create osquery client: %w", err)
	}
	defer client.Close()

	query := fmt.Sprintf("SELECT path, mtime, size FROM file WHERE directory = '%s' ORDER BY mtime DESC", a.MonitorDirectory)
	response, err := client.Query(query)
	if err != nil {
		return "", fmt.Errorf("failed to execute osquery query: %w", err)
	}

	var files []FileInfo

	for _, r := range response.Response {
		mtimeUnix, err := strconv.ParseInt(r["mtime"], 10, 64)
		if err != nil {
			fmt.Println("Failed to parse mtime: ", err)
			continue
		}
		mtime := time.Unix(mtimeUnix, 0).Format(time.RFC3339)
		size, err := strconv.ParseInt(r["size"], 10, 64)
		if err != nil {
			fmt.Println("Failed to parse size: ", err)
			continue
		}
		files = append(files, FileInfo{
			Path:         r["path"],
			ModifiedTime: mtime,
			Size:         size,
		})
	}
	jsonData, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal file info to JSON: %w", err)
	}

	log.Println("Updated files stats")
	return string(jsonData), nil
}

func (a *File) GetLatestFileModifications() string {
	stats, err := a.GetFileModificationStats()
	if err != nil {
		return fmt.Sprintf("Error getting file modification stats: %v", err)
	}
	return stats
}

func (a *File) SaveStatsToFile(fileStats string, systemStats string) error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	logFilePath := filepath.Join(homeDir, "logs", "stats.log")
	err = os.MkdirAll(filepath.Dir(logFilePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open stats log file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("File Stats: %s\nSystem Stats: %s\n", fileStats, systemStats))
	if err != nil {
		return fmt.Errorf("failed to write stats to file: %w", err)
	}

	a.Logger.Println("Successfully saved stats to file at", logFilePath)
	return nil
}
