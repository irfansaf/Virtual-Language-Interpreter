package signaling

import (
	"backend/internal/webrtc"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Websocket.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Create a new WebRTC session for each WebSocket connection
	sessionID := "unique_session_id" // You may generate a unique session ID here
	webrtc.NewWebRTCManager().CreateSession(sessionID)

	// Handle signaling messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message: %s", message)

		// Handle WebRTC signaling message
		handleSignalingMessage(sessionID, conn, messageType, message)
	}
}

func handleSignalingMessage(sessionID string, conn *websocket.Conn, messageType int, message []byte) {
	// retrieve WebRTC session
	session := webrtc.NewWebRTCManager().GetSession(sessionID)
	if session == nil {
		log.Printf("WebRTC session %s not found", sessionID)
		return
	}

	// Process signaling message based on message type
	// For now, we'll just echo the message back to the sender
	// You'll implement the actual WebRTC signaling logic here
	if err := conn.WriteMessage(messageType, message); err != nil {
		log.Println("Error writing message:", err)
	}
}
