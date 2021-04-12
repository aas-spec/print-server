package ps_event_bus

import "encoding/json"

type EventBusMessage struct {
	CommandType  string          `json:"type"`
	CommandParam string          `json:"param,omitempty"`
	StickerType  string          `json:"sticker_type,omitempty"`
	ClientId     string          `json:"client_id,omitempty"`
	Data         json.RawMessage `json:"data,omitempty"`
	IsBroadcast  bool            `json:"is_broadcast"`
}
