package signaling

import (
	"backend/internal/webrtc"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func WebSocketHandler(rtcManager *webrtc.WebRTCManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := Websocket.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}
		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
				log.Println("Error closing WebSocket:", err)
			}
		}(conn)

		// Create a new WebRTC session for each WebSocket connection
		sessionID := "unique_session_id" // You may generate a unique session ID here
		rtcManager.CreateSession(sessionID)

		// Handle signaling messages
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
			log.Printf("Received message: %s", message)

			// Handle WebRTC signaling message
			handleSignalingMessage(sessionID, rtcManager, conn, messageType, message)
		}
	}
}

func handleSignalingMessage(sessionID string, rtcManager *webrtc.WebRTCManager, conn *websocket.Conn, messageType int, message []byte) {
	// Retrieve WebRTC session
	session := rtcManager.GetSession(sessionID)
	if session == nil {
		log.Printf("WebRTC session %s not found", sessionID)
		return
	}

	// Process signaling message based on message type
	switch messageType {
	case websocket.TextMessage, websocket.BinaryMessage:
		// Handle offer, answer, or ICE candidate message
		// You'll need to implement this logic based on your signaling protocol
		if err := conn.HandleMessage(message); err != nil {
			log.Println("Error writing message:", err)
		}
	default:
		log.Printf("Unsupported message type: %d", messageType)
	}

	// For now, we'll just echo the message back to the sender
	// You'll implement the actual WebRTC signaling logic here
	if err := conn.WriteMessage(messageType, message); err != nil {
		log.Println("Error writing message:", err)
	}
}
