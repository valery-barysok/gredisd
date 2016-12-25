package handlers

import (
	"github.com/valery-barysok/gredisd/app"
	"github.com/valery-barysok/gredisd/app/cmd"
	"github.com/valery-barysok/resp"
)

// List of key value dict commands.
const (
	HSetCommand    = "hset"
	HGetCommand    = "hget"
	HDelCommand    = "hdel"
	HLenCommand    = "hlen"
	HExistsCommand = "hexists"
)

// BindAllKVDictHandlers binds all key value dict commands at once
func BindAllKVDictHandlers(app *app.App) {
	BindHSet(app)
	BindHGet(app)
	BindHDel(app)
	BindHLen(app)
	BindHExists(app)
}

func BindHSet(app *app.App) {
	app.Bind(HSetCommand, hSetCmd)
}

func BindHGet(app *app.App) {
	app.Bind(HGetCommand, hGetCmd)
}

func BindHDel(app *app.App) {
	app.Bind(HDelCommand, hDelCmd)
}

func BindHLen(app *app.App) {
	app.Bind(HLenCommand, hLenCmd)
}

func BindHExists(app *app.App) {
	app.Bind(HExistsCommand, hExistsCmd)
}

func hSetCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	l := len(cmd.Args)
	if l != 3 {
		w.WriteArityError(cmd.Cmd)
	} else {
		cnt, err := context.DB.HSet(cmd.Args[0].BulkString(), cmd.Args[1].BulkString(), cmd.Args[2].BulkString())
		if err != nil {
			w.WriteError(err)
		} else {
			w.WriteInteger(cnt)
		}
	}
	w.Flush()
	return nil
}

func hGetCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	l := len(cmd.Args)
	if l != 2 {
		w.WriteArityError(cmd.Cmd)
	} else {
		val, err := context.DB.HGet(cmd.Args[0].BulkString(), cmd.Args[1].BulkString())
		if err != nil {
			w.WriteError(err)
		} else if val != nil {
			w.WriteBulkString(val)
		} else {
			w.WriteNilBulk()
		}
	}
	w.Flush()
	return nil
}

func hDelCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	l := len(cmd.Args)
	if l < 2 {
		w.WriteArityError(cmd.Cmd)
	} else {
		l := len(cmd.Args)
		keys := make([][]byte, 0, l-1)
		for i := 1; i < l; i++ {
			keys = append(keys, cmd.Args[i].BulkString())
		}

		cnt, err := context.DB.HDel(cmd.Args[0].BulkString(), keys...)
		if err != nil {
			w.WriteError(err)
		} else {
			w.WriteInteger(cnt)
		}
	}
	w.Flush()
	return nil
}

func hLenCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	l := len(cmd.Args)
	if l != 1 {
		w.WriteArityError(cmd.Cmd)
	} else {
		l, err := context.DB.HLen(cmd.Args[0].BulkString())
		if err != nil {
			w.WriteError(err)
		} else {
			w.WriteInteger(l)
		}
	}
	w.Flush()
	return nil
}

func hExistsCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	l := len(cmd.Args)
	if l != 2 {
		w.WriteArityError(cmd.Cmd)
	} else {
		cnt, err := context.DB.HExists(cmd.Args[0].BulkString(), cmd.Args[1].BulkString())
		if err != nil {
			w.WriteError(err)
		} else {
			w.WriteInteger(cnt)
		}
	}
	w.Flush()
	return nil
}
