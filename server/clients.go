package server

type clientRegistry struct {
	clients map[uint64]Client
}

func newClientRegistry() *clientRegistry {
	return &clientRegistry{
		clients: make(map[uint64]Client),
	}
}

func (clients *clientRegistry) add(client Client) {
	clients.clients[client.ID()] = client
}

func (clients *clientRegistry) del(client Client) {
	delete(clients.clients, client.ID())
}

func (clients *clientRegistry) detach() map[uint64]Client {
	tmp := clients.clients
	clients.clients = nil
	return tmp
}
