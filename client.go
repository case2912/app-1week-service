package main

import "github.com/gollira/websocket"

type client struct {
	socket *websocket.Conn
	send   chan *Message
	room   *room
}

func (c *client) read() {
	for {
		msg := &Message{}
		err := c.socket.ReadJSON(msg)
		if err != nil {
			break
		}
		c.room.forward <- msg
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
	c.socket.Close()
}
