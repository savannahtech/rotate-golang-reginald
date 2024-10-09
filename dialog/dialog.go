package dialog

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Dialog interface {
	SaveFileDialog(ctx context.Context, options runtime.SaveDialogOptions) (string, error)
	OpenDirectoryDialog(ctx context.Context, options runtime.OpenDialogOptions) (string, error)
}

type WailsDialog struct{}

func (w *WailsDialog) SaveFileDialog(ctx context.Context, options runtime.SaveDialogOptions) (string, error) {
	return runtime.SaveFileDialog(ctx, options)
}

func (w *WailsDialog) OpenDirectoryDialog(ctx context.Context, options runtime.OpenDialogOptions) (string, error) {
	return runtime.OpenDirectoryDialog(ctx, options)
}
