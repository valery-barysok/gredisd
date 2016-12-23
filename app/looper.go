package app

import (
	"bytes"
	"errors"
	"io"

	"github.com/valery-barysok/resp"
)

var errInvalidRequest = errors.New("ERR Invalid Request")

type RespRequest struct {
	Cmd  string
	Args []*resp.Item
}

type looper struct {
	protocol *resp.Protocol
	router   *Router
}

func newLooper(protocol *resp.Protocol, router *Router) *looper {
	return &looper{
		protocol: protocol,
		router:   router,
	}
}

func (looper *looper) Loop(context *ClientContext, r io.Reader, w io.Writer) {
	responder := resp.NewWriter(w, looper.protocol)
	reader := resp.NewReader(r, looper.protocol)
	for {
		req, err := looper.ReadRequest(reader)
		if err != nil {
			if looper.router.errorHandler != nil {
				looper.router.errorHandler(context, err, responder)
			}
			return
		}

		if err := looper.router.Serve(context, req, responder); err != nil {
			if looper.router.errorHandler != nil {
				looper.router.errorHandler(context, err, responder)
			}
			return
		}
	}
}

func (looper *looper) ReadRequest(reader *resp.Reader) (*RespRequest, error) {
	item, err := reader.Read()
	if err != nil {
		return nil, err
	}

	switch it := item.Value.(type) {
	case []*resp.Item:
		items := it
		return &RespRequest{
			Cmd:  string(bytes.ToLower(items[0].BulkString())),
			Args: items[1:],
		}, nil
	default:
		return nil, errInvalidRequest
	}
}
