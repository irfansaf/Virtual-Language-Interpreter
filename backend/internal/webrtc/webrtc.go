package webrtc

import (
	"github.com/pion/webrtc/v4"
	"log"
	"sync"
)

type WebRTCManager struct {
	sessions map[string]*Session
	mutex    sync.Mutex
}

func NewWebRTCManager() *WebRTCManager {
	return &WebRTCManager{
		sessions: make(map[string]*Session),
	}
}

// CreateSession creates a new WebRTC session
func (w *WebRTCManager) CreateSession(sessionID string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if _, ok := w.sessions[sessionID]; !ok {
		w.sessions[sessionID] = &Session{
			Peers: make([]*Peer, 0),
		}
		log.Printf("Created WebRTC session: %s", sessionID)
	}
}

// GetSession retrieves a WebRTC session by its ID
func (w *WebRTCManager) GetSession(sessionID string) *Session {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if session, ok := w.sessions[sessionID]; ok {
		return session
	}

	return nil
}

// AddPeer adds a new peer to the WebRTC session
func (w *WebRTCManager) AddPeer(sessionID, peerID string, peerConnection *webrtc.PeerConnection) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	session := w.sessions[sessionID]
	if session == nil {
		log.Printf("WebRTC session %s not found", sessionID)
		return
	}

	peer := &Peer{
		ID:   peerID,
		Conn: peerConnection,
	}
	session.Peers = append(session.Peers, peer)
	log.Printf("Added peer %s to WebRTC session %s", peerID, sessionID)
}
