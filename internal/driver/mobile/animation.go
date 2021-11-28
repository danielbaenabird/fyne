package mobile

import "github.com/daninemonic/fyne/v2"

func (d *mobileDriver) StartAnimation(a *fyne.Animation) {
	d.animation.Start(a)
}

func (d *mobileDriver) StopAnimation(a *fyne.Animation) {
	d.animation.Stop(a)
}
