package app

import (
	"bytes"
	"context"
	"daemon/commands"
	"daemon/dialog"
	"daemon/internal/file"
	"daemon/internal/monitor"
	"daemon/internal/query"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

type App struct {
	ctx           context.Context
	config        Config
	logger        *log.Logger
	logBuffer     *bytes.Buffer
	osquery       query.Osquery
	WorkerQueue   chan string
	timerLogs     []string
	dialog        *dialog.WailsDialog
	mutex         sync.Mutex
	Server        *http.Server
	stopWorker    chan struct{}
	stopTimer     chan struct{}
	workerRunning bool
	timerRunning  bool
}

func NewApp() *App {
	logBuffer := new(bytes.Buffer)
	multiWriter := io.MultiWriter(os.Stdout, logBuffer)
	return &App{
		WorkerQueue: make(chan string, 100),
		timerLogs:   []string{},
		logBuffer:   logBuffer,
		logger:      log.New(multiWriter, "AppLogger: ", log.LstdFlags),
		stopWorker:  make(chan struct{}),
		stopTimer:   make(chan struct{}),
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	err := a.loadConfig(ctx)
	if err != nil {
		a.logger.Println("Could not load config:", err)
	} else {
		a.logger.Println("Config loaded:", a.config)
	}

	err = a.osquery.InitOsquery()
	if err != nil {
		a.logger.Println("Could not connect to osquery:", err)
	}

	a.logger.Printf("starting server on %s", a.Server.Addr)
	err = a.Server.ListenAndServe()
	a.logger.Fatal(err)

}

func (a *App) StartService() (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if !a.workerRunning {
		go a.workerThread()
		a.workerRunning = true
	}
	if !a.timerRunning {
		go a.timerThread()
		a.timerRunning = true
	}
	return "Service started", nil
}

func (a *App) FetchLogs() (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	logs := a.logBuffer.String()

	a.logBuffer.Reset()

	return logs, nil
}

func (a *App) StopService() (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.workerRunning {
		a.stopWorker <- struct{}{}
		a.workerRunning = false
	}
	if a.timerRunning {
		a.stopTimer <- struct{}{}
		a.timerRunning = false
	}
	return "Service stopped", nil
}

func (a *App) workerThread() {
	a.logger.Println("Worker thread started")
	var frequency int
	if a.config.CheckFrequency == 0 {
		frequency = 1
	} else {
		frequency = a.config.CheckFrequency
	}
	whitelist := commands.GetWhitelist()
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case cmdStr := <-a.WorkerQueue:
			a.logger.Printf("Received command: %s", cmdStr)

			if _, allowed := whitelist[cmdStr]; !allowed {
				a.logger.Printf("Command not allowed: %s", cmdStr)
				continue
			}
			cmd := exec.Command("sh", "-c", cmdStr)
			output, err := cmd.CombinedOutput()
			if err != nil {
				a.logger.Printf("Error executing command: %v, output: %s", err, output)
				continue
			}
			a.logger.Printf("Command output: %s", output)
		case <-a.stopWorker:
			a.logger.Println("Worker thread stopped")
			return
		}
	}
}

func (a *App) timerThread() {
	var frequency int
	m := monitor.Monitor{
		a.osquery.OsqueryInstance,
		a.osquery.OsquerySocketPath,
		a.config.MonitorDirectory,
	}

	f := file.File{
		a.osquery.OsqueryInstance,
		a.osquery.OsquerySocketPath,
		a.config.MonitorDirectory,
		a.mutex,
		a.logger,
	}

	if a.config.CheckFrequency == 0 {
		frequency = 1
	} else {
		frequency = a.config.CheckFrequency
	}
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
	defer ticker.Stop()
	a.logger.Println("Timer thread started")

	for {
		select {
		case <-ticker.C:
			stats, err := f.GetFileModificationStats()
			if err != nil {
				a.logger.Printf("Error getting file modification stats: %v", err)
				continue
			}

			systemStats, err := m.GetSystemMonitoringData()
			if err != nil {
				a.logger.Printf("Error getting system monitoring data: %v", err)
				continue
			}

			a.mutex.Lock()
			a.timerLogs = append(a.timerLogs, stats)
			a.mutex.Unlock()

			if err := f.SaveStatsToFile(stats, systemStats); err != nil {
				a.logger.Printf("Error saving stats to file: %v", err)
			}

			if err := a.sendStatsToAPI(stats, systemStats); err != nil {
				a.logger.Printf("Error sending stats to API: %v", err)
			}
		case <-a.stopTimer:
			a.logger.Println("Timer thread stopped")
			return
		}
	}
}

func (a *App) sendStatsToAPI(fileStats string, systemMonitor string) error {
	payload := map[string]string{
		"file_stats":     fileStats,
		"system_monitor": systemMonitor,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal stats to JSON: %w", err)
	}

	apiEndpoint := a.config.APIEndpoint

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send stats to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API responded with status: %s", resp.Status)
	}

	a.logger.Println("Successfully sent stats to API")
	return nil
}
