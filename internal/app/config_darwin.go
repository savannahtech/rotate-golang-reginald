package app

import (
	"context"
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Config struct {
	MonitorDirectory string `mapstructure:"monitor_directory" validate:"required,dir"`
	CheckFrequency   int    `mapstructure:"check_frequency" validate:"required,min=1,max=60"`
	APIEndpoint      string `mapstructure:"api_endpoint" validate:"required,url"`
}

func (a *App) loadConfig(ctx context.Context) error {
	configPath, err := a.dialog.SaveFileDialog(ctx, runtime.SaveDialogOptions{
		DefaultFilename: "config.yaml",
		Title:           "Choose config file location",
	})
	if err != nil {
		return fmt.Errorf("failed to get config file path: %w", err)
	}

	a.logger.Printf("Config file path selected: %s", configPath)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		a.logger.Printf("Config file does not exist, creating default config at: %s", configPath)
		monitorDir, err := a.dialog.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{
			Title: "Select directory to monitor",
		})
		if err != nil {
			return fmt.Errorf("failed to select directory: %w", err)
		}
		if monitorDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			monitorDir = fmt.Sprintf("%s/Documents", homeDir)
			a.logger.Printf("No directory selected. Using default directory: %s", monitorDir)
		}

		if err := a.createDefaultConfig(configPath, monitorDir); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
	}

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&a.config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(a.config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	a.logger.Println("Config loaded successfully")
	return nil
}

func (a *App) createDefaultConfig(configPath string, monitorDir string) error {
	defaultConfig := fmt.Sprintf(`
monitor_directory: "%s"
check_frequency: 60
api_endpoint: "https://eo13t4hn4shbd6x.m.pipedream.net"
`, monitorDir)

	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}
