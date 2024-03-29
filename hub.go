package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gollira/websocket"
	uuid "github.com/satori/go.uuid"
)

type hub struct {
	rooms      map[string]*room
	register   chan *room
	unregister chan *room
}

func newHub() *hub {
	return &hub{
		rooms: make(map[string]*room),
	}
}
func (h *hub) run(id string) {
	r := h.rooms[id]
	for {
		select {
		case client := <-r.join:
			fmt.Printf("client joined ROOM:%s\n", id)
			r.clients[client] = true
		case client := <-r.leave:
			fmt.Println("client left")
			delete(r.clients, client)
			close(client.send)
			if len(r.clients) == 0 {
				fmt.Printf("room deleted ROOM:%s\n", id)
				delete(h.rooms, r.id)
			}
		case msg := <-r.forward:
			fmt.Printf("client forward message %+v\n", msg)
			if msg.MessageType == "Haiku" {
				r.renga = append(r.renga, msg)
			}

			fmt.Printf("renga %+v\n", r.renga)
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (h *hub) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	_, ok := h.rooms[id]
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	upgrader.CheckOrigin = func(req *http.Request) bool {
		return true
	}
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan *Message, messageBufferSize),
		room:   h.rooms[id],
	}
	h.rooms[id].join <- client
	h.rooms[id].forward <- &Message{
		Message:     "client joined",
		MessageType: "Comment",
	}
	for _, msg := range h.rooms[id].renga {
		client.socket.WriteJSON(msg)
	}
	defer func() {
		h.rooms[id].leave <- client
		h.rooms[id].forward <- &Message{
			Message:     "client left",
			MessageType: "Comment",
		}
	}()
	go client.write()
	client.read()
}

type Room struct {
	RoomID string `json:"room_id"`
}

func (h *hub) CreateRoom(w http.ResponseWriter, req *http.Request) {
	id := uuid.NewV4().String()
	new := newRoom(id)
	h.rooms[new.id] = new
	go h.run(id)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Println("request create room")
	json.NewEncoder(w).Encode(Room{
		RoomID: id,
	})
}

type RoomList struct {
	Rooms []string `json:"rooms"`
}

func (h *hub) GetRoomList(w http.ResponseWriter, req *http.Request) {
	keys := make([]string, 0, len(h.rooms))
	for roomID := range h.rooms {
		keys = append(keys, roomID)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Println("request get room list")
	json.NewEncoder(w).Encode(RoomList{
		Rooms: keys,
	})
}
