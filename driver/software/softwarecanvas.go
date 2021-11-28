package software

import (
	"github.com/daninemonic/fyne/v2/internal/painter/software"
	"github.com/daninemonic/fyne/v2/test"
)

// NewCanvas creates a new canvas in memory that can render without hardware support
func NewCanvas() test.WindowlessCanvas {
	return test.NewCanvasWithPainter(software.NewPainter())
}
