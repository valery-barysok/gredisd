package app

import (
	"strings"

	"github.com/valery-barysok/resp"
)

type ErrorHandler func(context *ClientContext, err error, res *resp.Writer)

type Filter func(context *ClientContext, req *RespRequest, res *resp.Writer) (bool, error)
type Handler func(context *ClientContext, req *RespRequest, res *resp.Writer) error

type Router struct {
	filters []Filter
	routes  map[string]Handler

	notFound     Handler
	errorHandler ErrorHandler
}

func NewRouter() *Router {
	return &Router{
		filters: make([]Filter, 0),
		routes:  make(map[string]Handler),
	}
}

// BindFilter binds precondition filter
func (router *Router) BindFilter(filter Filter) {
	router.filters = append(router.filters, filter)
}

func (router *Router) Bind(cmd string, handler Handler) Handler {
	cmd = strings.ToLower(cmd)
	oldHandler := router.routes[cmd]
	router.routes[cmd] = handler
	return oldHandler
}

func (router *Router) BindNotFound(handler Handler) Handler {
	oldHandler := router.notFound
	router.notFound = handler
	return oldHandler
}

func (router *Router) BindError(errorHandler ErrorHandler) ErrorHandler {
	oldHandler := router.errorHandler
	router.errorHandler = errorHandler
	return oldHandler
}

func (router *Router) Serve(context *ClientContext, req *RespRequest, res *resp.Writer) error {
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

func (router *Router) serve(context *ClientContext, req *RespRequest, res *resp.Writer) error {
	cmd := req.Cmd

	if handle := router.routes[cmd]; handle != nil {
		return handle(context, req, res)
	} else if router.notFound != nil {
		return router.notFound(context, req, res)
	}

	return nil
}
