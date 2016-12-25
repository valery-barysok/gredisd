package app

import (
	"bytes"
	"errors"
	"io"

	"github.com/valery-barysok/resp"
)

var errInvalidRequest = errors.New("ERR Invalid Request")

type RespCommand struct {
	Cmd  string
	Args []*resp.Item
}

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
	responder := resp.NewWriter(w, looper.protocol)
	reader := resp.NewReader(r, looper.protocol)
	for {
		req, err := looper.readRespCommand(reader)
		if err != nil {
			if looper.router.errorHandler != nil {
				looper.router.errorHandler(context, err, responder)
			}
			return
		}

		if err := looper.router.serve(context, req, responder); err != nil {
			if looper.router.errorHandler != nil {
				looper.router.errorHandler(context, err, responder)
			}
			return
		}
	}
}

func (looper *looper) readRespCommand(reader *resp.Reader) (*RespCommand, error) {
	item, err := reader.Read()
	if err != nil {
		return nil, err
	}

	switch it := item.Value.(type) {
	case []*resp.Item:
		items := it
		return &RespCommand{
			Cmd:  string(bytes.ToLower(items[0].BulkString())),
			Args: items[1:],
		}, nil
	default:
		return nil, errInvalidRequest
	}
}
