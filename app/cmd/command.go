package cmd

import (
	"bytes"
	"errors"

	"github.com/valery-barysok/resp"
)

var errInvalidRequest = errors.New("ERR Invalid Request")

type Command struct {
	Cmd  string
	Args []*resp.Message
}

func ReadCommand(reader *resp.Reader) (*Command, error) {
	msg, err := reader.Read()
	if err != nil {
		return nil, err
	}

	switch msgs := msg.Value.(type) {
	case []*resp.Message:
		return &Command{
			Cmd:  string(bytes.ToLower(msgs[0].BulkString())),
			Args: msgs[1:],
		}, nil
	}

	return nil, errInvalidRequest
}
