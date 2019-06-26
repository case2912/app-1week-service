package main

type Message struct {
	Message     string `json:"message"`
	From        string `json:"from"`
	MessageType string `json:"messageType"`
}
