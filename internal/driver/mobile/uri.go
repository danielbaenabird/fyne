//go:build !android
// +build !android

package mobile

import (
	"github.com/daninemonic/fyne/v2"
	"github.com/daninemonic/fyne/v2/storage"
)

func nativeURI(uri string) fyne.URI {
	return storage.NewURI(uri)
}
