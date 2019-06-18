package main

type room struct {
	id      string
	forward chan *Message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func newRoom(id string) *room {
	return &room{
		id:      id,
		forward: make(chan *Message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
