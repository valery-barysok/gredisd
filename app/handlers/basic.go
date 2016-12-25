package handlers

import (
	"errors"
	"log"

	"github.com/valery-barysok/gredisd/app"
	"github.com/valery-barysok/resp"
)

// List of basic commands.
const (
	AuthCommand     = "auth"
	SelectCommand   = "select"
	EchoCommand     = "echo"
	PingCommand     = "ping"
	ShutdownCommand = "shutdown"
	// TODO: use custom name "COMMANDS" instead of "COMMAND" due to incompatibility with redis-cli
	CommandCommand = "commands"
	KeysCommand    = "keys"
	ExistsCommand  = "exists"
	ExpireCommand  = "expire"
)

// BindAllBasicHandlers binds all basic commands at once
func BindAllBasicHandlers(app *app.App) {
	BindAuth(app)
	BindSelect(app)
	BindEcho(app)
	BindPing(app)
	BindShutdown(app)
	BindCommand(app)
	BindKeys(app)
	BindExists(app)
	BindExpire(app)

	BindNotFound(app)
	BindError(app)
}

// BindAuth binds auth command and required auth filter
func BindAuth(app *app.App) {
	app.BindFilter(authFilter)
	app.Bind(AuthCommand, authCmd)
}

// BindSelect binds Select command that select current database for specified client
func BindSelect(app *app.App) {
	app.Bind(SelectCommand, selectCmd)
}

// BindEcho binds Echo command that response with message back
func BindEcho(app *app.App) {
	app.Bind(EchoCommand, echoCmd)
}

// BindPing binds Ping command
func BindPing(app *app.App) {
	app.Bind(PingCommand, pingCmd)
}

// BindShutdown binds Select command that shutdown App
func BindShutdown(app *app.App) {
	app.Bind(ShutdownCommand, shutdownCmd)
}

// BindCommand binds Command command that list all available commands
func BindCommand(app *app.App) {
	app.Bind(CommandCommand, commandCmd)
}

func BindKeys(app *app.App) {
	app.Bind(KeysCommand, keysCmd)
}

func BindExists(app *app.App) {
	app.Bind(ExistsCommand, existsCmd)
}

func BindExpire(app *app.App) {
	app.Bind(ExpireCommand, expireCmd)
}

// BindNotFound binds handler for handling all unknown commands
func BindNotFound(appl *app.App) {
	appl.BindNotFound(func(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
		res.WriteUnknownCommandError(req.Cmd)
		res.Flush()
		return nil
	})
}

// BindError binds handler for handling all unknown commands
func BindError(appl *app.App) {
	appl.BindError(func(context *app.ClientContext, err error, res *resp.Writer) {
		log.Println(err)
		res.Flush()
	})
}

func authFilter(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) (bool, error) {
	if req.Cmd != AuthCommand && context.RequireAuth {
		res.WriteErrorString("NOAUTH Authentication required.")
		res.Flush()
		return true, nil
	}
	return false, nil
}

func authCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l != 1 {
		res.WriteArityError(req.Cmd)
	} else {
		context.RequireAuth = !context.App.Auth(string(req.Args[0].BulkString()))
		if context.RequireAuth {
			res.WriteErrorString("ERR invalid password")
		} else {
			res.WriteOK()
		}
	}
	res.Flush()
	return nil
}

func selectCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l != 1 {
		res.WriteArityError(req.Cmd)
	} else {
		db, err := context.App.Select(string(req.Args[0].BulkString()))
		if err != nil {
			res.WriteErrorString(err.Error())
		} else {
			context.DB = db
			res.WriteOK()
		}
	}
	res.Flush()
	return nil
}

func echoCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l != 1 {
		res.WriteArityError(req.Cmd)
	} else {
		res.WriteBulkString(req.Args[0].BulkString())
	}
	res.Flush()
	return nil
}

func pingCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l > 1 {
		res.WriteArityError(req.Cmd)
	} else if l == 1 {
		res.WriteBulkString(req.Args[0].BulkString())
	} else {
		res.WritePong()
	}
	res.Flush()
	return nil
}

func shutdownCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	context.App.Shutdown()
	return errors.New("Shutdown command received from client")
}

func commandCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	commands := context.App.Commands()
	res.WriteArray(commands)
	res.Flush()
	return nil
}

func keysCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l != 1 {
		res.WriteArityError(req.Cmd)
	} else {
		keys, err := context.DB.Keys(req.Args[0].BulkString())
		if err != nil {
			res.WriteError(err)
		} else {
			res.WriteArray(keys)
		}
	}
	res.Flush()
	return nil
}

func existsCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l < 1 {
		res.WriteArityError(req.Cmd)
	} else {
		keys := make([][]byte, 0, len(req.Args))
		for _, arg := range req.Args {
			keys = append(keys, arg.BulkString())
		}

		res.WriteInteger(context.DB.Exists(keys...))
	}
	res.Flush()
	return nil
}

func expireCmd(context *app.ClientContext, req *app.RespCommand, res *resp.Writer) error {
	l := len(req.Args)
	if l != 2 {
		res.WriteArityError(req.Cmd)
	} else {
		res.WriteInteger(context.DB.Expire(req.Args[0].BulkString(), req.Args[1].BulkString()))
	}
	res.Flush()
	return nil
}
