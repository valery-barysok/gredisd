package handlers

import "github.com/valery-barysok/gredisd/app"

func BindAllHanlders(app *app.App) {
	BindAllBasicHandlers(app)
	BindAllKVHandlers(app)
	BindAllKVListHandlers(app)
	BindAllKVDictHandlers(app)
}
