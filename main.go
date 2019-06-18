package main

import (
	"log"
	"net/http"
)

func main() {
	h := newHub()

	http.Handle("/room", h)
	http.HandleFunc("/room/create", h.CreateRoom)
	http.HandleFunc("/roomlist", h.GetRoomList)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
