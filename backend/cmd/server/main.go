package main

import (
	"backend/internal/signaling"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", signaling.WebSocketHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
