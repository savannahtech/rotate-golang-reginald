//go:build windows
// +build windows

package tray

import (
	"context"
	_ "embed"
	"os"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed logo-universal.png

var wailsIcon []byte

func CreateSystemTray(ctx context.Context) func() {

	return func() {
		systray.SetIcon(wailsIcon)
		show := systray.AddMenuItem("Show", "Show The Window")
		systray.AddSeparator()
		exit := systray.AddMenuItem("Exit", "Quit The Program")
		show.Click(func() {
			runtime.WindowShow(ctx)
		})

		exit.Click(func() {
			os.Exit(0)
		})
		systray.SetOnClick(func(menu systray.IMenu) {
			runtime.WindowShow(ctx)
			menu.ShowMenu()
		})
	}
}
