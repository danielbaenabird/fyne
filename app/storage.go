package app

import (
	"os"

	"github.com/daninemonic/fyne/v2"
	"github.com/daninemonic/fyne/v2/internal"
	"github.com/daninemonic/fyne/v2/storage"
)

type store struct {
	*internal.Docs
	a *fyneApp
}

func (s *store) RootURI() fyne.URI {
	if s.a.uniqueID == "" {
		fyne.LogError("Storage API requires a unique ID, use app.NewWithID()", nil)
		return storage.NewURI("file://" + os.TempDir())
	}

	return storage.NewURI("file://" + s.a.storageRoot())
}

func (s *store) docRootURI() (fyne.URI, error) {
	return storage.Child(s.RootURI(), "Documents")
}
