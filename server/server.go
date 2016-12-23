package server

import (
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	acceptMinSleep = 10 * time.Millisecond
	acceptMaxSleep = 1 * time.Second
)

// Client serves specific incoming connection
type Client interface {
	ID() uint64
	Loop()
	CloseConnection()
}

// ClientProvider creates Client for specific incoming connection
type ClientProvider interface {
	CreateClient(server *Server, conn net.Conn) (Client, error)
}

// Server contains various details about specific server like connected clients, server options etc
type Server struct {
	// last client id. used for generate next client id
	cid            uint64
	mu             sync.Mutex
	opts           *Options
	running        bool
	tcp            net.Listener
	startTime      time.Time
	clients        *clientRegistry
	clientProvider ClientProvider

	grMu      sync.Mutex
	grRunning bool
	grWG      sync.WaitGroup
}

// NewServer creates server based on provided options and client provider
func NewServer(opts *Options, clientProvider ClientProvider) *Server {
	s := &Server{
		opts:           opts,
		startTime:      time.Now(),
		clientProvider: clientProvider,
		clients:        newClientRegistry(),
	}

	return s
}

// GenerateClientID generates next unique client id for specified server
func (server *Server) GenerateClientID() uint64 {
	return atomic.AddUint64(&server.cid, 1)
}

// Start begins to listen for incoming connections
func (server *Server) Start() {
	server.mu.Lock()
	server.running = true
	server.mu.Unlock()

	server.grMu.Lock()
	server.grRunning = true
	server.grMu.Unlock()

	server.acceptLoop()
}

// Shutdown stops to listen for incoming connections and closes all opened connections
func (server *Server) Shutdown() {
	server.mu.Lock()

	if !server.running {
		server.mu.Unlock()
		return
	}

	server.running = false
	server.grMu.Lock()
	server.grRunning = false
	server.grMu.Unlock()

	clients := server.clients.detach()

	server.mu.Unlock()

	for _, c := range clients {
		c.CloseConnection()
	}

	server.grWG.Wait()
	log.Print("Server shutdown")
}

func (server *Server) acceptLoop() {
	hp := net.JoinHostPort(server.opts.Host, strconv.Itoa(server.opts.Port))
	l, e := net.Listen("tcp", hp)
	if e != nil {
		return
	}

	log.Printf("Listening for client connections on %s", hp)

	server.mu.Lock()
	server.tcp = l
	server.mu.Unlock()

	tmpDelay := acceptMinSleep

	for server.isRunning() {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				time.Sleep(tmpDelay)
				tmpDelay *= 2
				if tmpDelay > acceptMaxSleep {
					tmpDelay = acceptMaxSleep
				}
			}
			continue
		}

		tmpDelay = acceptMinSleep
		server.startGoRoutine(func() {
			log.Printf("Server accepted connection on %s", conn.RemoteAddr())

			defer server.grWG.Done()
			server.createClient(conn)
		})
	}
}

func (server *Server) isRunning() bool {
	server.mu.Lock()
	defer server.mu.Unlock()
	return server.running
}

func (server *Server) createClient(conn net.Conn) (Client, error) {
	client, err := server.clientProvider.CreateClient(server, conn)
	if err != nil {
		return nil, err
	}

	server.mu.Lock()
	if !server.running {
		server.mu.Unlock()
		return client, nil
	}
	server.clients.add(client)
	server.mu.Unlock()

	server.startGoRoutine(func() {
		defer func() {
			server.removeClient(client)
			server.grWG.Done()
		}()
		client.Loop()
	})

	return client, nil
}

func (server *Server) removeClient(client Client) {
	server.mu.Lock()
	server.clients.del(client)
	server.mu.Unlock()

	client.CloseConnection()
}

func (server *Server) startGoRoutine(f func()) {
	server.grMu.Lock()
	if server.grRunning {
		server.grWG.Add(1)
		go f()
	}
	server.grMu.Unlock()
}
