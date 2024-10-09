package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/osquery/osquery-go"
)

type Monitor struct {
	OsqueryInstance   *osquery.ExtensionManagerServer
	OsquerySocketPath string
	MonitorDirectory  string
}

type SystemStats struct {
	CPUUsage     string `json:"cpu_usage"`
	MemoryUsage  string `json:"memory_usage"`
	DiskUsage    string `json:"disk_usage"`
	SystemUptime string `json:"system_uptime"`
}

func (a *Monitor) GetSystemMonitoringData() (string, error) {
	if a.OsqueryInstance == nil {
		return "", fmt.Errorf("osquery instance not initialized")
	}

	client, err := osquery.NewClient(a.OsquerySocketPath, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to create osquery client: %w", err)
	}
	defer client.Close()

	// Query CPU usage
	cpuQuery := "SELECT cpu_time_user + cpu_time_system AS cpu_usage FROM cpu_time"
	cpuResponse, err := client.Query(cpuQuery)
	if err != nil {
		return "", fmt.Errorf("failed to query CPU usage: %w", err)
	}

	cpuUsage := "Unknown"
	if len(cpuResponse.Response) > 0 {
		cpuUsage = cpuResponse.Response[0]["cpu_usage"]
	}

	memoryQuery := "SELECT (total_available_bytes * 100.0) / total_bytes AS memory_usage FROM memory_info"
	memoryResponse, err := client.Query(memoryQuery)
	if err != nil {
		return "", fmt.Errorf("failed to query memory usage: %w", err)
	}

	memoryUsage := "Unknown"
	if len(memoryResponse.Response) > 0 {
		memoryUsage = memoryResponse.Response[0]["memory_usage"]
	}

	diskQuery := fmt.Sprintf("SELECT (blocks_available * 100.0) / blocks_size AS disk_usage FROM mounts WHERE path = '%s'", a.MonitorDirectory)
	diskResponse, err := client.Query(diskQuery)
	if err != nil {
		return "", fmt.Errorf("failed to query disk usage: %w", err)
	}

	diskUsage := "Unknown"
	if len(diskResponse.Response) > 0 {
		diskUsage = diskResponse.Response[0]["disk_usage"]
	}

	uptimeQuery := "SELECT total_seconds AS system_uptime FROM uptime"
	uptimeResponse, err := client.Query(uptimeQuery)
	if err != nil {
		return "", fmt.Errorf("failed to query system uptime: %w", err)
	}

	systemUptime := "Unknown"
	if len(uptimeResponse.Response) > 0 {
		systemUptime = uptimeResponse.Response[0]["system_uptime"]
	}

	systemStats := SystemStats{
		CPUUsage:     cpuUsage,
		MemoryUsage:  memoryUsage,
		DiskUsage:    diskUsage,
		SystemUptime: systemUptime,
	}

	jsonData, err := json.MarshalIndent(systemStats, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal system stats to JSON: %w", err)
	}

	log.Println("Updated system stats")
	return string(jsonData), nil
}
