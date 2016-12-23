package app

import "time"

const (
	// Version is current version of app
	Version = "0.0.1"

	// DefaultPort is port for incoming connections by default
	DefaultPort = 16379

	// DefaultHost is address to listen for incoming connections by default
	DefaultHost = "0.0.0.0"

	// DefaultDatabases is number of available databases by default
	DefaultDatabases = 16

	// DefaultFlushDeadline is timeout for writing to the client by default
	DefaultFlushDeadline = 2 * time.Second
)
