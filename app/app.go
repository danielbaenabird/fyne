// Package app provides app implementations for working with Fyne graphical interfaces.
// The fastest way to get started is to call app.New() which will normally load a new desktop application.
// If the "ci" tag is passed to go (go run -tags ci myapp.go) it will run an in-memory application.
package app // import "github.com/daninemonic/fyne/v2/app"

import (
	"strconv"
	"sync"
	"time"

	"github.com/daninemonic/fyne/v2"
	"github.com/daninemonic/fyne/v2/internal"
	"github.com/daninemonic/fyne/v2/internal/app"
	intRepo "github.com/daninemonic/fyne/v2/internal/repository"
	"github.com/daninemonic/fyne/v2/storage/repository"

	"golang.org/x/sys/execabs"
)

// Declare conformity with App interface
var _ fyne.App = (*fyneApp)(nil)

type fyneApp struct {
	driver   fyne.Driver
	icon     fyne.Resource
	uniqueID string

	lifecycle fyne.Lifecycle
	settings  *settings
	storage   *store
	prefs     fyne.Preferences

	running  bool
	runMutex sync.Mutex
	exec     func(name string, arg ...string) *execabs.Cmd
}

func (a *fyneApp) Icon() fyne.Resource {
	return a.icon
}

func (a *fyneApp) SetIcon(icon fyne.Resource) {
	a.icon = icon
}

func (a *fyneApp) UniqueID() string {
	if a.uniqueID != "" {
		return a.uniqueID
	}

	fyne.LogError("Preferences API requires a unique ID, use app.NewWithID()", nil)
	a.uniqueID = "missing-id-" + strconv.FormatInt(time.Now().Unix(), 10) // This is a fake unique - it just has to not be reused...
	return a.uniqueID
}

func (a *fyneApp) NewWindow(title string) fyne.Window {
	return a.driver.CreateWindow(title)
}

func (a *fyneApp) Run() {
	a.runMutex.Lock()

	if a.running {
		a.runMutex.Unlock()
		return
	}

	a.running = true
	a.runMutex.Unlock()

	a.driver.Run()
}

func (a *fyneApp) Quit() {
	for _, window := range a.driver.AllWindows() {
		window.Close()
	}

	a.driver.Quit()
	a.settings.stopWatching()
	a.running = false
}

func (a *fyneApp) Driver() fyne.Driver {
	return a.driver
}

// Settings returns the application settings currently configured.
func (a *fyneApp) Settings() fyne.Settings {
	return a.settings
}

func (a *fyneApp) Storage() fyne.Storage {
	return a.storage
}

func (a *fyneApp) Preferences() fyne.Preferences {
	if a.uniqueID == "" {
		fyne.LogError("Preferences API requires a unique ID, use app.NewWithID()", nil)
	}
	return a.prefs
}

func (a *fyneApp) Lifecycle() fyne.Lifecycle {
	return a.lifecycle
}

// New returns a new application instance with the default driver and no unique ID
func New() fyne.App {
	internal.LogHint("Applications should be created with a unique ID using app.NewWithID()")
	return NewWithID("")
}

func newAppWithDriver(d fyne.Driver, id string) fyne.App {
	newApp := &fyneApp{uniqueID: id, driver: d, exec: execabs.Command, lifecycle: &app.Lifecycle{}}
	fyne.SetCurrentApp(newApp)
	newApp.settings = loadSettings()

	newApp.prefs = newPreferences(newApp)
	newApp.storage = &store{a: newApp}
	if id != "" {
		if pref, ok := newApp.prefs.(interface{ load() }); ok {
			pref.load()
		}

		root, _ := newApp.storage.docRootURI()
		newApp.storage.Docs = &internal.Docs{RootDocURI: root}
	} else {
		newApp.storage.Docs = &internal.Docs{} // an empty impl to avoid crashes
	}

	if !d.Device().IsMobile() {
		newApp.settings.watchSettings()
	}

	repository.Register("http", intRepo.NewHTTPRepository())
	repository.Register("https", intRepo.NewHTTPRepository())

	return newApp
}
