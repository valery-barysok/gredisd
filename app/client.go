package app

import (
	"bufio"
	"net"
	"os"
	"sync"
	"time"

	"github.com/valery-barysok/gredisd/app/model"
	"github.com/valery-barysok/gredisd/server"
	"github.com/valery-barysok/resp"
)

type clientProvider struct {
	app    *App
	looper *looper
}

func NewClientProvider(app *App) server.ClientProvider {
	var protocol *resp.Protocol
	if app.opts.TraceProtocol {
		protocol = resp.NewProtocolWithLogging(os.Stdout)
	} else {
		protocol = resp.NewProtocol()
	}

	return &clientProvider{
		app:    app,
		looper: newLooper(protocol, app.router),
	}
}

func (cp *clientProvider) CreateClient(server *server.Server, conn net.Conn) (server.Client, error) {
	return newClient(server, conn, cp), nil
}

type ClientContext struct {
	App         *App
	DB          *model.DBModel
	RequireAuth bool
}

type client struct {
	mu        sync.Mutex
	id        uint64
	server    *server.Server
	conn      net.Conn
	bufReader *bufio.Reader
	bufWriter *bufio.Writer
	startTime time.Time
	last      time.Time
	looper    *looper
	context   *ClientContext
}

func newClientContext(app *App) *ClientContext {
	db, _ := app.SelectIndex(0)
	context := ClientContext{
		App:         app,
		DB:          db,
		RequireAuth: app.RequireAuth(),
	}
	return &context
}

func newClient(server *server.Server, conn net.Conn, cp *clientProvider) *client {
	client := &client{
		id:        server.GenerateClientID(),
		server:    server,
		conn:      conn,
		bufReader: bufio.NewReader(conn),
		bufWriter: bufio.NewWriter(conn),
		startTime: time.Now(),
		looper:    cp.looper,
		context:   newClientContext(cp.app),
	}

	return client
}

func (client *client) ID() uint64 {
	return client.id
}

func (client *client) Loop() {
	client.mu.Lock()
	conn := client.conn
	br := client.bufReader
	bw := client.bufWriter
	client.mu.Unlock()

	if conn == nil {
		return
	}

	client.looper.loop(client.context, br, bw)
}

func (client *client) CloseConnection() {
	client.mu.Lock()
	if client.conn == nil {
		client.mu.Unlock()
		return
	}

	client.clearConnection()

	client.conn = nil
	client.mu.Unlock()
}

func (client *client) clearConnection() {
	if client.conn == nil {
		return
	}

	client.conn.SetWriteDeadline(time.Now().Add(DefaultFlushDeadline))
	if client.bufWriter != nil {
		client.bufWriter.Flush()
	}
	client.conn.Close()
	client.conn.SetWriteDeadline(time.Time{})
}
