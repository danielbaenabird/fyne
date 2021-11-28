//go:build !darwin || no_native_menus
// +build !darwin no_native_menus

package glfw

import "github.com/danielbaenabird/fyne/v2"

func hasNativeMenu() bool {
	return false
}

func setupNativeMenu(_ *window, _ *fyne.MainMenu) {
	// no-op
}
