package app

import (
	"strings"

	"github.com/valery-barysok/gredisd/app/cmd"
	"github.com/valery-barysok/resp"
)

type router struct {
	filters []Filter
	routes  map[string]Handler

	notFound     Handler
	errorHandler ErrorHandler
}

func newRouter() *router {
	return &router{
		filters: make([]Filter, 0),
		routes:  make(map[string]Handler),
	}
}

// BindFilter binds precondition filter
func (router *router) bindFilter(filter Filter) {
	router.filters = append(router.filters, filter)
}

func (router *router) bind(cmd string, handler Handler) Handler {
	cmd = strings.ToLower(cmd)
	oldHandler := router.routes[cmd]
	router.routes[cmd] = handler
	return oldHandler
}

func (router *router) bindNotFound(handler Handler) Handler {
	oldHandler := router.notFound
	router.notFound = handler
	return oldHandler
}

func (router *router) bindError(errorHandler ErrorHandler) ErrorHandler {
	oldHandler := router.errorHandler
	router.errorHandler = errorHandler
	return oldHandler
}

func (router *router) serve(context *ClientContext, cmd *cmd.Command, res *resp.Writer) error {
	for _, filter := range router.filters {
		done, err := filter(context, cmd, res)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}

	return router.handle(context, cmd, res)
}

func (router *router) handle(context *ClientContext, cmd *cmd.Command, res *resp.Writer) error {
	if handle := router.routes[cmd.Cmd]; handle != nil {
		return handle(context, cmd, res)
	} else if router.notFound != nil {
		return router.notFound(context, cmd, res)
	}

	return nil
}
