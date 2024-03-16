package webrtc

import "github.com/pion/webrtc/v4"

type Session struct {
	Peers          []*Peer
	PeerConnection *webrtc.PeerConnection // Add PeerConnection field
	rtcManager     *WebRTCManager
}

type Peer struct {
	ID   string
	Conn *webrtc.PeerConnection // Use PeerConnection from pion/webrtc/v4
}
