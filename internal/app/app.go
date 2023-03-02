package app

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// App struct
type App struct {
	ctx    context.Context
	log    *wails.CustomLogger
	Server *Server
	Menu   *menu.Menu
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{
		ctx:    nil,
		log:    logger.NewCustomLogger("APP"),
		Server: NewServer(),
		Menu:   &menu.Menu{},
	}
	app.initMenu()
	return app
}

// startup is called at application startup
func (a *App) Startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a *App) DomReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
