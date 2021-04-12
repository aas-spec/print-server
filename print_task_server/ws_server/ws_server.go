package ws_server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"print_task_server/ps_event_bus"
)


var webSocketUpgrader = websocket.Upgrader{

	ReadBufferSize: 1024,

	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func PrintWebSocketHandler(eventBus *ps_event_bus.EventBus, c *gin.Context) {
	conn, err := webSocketUpgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Printf("failed to set websocket upgrade: %+v", err)
		return
	}

	clientId := c.Param("key")

	log.Printf("new client connected, id: %s", clientId)

	client := &ps_event_bus.EventBusClient{

		EventBus: eventBus,

		ClientId: clientId,

		WebSocketConnection: conn,

		Send: make(chan ps_event_bus.EventBusMessage),
	}

	client.EventBus.Register <- client

	client.EventBus.SubscribeEvent("sticker", client)
	client.EventBus.SubscribeEvent("document", client)

	go client.Read()

	go client.Write()

}
