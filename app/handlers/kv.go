package handlers

import (
	"github.com/valery-barysok/gredisd/app"
	"github.com/valery-barysok/resp"
)

// List of key value commands.
const (
	SetCommand = "set"
	GetCommand = "get"
	DelCommand = "del"
)

// BindAllKVHandlers binds all key value commands at once
func BindAllKVHandlers(app *app.App) {
	BindSet(app)
	BindGet(app)
	BindDel(app)
}

func BindSet(app *app.App) {
	app.Bind(SetCommand, setCmd)
}

func BindGet(app *app.App) {
	app.Bind(GetCommand, getCmd)
}

func BindDel(app *app.App) {
	app.Bind(DelCommand, delCmd)
}

func setCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l < 2 {
		res.WriteArityError(req.Cmd)
	} else {
		context.DB.Set(req.Args[0].BulkString(), req.Args[1].BulkString())
		res.WriteOK()
	}
	res.End()
	return nil
}

func getCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l != 1 {
		res.WriteArityError(req.Cmd)
	} else {
		val, err := context.DB.Get(req.Args[0].BulkString())
		if err != nil {
			res.WriteError(err)
		} else if val != nil {
			res.WriteBulkString(val)
		} else {
			res.WriteNilBulk()
		}
	}
	res.End()
	return nil
}

func delCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l < 1 {
		res.WriteArityError(req.Cmd)
	} else {
		keys := make([][]byte, 0, len(req.Args))
		for _, arg := range req.Args {
			keys = append(keys, arg.BulkString())
		}

		cnt := context.DB.Del(keys...)
		res.WriteInteger(cnt)
	}
	res.End()
	return nil
}
