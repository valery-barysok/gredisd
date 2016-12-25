package handlers

import (
	"github.com/valery-barysok/gredisd/app"
	"github.com/valery-barysok/gredisd/app/cmd"
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

func lrpushCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer, push lrPush) error {
	l := len(cmd.Args)
	if l < 1 {
		w.WriteArityError(cmd.Cmd)
	} else {
		l := len(cmd.Args)
		values := make([][]byte, 0, l-1)
		for i := 1; i < l; i++ {
			values = append(values, cmd.Args[i].BulkString())
		}
		cnt, err := push(context.DB, cmd.Args[0].BulkString(), values...)
		if err != nil {
			w.WriteError(err)
		} else {
			w.WriteInteger(cnt)
		}
	}
	w.Flush()
	return nil
}

func lpushCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	return lrpushCmd(context, cmd, w, func(db *model.DBModel, key []byte, values ...[]byte) (int, error) {
		return db.LPush(key, values...)
	})
}

func rpushCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	return lrpushCmd(context, cmd, w, func(db *model.DBModel, key []byte, values ...[]byte) (int, error) {
		return db.RPush(key, values...)
	})
}

func lrpopCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer, pop lrPop) error {
	l := len(cmd.Args)
	if l != 1 {
		w.WriteArityError(cmd.Cmd)
	} else {
		value, err := pop(context.DB, cmd.Args[0].BulkString())
		if err != nil {
			w.WriteError(err)
		} else if value != nil {
			w.WriteBulkString(value)
		} else {
			w.WriteNilBulk()
		}
	}
	w.Flush()
	return nil
}

func lpopCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	return lrpopCmd(context, cmd, w, func(db *model.DBModel, key []byte) ([]byte, error) {
		return db.LPop(key)
	})
}

func rpopCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	return lrpopCmd(context, cmd, w, func(db *model.DBModel, key []byte) ([]byte, error) {
		return db.RPop(key)
	})
}

func llenCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	l := len(cmd.Args)
	if l != 1 {
		w.WriteArityError(cmd.Cmd)
	} else {
		cnt, err := context.DB.LLen(cmd.Args[0].BulkString())
		if err != nil {
			w.WriteError(err)
		} else {
			w.WriteInteger(cnt)
		}
	}
	w.Flush()
	return nil
}

func linsertCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	l := len(cmd.Args)
	if l != 4 {
		w.WriteArityError(cmd.Cmd)
	} else {
		l, err := context.DB.LInsert(cmd.Args[0].BulkString(), cmd.Args[1].BulkString(),
			cmd.Args[2].BulkString(), cmd.Args[3].BulkString())
		if err != nil {
			w.WriteError(err)
		} else {
			w.WriteInteger(l)
		}
	}
	w.Flush()
	return nil
}

func lindexCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	l := len(cmd.Args)
	if l != 2 {
		w.WriteArityError(cmd.Cmd)
	} else {
		s, err := context.DB.LIndex(cmd.Args[0].BulkString(), cmd.Args[1].BulkString())
		if err != nil {
			w.WriteError(err)
		} else if s != nil {
			w.WriteBulkString(s)
		} else {
			w.WriteNilBulk()
		}
	}
	w.Flush()
	return nil
}

func lrangeCmd(context *app.ClientContext, cmd *cmd.Command, w *resp.Writer) error {
	l := len(cmd.Args)
	if l != 3 {
		w.WriteArityError(cmd.Cmd)
	} else {
		values, err := context.DB.LRange(cmd.Args[0].BulkString(), cmd.Args[1].BulkString(), cmd.Args[2].BulkString())
		if err != nil {
			w.WriteError(err)
		} else if values != nil {
			w.WriteArray(values)
		} else {
			w.WriteNilBulk()
		}
	}
	w.Flush()
	return nil
}
