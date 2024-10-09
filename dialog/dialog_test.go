package dialog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWailsDialogStruct(t *testing.T) {
	wailsDialog := &WailsDialog{}
	assert.NotNil(t, wailsDialog, "Expected WailsDialog to be created")
}
