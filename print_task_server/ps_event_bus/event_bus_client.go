package ps_event_bus

import (
	"github.com/gorilla/websocket"
	"log"
)

/*
	Структура обрабатывает запросы конкретного клиента
*/
type EventBusClient struct {
	EventBus            *EventBus
	ClientId            string
	WebSocketConnection *websocket.Conn
	Send                chan EventBusMessage
}

func (item *EventBusClient) Read() {
	defer func() {
		item.EventBus.Unregister <- item
		item.WebSocketConnection.Close()
		log.Printf("Client %s disconnected", item.ClientId)
	}()

	for {
		message := EventBusMessage{}
		err := item.WebSocketConnection.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("web socket read error: %v", err)
			}
			break
		}

		item.ProcessCommand(&message)
	}
}

func (item *EventBusClient) ProcessCommand(message *EventBusMessage) {
	switch message.CommandType {
	case "subscribeEvent":
		item.EventBus.SubscribeEvent(message.CommandParam, item)
	case "unsubscribeEvent":
		item.EventBus.UnsubscribeEvent(message.CommandParam, item)
	case "sendMessage":
		item.EventBus.SendMessage(message, item)
	}
}

func (item *EventBusClient) SendMessage(message *EventBusMessage) {
	item.Send <- *message
}

func (item *EventBusClient) Write() {
	defer func() {
		item.WebSocketConnection.Close()
	}()

	for {
		select {
		case message, ok := <-item.Send:
			if !ok {
				// The hub closed the channel.
				err := item.WebSocketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("unable to write message for client %s: %s", item.ClientId, err)
					return
				}
			}

			err := item.WebSocketConnection.WriteJSON(&message)
			if err != nil {
				log.Printf("unable to write json for client %s: %s", item.ClientId, err)
				return
			}
		}
	}

}
