package signaling

import (
	"backend/internal/webrtc"
	"encoding/json"
	"github.com/gorilla/websocket"
	pion "github.com/pion/webrtc/v4"
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

	// Unmarshal the signaling message
	var signal Signal
	if err := json.Unmarshal(message, &signal); err != nil {
		log.Println("Error unmarshalling signaling message:", err)
		return
	}

	switch signal.Type {
	case "offer":
		// Handle offer message
		handleOffer(session, signal.Data, conn)
	case "answer":
		// Handle answer message
		handleAnswer(session, signal.Data)
	case "iceCandidate":
		// Handle ICE candidate message
		handleICECandidate(session, signal.Data)
	default:
		log.Printf("Unsupported message type: %s", signal.Type)
	}
}

func handleOffer(session *webrtc.Session, offerData json.RawMessage, conn *websocket.Conn) {
	// Unmarshal the offer data
	var offer pion.SessionDescription
	if err := json.Unmarshal(offerData, &offer); err != nil {
		log.Println("Error unmarshalling offer:", err)
		return
	}

	// Create a new peer connection
	pc, err := pion.NewPeerConnection(pion.Configuration{})
	if err != nil {
		log.Println("Error creating peer connection:", err)
		return
	}
	session.PeerConnection = pc // Set peer connection in the session

	// Set remote description
	if err := pc.SetRemoteDescription(offer); err != nil {
		log.Println("Error setting remote description:", err)
		return
	}

	// Create answer
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		log.Println("Error creating answer:", err)
		return
	}

	// Set local description
	if err := pc.SetLocalDescription(answer); err != nil {
		log.Println("Error setting local description:", err)
		return
	}

	// Send answer to the remote peer
	if err := sendSignalingMessage(conn, "answer", answer); err != nil {
		log.Println("Error sending answer:", err)
		return
	}
}

func handleAnswer(session *webrtc.Session, answerData json.RawMessage) {
	// Unmarshal the answer data
	var answer pion.SessionDescription
	if err := json.Unmarshal(answerData, &answer); err != nil {
		log.Println("Error unmarshalling answer:", err)
		return
	}

	// Set remote description
	if err := session.PeerConnection.SetRemoteDescription(answer); err != nil {
		log.Println("Error setting remote description:", err)
		return
	}
}

func handleICECandidate(session *webrtc.Session, candidateData json.RawMessage) {
	// Unmarshal the ICE candidate data
	var candidate pion.ICECandidateInit
	if err := json.Unmarshal(candidateData, &candidate); err != nil {
		log.Println("Error unmarshalling ICE candidate:", err)
		return
	}

	// Add ICE candidate to the peer connection
	if err := session.PeerConnection.AddICECandidate(candidate); err != nil {
		log.Println("Error adding ICE candidate:", err)
		return
	}
}

func sendSignalingMessage(conn *websocket.Conn, signalType string, data interface{}) error {
	// Marshal the data into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Send the JSON message over WebSocket
	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		return err
	}

	return nil
}
