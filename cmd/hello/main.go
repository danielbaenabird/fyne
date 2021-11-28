// Package main loads a very basic Hello World graphical application.
package main

import (
	"github.com/daninemonic/fyne/v2/app"
	"github.com/daninemonic/fyne/v2/container"
	"github.com/daninemonic/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
}
