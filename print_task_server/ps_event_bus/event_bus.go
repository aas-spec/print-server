package ps_event_bus

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"strings"
	"time"
)

/*
	Структура которая будет осуществлять обмен сообщениями между бекендом и фронтом через веб-сокет
*/
type EventBus struct {
	Clients        map[*EventBusClient]bool
	EventListeners map[string]map[*EventBusClient]bool
	Register       chan *EventBusClient
	Unregister     chan *EventBusClient
	RedisClient    *redis.Client
	RedisChannel   string // Topic, на который надо подписываться
}

/*
	Return new event bus
*/
func NewEventBus(redisClient *redis.Client, redisChannel string) *EventBus {

	// Check redis connection and write warning message
	if err := redisClient.Ping().Err(); err != nil {
		log.Printf("unable connect to redis %s", err)
	}

	return &EventBus{
		Clients:        make(map[*EventBusClient]bool),
		Register:       make(chan *EventBusClient),
		Unregister:     make(chan *EventBusClient),
		EventListeners: make(map[string]map[*EventBusClient]bool),
		RedisClient:    redisClient,
		RedisChannel:   redisChannel,
	}
}

/*
	Прослушиваем Redis очередь для получения новых сообщений
*/

func (item *EventBus) ListenRedisQueue() {
	topic := item.RedisClient.PSubscribe(item.RedisChannel)
	defer topic.Close()

	_, err := topic.Receive()
	if err != nil {
		log.Printf("redis pub sub failure %s", err)
	}
	ch := topic.Channel()
	for msg := range ch {
		message := EventBusMessage{}
		// Получаю client id
		pathElements := strings.Split(msg.Channel, "/")
		clientId := pathElements[len(pathElements)-1]
		log.Printf("New redis message to %s: %v", clientId, msg)
		err = json.Unmarshal([]byte(msg.Payload), &message)
		if err != nil {
			log.Printf("unable to unmarshal message %s", err)
			continue
		}
		// Получаю client id
		message.ClientId = clientId
		log.Printf("Start sending message to client %s", message.ClientId)
		item.SendMessage(&message, nil)
	}
}

/*
	Подписаться на событие
*/
func (item *EventBus) SubscribeEvent(eventType string, client *EventBusClient) {
	if _, ok := item.EventListeners[eventType]; !ok {
		item.EventListeners[eventType] = make(map[*EventBusClient]bool)
	}
	item.EventListeners[eventType][client] = true
}

/*
	Отписаться от события
*/
func (item *EventBus) UnsubscribeEvent(eventType string, client *EventBusClient) {
	if _, ok := item.EventListeners[eventType]; ok {
		if _, ok := item.EventListeners[eventType][client]; ok {
			delete(item.EventListeners[eventType], client)
		}
	}
}

/*
	SendMessage сообщение
*/
func (item *EventBus) SendMessage(message *EventBusMessage, startEventClient *EventBusClient) {
	for client := range item.Clients {
		if client != nil && client != startEventClient {
			if message.IsBroadcast == true {
				client.SendMessage(message)
			} else {
				if client.ClientId == message.ClientId { // Шлю клиенту, которому предназначалось сообщение
					client.SendMessage(message)
				}
			}
		}
	}
}

/*
 Ping - раз в 30 секунд шлем сообщение всем клиентам, что бы поддерживать коннект
*/
func (item *EventBus) Ping() {
	// Для отладки - эмуляция сообщений печати
	// pingMessage := GetTestPrintMessage("Test")
	// pingMessage := GetTestDocumentMessage("Test")
	pingMessage := EventBusMessage{CommandType: "PING", IsBroadcast: true}
	for {
		item.SendMessage(&pingMessage, nil)
		time.Sleep(30 * time.Second)
	}
}

func (item *EventBus) Run() {
	go item.ListenRedisQueue()
	go item.Ping()

	for {
		select {
		case eventBusClient := <-item.Register:
			log.Printf("client %s registered", eventBusClient.ClientId)
			item.Clients[eventBusClient] = true
		case eventBusClient := <-item.Unregister:
			if _, ok := item.Clients[eventBusClient]; ok {
				log.Printf("client %s unregistered", eventBusClient.ClientId)
				delete(item.Clients, eventBusClient)
			}
		}

	}
}
