package main

import (
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const broker string = "tls://cf4585ba36124b2cbc748a34dce5431d.s1.eu.hivemq.cloud"
const port int = 8883
const username string = "hivemq.webclient.1668925004357"
const password string = "0A:1yRZa2V?xu&7JlKv*"

var knt int
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())
	text := fmt.Sprintf("this is result msg #%d!", knt)
	knt++
	token := client.Publish("nn/result", 0, false, text)
	token.Wait()
}

type todo struct {
	ID        string `json:"id"`
	Item      string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{ID: "1", Item: "Clean Room", Completed: false},
	{ID: "2", Item: "Read Book", Completed: true},
}

func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func main() {
	knt = 0
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fmt.Println(password)

	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("%s:%d", broker, port))
	opts.SetClientID("mac-go")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(f)
	topic := "register-service"

	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}
	// <-c

	router := gin.Default()
	router.GET("/todos", getTodos)
	router.Run("localhost:8080")
}
