//go:build ci
// +build ci

package app

import (
	"github.com/daninemonic/fyne/v2"
	"github.com/daninemonic/fyne/v2/internal/painter/software"
	"github.com/daninemonic/fyne/v2/test"
)

// NewWithID returns a new app instance using the test (headless) driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	return newAppWithDriver(test.NewDriverWithPainter(software.NewPainter()), id)
}
