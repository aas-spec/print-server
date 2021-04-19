package ps_event_bus

import "encoding/json"

type DocumentDataUrl struct {
	Url string `json:"url"`
}

type DocumentsMessage struct {
	Documents []DocumentDataUrl `json:"documents"`
}

var testDocumentMessage = DocumentsMessage{
	Documents: []DocumentDataUrl{
		{
			Url: "c:\\temp\\1.pdf",
		},
		{
			Url: "c:\\temp\\2.pdf",
		},
	},
}

func GetTestDocumentMessage(clientId string) EventBusMessage {

	data, _ := json.Marshal(&testDocumentMessage)
	return EventBusMessage{
		CommandType:  "documents",
		CommandParam: "",
		ClientId:     clientId,
		Data:         data,
		IsBroadcast:  false,
	}
}
