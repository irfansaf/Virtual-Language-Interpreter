package main

import (
	"backend/internal/signaling"
	"backend/internal/webrtc"
	"net/http"
)

func main() {
	rtcManager := webrtc.NewWebRTCManager()

	http.HandleFunc("/ws", signaling.WebSocketHandler(rtcManager))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
