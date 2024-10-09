package query

import (
	"context"
	"fmt"
	"os"

	"github.com/osquery/osquery-go"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Osquery struct {
	OsquerySocketPath string
	OsqueryInstance   *osquery.ExtensionManagerServer
	ctx               context.Context
}

func (a *Osquery) InitOsquery() error {
	namedPipePath, err := a.discoverOsqueryPipe()
	if err != nil {
		return fmt.Errorf("failed to discover osquery named pipe: %w", err)
	}

	a.OsquerySocketPath = namedPipePath

	server, err := osquery.NewExtensionManagerServer("file_monitor", namedPipePath)
	if err != nil {
		return fmt.Errorf("failed to create osquery extension: %w", err)
	}

	a.OsqueryInstance = server

	go func() {
		if err := server.Run(); err != nil {
			wailsRuntime.LogErrorf(a.ctx, "osquery extension server stopped: %v", err)
		}
	}()

	return nil
}

func (a *Osquery) discoverOsqueryPipe() (string, error) {
	namedPipe := `\\.\pipe\shell.em`

	if _, err := os.Stat(namedPipe); err == nil {
		return namedPipe, nil
	}

	return "", fmt.Errorf("could not find osquery named pipe at %s", namedPipe)
}

func (a *Osquery) GetOsqueryStatus() string {
	if a.OsquerySocketPath == "" {
		return "Osquery named pipe not set. Has initOsquery been called?"
	}

	if _, err := os.Stat(a.OsquerySocketPath); os.IsNotExist(err) {
		return fmt.Sprintf("Osquery named pipe not found at %s. Is osqueryd running?", a.OsquerySocketPath)
	}

	if a.OsqueryInstance == nil {
		return "Osquery named pipe found, but extension is not initialized"
	}

	return fmt.Sprintf("Osquery is running and extension is initialized. Pipe: %s", a.OsquerySocketPath)
}
