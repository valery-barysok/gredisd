package handlers

import (
	"github.com/valery-barysok/gredisd/app"
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

func hSetCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l != 3 {
		res.WriteArityError(req.Cmd)
	} else {
		cnt, err := context.DB.HSet(req.Args[0].BulkString(), req.Args[1].BulkString(), req.Args[2].BulkString())
		if err != nil {
			res.WriteError(err)
		} else {
			res.WriteInteger(cnt)
		}
	}
	res.Flush()
	return nil
}

func hGetCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l != 2 {
		res.WriteArityError(req.Cmd)
	} else {
		val, err := context.DB.HGet(req.Args[0].BulkString(), req.Args[1].BulkString())
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

func hDelCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l < 2 {
		res.WriteArityError(req.Cmd)
	} else {
		l := len(req.Args)
		keys := make([][]byte, 0, l-1)
		for i := 1; i < l; i++ {
			keys = append(keys, req.Args[i].BulkString())
		}

		cnt, err := context.DB.HDel(req.Args[0].BulkString(), keys...)
		if err != nil {
			res.WriteError(err)
		} else {
			res.WriteInteger(cnt)
		}
	}
	res.Flush()
	return nil
}

func hLenCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l != 1 {
		res.WriteArityError(req.Cmd)
	} else {
		l, err := context.DB.HLen(req.Args[0].BulkString())
		if err != nil {
			res.WriteError(err)
		} else {
			res.WriteInteger(l)
		}
	}
	res.Flush()
	return nil
}

func hExistsCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l != 2 {
		res.WriteArityError(req.Cmd)
	} else {
		cnt, err := context.DB.HExists(req.Args[0].BulkString(), req.Args[1].BulkString())
		if err != nil {
			res.WriteError(err)
		} else {
			res.WriteInteger(cnt)
		}
	}
	res.Flush()
	return nil
}
