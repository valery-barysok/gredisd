/*
Package server provides implementation of simple tcp server

How to create server:

	opts := server.Options{
		Host: "localhost",
		Port: 8080,
	}

	srv := server.NewServer(&opts, NewClientProvider())
	srv.Start()

	time.Sleep(5 * time.Second)

	srv.Shutdown()

where ClientProvider is intended for supply different types of Clients
based on provided Server and accepted connection

Client lifecycle introduced by this package is:
	Create Client on accepted connection
	Loop for incoming messages
	Close connection
*/
package server // import "github.com/valery-barysok/gredisd/server"
