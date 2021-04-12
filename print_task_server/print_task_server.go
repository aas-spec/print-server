package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"print_task_server/ps_event_bus"
	"print_task_server/ws_server"
	"strconv"
)

var redisAddr = flag.String("redis_addr", "localhost", "redis address")
var redisPort = flag.Int("redis_port", 6379, "redis port")
var redisPwd = flag.String("redis_pwd", "", "redis password")
var redisDb = flag.Int("redis_db", 0, "redis database")
var redisChannel = flag.String("redis_channel", "queue/print/*", "print tasks channel")
var webPort = flag.Int("web_port", 12345, "web port")
var webEndpoint = flag.String("web_endpoint", "/queue/print/:key", "web endpoint")

var redisClient *redis.Client

func newRedisClient() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     *redisAddr + ":" + strconv.Itoa(*redisPort),
		Password: *redisPwd, // password set
		DB:       *redisDb,  // use default DB
	})
	return redisClient
}

func main() {
	log.Print("Print task server started")
	flag.Parse()
	log.Printf( "Redis Address: %s:%d", *redisAddr, *redisPort)
	log.Printf( "Redis Db: %d", *redisDb)
	log.Printf( "Redis channel %s", *redisChannel)
	log.Printf( "Web Endpoint %s", *webEndpoint)

	router := gin.Default()
	eventBus := ps_event_bus.NewEventBus(newRedisClient(), *redisChannel)
	go eventBus.Run()
	router.GET(*webEndpoint, func(c *gin.Context) {
		ws_server.PrintWebSocketHandler(eventBus, c)
	})
	log.Printf("Start server on port :%d", *webPort)
	err := router.Run(":" + strconv.Itoa(*webPort))
	if err != nil {
		log.Fatalf("unable to run router %s", err)
	}
}
