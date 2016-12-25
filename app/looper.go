package app

import (
	"io"

	"github.com/valery-barysok/gredisd/app/cmd"
	"github.com/valery-barysok/resp"
)

type looper struct {
	protocol *resp.Protocol
	router   *router
}

func newLooper(protocol *resp.Protocol, router *router) *looper {
	return &looper{
		protocol: protocol,
		router:   router,
	}
}

func (looper *looper) loop(context *ClientContext, r io.Reader, w io.Writer) {
	reader := resp.NewReader(r, looper.protocol)
	writer := resp.NewWriter(w, looper.protocol)
	for {
		cmd, err := cmd.ReadCommand(reader)
		if err != nil {
			if looper.router.errorHandler != nil {
				looper.router.errorHandler(context, err, writer)
			}
			return
		}

		if err := looper.router.serve(context, cmd, writer); err != nil {
			if looper.router.errorHandler != nil {
				looper.router.errorHandler(context, err, writer)
			}
			return
		}
	}
}
