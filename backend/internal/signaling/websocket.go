package signaling

import (
	"github.com/gorilla/websocket"
)

var Websocket = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
