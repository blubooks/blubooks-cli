package app

import mainApp "github.com/blubooks/blubooks-cli/internal/app"

type App struct {
	mainApp *mainApp.App
}

func New(app *mainApp.App) *App {
	return &App{
		mainApp: app,
	}
}
