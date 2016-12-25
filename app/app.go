package app

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"

	"github.com/valery-barysok/gredisd/app/cmd"
	"github.com/valery-barysok/gredisd/app/model"
	"github.com/valery-barysok/gredisd/server"
	"github.com/valery-barysok/resp"
)

type ErrorHandler func(context *ClientContext, err error, w *resp.Writer)

type Filter func(context *ClientContext, cmd *cmd.Command, w *resp.Writer) (bool, error)
type Handler func(context *ClientContext, cmd *cmd.Command, w *resp.Writer) error

type Info struct {
	ID        string `json:"server_id"`
	Version   string `json:"version"`
	GoVersion string `json:"go"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
}

type Options struct {
	Host          string `json:"addr"`
	Port          int    `json:"port"`
	Auth          string `json:"-"`
	Databases     int    `json:"databases"`
	TraceProtocol bool   `json:"trace_protocol"`
}

type App struct {
	info   Info
	opts   *Options
	server *server.Server
	router *router
	model  *model.AppModel
}

func NewApp(opts *Options) *App {
	normalizeOptions(opts)

	info := Info{
		ID:        genID(),
		Version:   Version,
		GoVersion: runtime.Version(),
		Host:      opts.Host,
		Port:      opts.Port,
	}
	app := &App{
		info:   info,
		opts:   opts,
		router: newRouter(),
		model:  model.NewAppModel(opts.Databases),
	}

	return app
}

func genID() string {
	// TODO: uuid
	return strconv.FormatInt(rand.Int63(), 10)
}

func (app *App) ShowVersion() {
	fmt.Printf("GRedis version %s\n", app.info.Version)
	fmt.Printf("Go runtime version %s\n", app.info.GoVersion)
}

func (app *App) Run() error {
	if app.server != nil {
		return errors.New("Already started")
	}

	log.Printf("GRedis version %s\n", app.info.Version)

	opts := server.Options{
		Host: app.opts.Host,
		Port: app.opts.Port,
	}

	if app.RequireAuth() {
		log.Println("App requires authentication")
	}

	//app.server = server.NewServer(&opts, NewClientProvider(app))
	app.server = server.NewServer(&opts, NewClientProvider(app))
	app.server.Start()

	return nil
}

func (app *App) Shutdown() {
	go func() {
		app.server.Shutdown()
		os.Exit(0)
	}()
}

func (app *App) Commands() []interface{} {
	return app.model.Commands()
}

func (app *App) BindFilter(filter Filter) {
	app.router.bindFilter(filter)
}

func (app *App) Bind(cmd string, handler Handler) Handler {
	app.model.AddCmd(cmd)
	return app.router.bind(cmd, handler)
}

func (app *App) BindNotFound(handler Handler) Handler {
	return app.router.bindNotFound(handler)
}

func (app *App) BindError(errorHandler ErrorHandler) ErrorHandler {
	return app.router.bindError(errorHandler)
}

func (app *App) RequireAuth() bool {
	return len(app.opts.Auth) > 0
}

func (app *App) Auth(pass string) bool {
	return !app.RequireAuth() || app.opts.Auth == pass
}

func (app *App) Select(index string) (*model.DBModel, error) {
	return app.model.Select(index)
}

func (app *App) SelectIndex(index int) (*model.DBModel, error) {
	return app.model.SelectIndex(index)
}

func normalizeOptions(opts *Options) {
	if opts.Host == "" {
		opts.Host = DefaultHost
	}
	if opts.Port == 0 {
		opts.Port = DefaultPort
	}
	if opts.Databases <= 0 {
		opts.Databases = DefaultDatabases
	}
}
