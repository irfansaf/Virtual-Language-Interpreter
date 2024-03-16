package signaling

import "encoding/json"

type Signal struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
