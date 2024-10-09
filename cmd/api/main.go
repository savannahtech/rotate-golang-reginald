package main

import (
	"daemon/internal/app"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

type config struct {
	port int
	env  string
}

type serverApplication struct {
	config config
	logger *log.Logger
	logDir string
	app    *app.App
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := app.NewApp()
	srvApp := &serverApplication{
		config: cfg,
		logger: logger,
		app:    app,
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      srvApp.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.Server = srv

	err := wails.Run(&options.App{
		Title:             "daemon",
		Width:             640,
		HideWindowOnClose: true,
		Height:            480,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

}
