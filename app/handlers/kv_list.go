package handlers

import (
	"github.com/valery-barysok/gredisd/app"
	"github.com/valery-barysok/gredisd/app/model"
	"github.com/valery-barysok/resp"
)

// List of key value list commands.
const (
	LPushCommand   = "lpush"
	RPushCommand   = "rpush"
	LPopCommand    = "lpop"
	RPopCommand    = "rpop"
	LLenCommand    = "llen"
	LInsertCommand = "linsert"
	LIndexCommand  = "lindex"
	LRangeCommand  = "lrange"
)

// BindAllKVListHandlers binds all key value list commands at once
func BindAllKVListHandlers(app *app.App) {
	BindLPush(app)
	BindRPush(app)
	BindLPop(app)
	BindRPop(app)
	BindLLen(app)
	BindLInsert(app)
	BindLIndex(app)
	BindLRange(app)
}

func BindLPush(app *app.App) {
	app.Bind(LPushCommand, lpushCmd)
}

func BindRPush(app *app.App) {
	app.Bind(RPushCommand, rpushCmd)
}

func BindLPop(app *app.App) {
	app.Bind(LPopCommand, lpopCmd)
}

func BindRPop(app *app.App) {
	app.Bind(RPopCommand, rpopCmd)
}

func BindLLen(app *app.App) {
	app.Bind(LLenCommand, llenCmd)
}

func BindLInsert(app *app.App) {
	app.Bind(LInsertCommand, linsertCmd)
}

func BindLIndex(app *app.App) {
	app.Bind(LIndexCommand, lindexCmd)
}

func BindLRange(app *app.App) {
	app.Bind(LRangeCommand, lrangeCmd)
}

type lrPush func(db *model.DBModel, key []byte, values ...[]byte) (int, error)
type lrPop func(db *model.DBModel, key []byte) ([]byte, error)

func lrpushCmd(context *app.ClientContext, req *app.RespRequest, res *resp.Writer, push lrPush) error {
	l := len(req.Args)
	if l < 1 {
		res.WriteArityError(req.Cmd)
	} else {
		l := len(req.Args)
		values := make([][]byte, 0, l-1)
		for i := 1; i < l; i++ {
			values = append(values, req.Args[i].BulkString())
		}
		cnt, err := push(context.DB, req.Args[0].BulkString(), values...)
		if err != nil {
			res.WriteError(err)
		} else {
			res.WriteInteger(cnt)
		}
	}
	res.End()
	return nil
}

func lpushCmd(context *app.ClientContext, req *app.RespRequest, res *resp.Writer) error {
	return lrpushCmd(context, req, res, func(db *model.DBModel, key []byte, values ...[]byte) (int, error) {
		return db.LPush(key, values...)
	})
}

func rpushCmd(context *app.ClientContext, req *app.RespRequest, res *resp.Writer) error {
	return lrpushCmd(context, req, res, func(db *model.DBModel, key []byte, values ...[]byte) (int, error) {
		return db.RPush(key, values...)
	})
}

func lrpopCmd(context *app.ClientContext, req *app.RespRequest, res *resp.Writer, pop lrPop) error {
	l := len(req.Args)
	if l != 1 {
		res.WriteArityError(req.Cmd)
	} else {
		value, err := pop(context.DB, req.Args[0].BulkString())
		if err != nil {
			res.WriteError(err)
		} else if value != nil {
			res.WriteBulkString(value)
		} else {
			res.WriteNilBulk()
		}
	}
	res.End()
	return nil
}

func lpopCmd(context *app.ClientContext, req *app.RespRequest, res *resp.Writer) error {
	return lrpopCmd(context, req, res, func(db *model.DBModel, key []byte) ([]byte, error) {
		return db.LPop(key)
	})
}

func rpopCmd(context *app.ClientContext, req *app.RespRequest, res *resp.Writer) error {
	return lrpopCmd(context, req, res, func(db *model.DBModel, key []byte) ([]byte, error) {
		return db.RPop(key)
	})
}

func llenCmd(context *app.ClientContext, req *app.RespRequest, res *resp.Writer) error {
	l := len(req.Args)
	if l != 1 {
		res.WriteArityError(req.Cmd)
	} else {
		cnt, err := context.DB.LLen(req.Args[0].BulkString())
		if err != nil {
			res.WriteError(err)
		} else {
			res.WriteInteger(cnt)
		}
	}
	res.End()
	return nil
}

func linsertCmd(context *app.ClientContext, req *app.RespRequest, res *resp.Writer) error {
	l := len(req.Args)
	if l != 4 {
		res.WriteArityError(req.Cmd)
	} else {
		l, err := context.DB.LInsert(req.Args[0].BulkString(), req.Args[1].BulkString(),
			req.Args[2].BulkString(), req.Args[3].BulkString())
		if err != nil {
			res.WriteError(err)
		} else {
			res.WriteInteger(l)
		}
	}
	res.End()
	return nil
}

func lindexCmd(context *app.ClientContext, req *app.RespRequest, res *resp.Writer) error {
	l := len(req.Args)
	if l != 2 {
		res.WriteArityError(req.Cmd)
	} else {
		s, err := context.DB.LIndex(req.Args[0].BulkString(), req.Args[1].BulkString())
		if err != nil {
			res.WriteError(err)
		} else if s != nil {
			res.WriteBulkString(s)
		} else {
			res.WriteNilBulk()
		}
	}
	res.End()
	return nil
}

func lrangeCmd(context *app.ClientContext, req *app.RespRequest, res *resp.Writer) error {
	l := len(req.Args)
	if l != 3 {
		res.WriteArityError(req.Cmd)
	} else {
		values, err := context.DB.LRange(req.Args[0].BulkString(), req.Args[1].BulkString(), req.Args[2].BulkString())
		if err != nil {
			res.WriteError(err)
		} else if values != nil {
			res.WriteArray(values)
		} else {
			res.WriteNilBulk()
		}
	}
	res.End()
	return nil
}
