package app

import (
	"strings"

	"github.com/valery-barysok/resp"
)

type ErrorHandler func(context *ClientContext, err error, res *resp.Writer)

type Filter func(context *ClientContext, req *RespCommand, res *resp.Writer) (bool, error)
type Handler func(context *ClientContext, req *RespCommand, res *resp.Writer) error

type router struct {
	filters []Filter
	routes  map[string]Handler

	notFound     Handler
	errorHandler ErrorHandler
}

func NewRouter() *router {
	return &router{
		filters: make([]Filter, 0),
		routes:  make(map[string]Handler),
	}
}

// BindFilter binds precondition filter
func (router *router) BindFilter(filter Filter) {
	router.filters = append(router.filters, filter)
}

func (router *router) Bind(cmd string, handler Handler) Handler {
	cmd = strings.ToLower(cmd)
	oldHandler := router.routes[cmd]
	router.routes[cmd] = handler
	return oldHandler
}

func (router *router) BindNotFound(handler Handler) Handler {
	oldHandler := router.notFound
	router.notFound = handler
	return oldHandler
}

func (router *router) BindError(errorHandler ErrorHandler) ErrorHandler {
	oldHandler := router.errorHandler
	router.errorHandler = errorHandler
	return oldHandler
}

func (router *router) Serve(context *ClientContext, req *RespCommand, res *resp.Writer) error {
	for _, filter := range router.filters {
		done, err := filter(context, req, res)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}

	return router.serve(context, req, res)
}

func (router *router) serve(context *ClientContext, req *RespCommand, res *resp.Writer) error {
	cmd := req.Cmd

	if handle := router.routes[cmd]; handle != nil {
		return handle(context, req, res)
	} else if router.notFound != nil {
		return router.notFound(context, req, res)
	}

	return nil
}
