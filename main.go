package main

import (
	"log"
	"net/http"
)

func main() {
	h := newHub()
	http.HandleFunc("/roomlist", h.GetRoomList)
	http.Handle("/room", h)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
