package main

import (
	"github.com/valery-barysok/gredisd/app"
	"github.com/valery-barysok/gredisd/app/handlers"
)

func NewApp(opts *app.Options) *app.App {
	app := app.NewApp(opts)
	handlers.BindAllHanlders(app)
	return app
}
