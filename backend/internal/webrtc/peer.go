package webrtc

import "github.com/pion/webrtc/v4"

type Peer struct {
	ID   string
	Conn *webrtc.PeerConnection
}
