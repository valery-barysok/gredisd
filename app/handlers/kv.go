package handlers

import (
	"github.com/valery-barysok/gredisd/app"
	"github.com/valery-barysok/gredisd/app/cmd"
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

func setCmd(context *app.ClientContext, cmd *cmd.Command, res *resp.Writer) error {
	l := len(cmd.Args)
	if l < 2 {
		res.WriteArityError(cmd.Cmd)
	} else {
		context.DB.Set(cmd.Args[0].BulkString(), cmd.Args[1].BulkString())
		res.WriteOK()
	}
	res.Flush()
	return nil
}

func getCmd(context *app.ClientContext, cmd *cmd.Command, res *resp.Writer) error {
	l := len(cmd.Args)
	if l != 1 {
		res.WriteArityError(cmd.Cmd)
	} else {
		val, err := context.DB.Get(cmd.Args[0].BulkString())
		if err != nil {
			res.WriteError(err)
		} else if val != nil {
			res.WriteBulkString(val)
		} else {
			res.WriteNilBulk()
		}
	}
	res.Flush()
	return nil
}

func delCmd(context *app.ClientContext, cmd *cmd.Command, res *resp.Writer) error {
	l := len(cmd.Args)
	if l < 1 {
		res.WriteArityError(cmd.Cmd)
	} else {
		keys := make([][]byte, 0, len(cmd.Args))
		for _, arg := range cmd.Args {
			keys = append(keys, arg.BulkString())
		}

		cnt := context.DB.Del(keys...)
		res.WriteInteger(cnt)
	}
	res.Flush()
	return nil
}
