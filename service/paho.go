package service

import (
	"fmt"
	"net/http"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var knt int

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())
	fmt.Printf("this is result msg #%d!", knt)
	message := fmt.Sprintf("Message: %s", msg.Payload())
	httpReq(message)
	knt++
}

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Println("Connected")
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
	time.Sleep(2 * time.Second)
	client.Connect()
}

type MQTTConnector struct {
	Client MQTT.Client
	SubCh  string
}

const broker string = "tls://fc61e06e9fda466eb883fa570fe337d4.s1.eu.hivemq.cloud"
const port int = 8883
const username string = "QuantumWaterBot"
const password string = "Quantum#123"
const baseUrl string = "https://api.telegram.org/bot1638003720:AAG1JD9I4XjQYEkYiUTa7An3rOGiVk9sq4M/sendMessage?chat_id=-568647766&text="

func (c *MQTTConnector) Start() {
	fmt.Println("MQTTConnector.start()")

	knt = 0
	// configure the mqtt client
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", broker, port))
	opts.SetClientID("backend")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnectionLost = connectLostHandler
	opts.OnConnect = func(cl MQTT.Client) {
		// on connect will subscribe to default topic
		if token := cl.Subscribe(c.SubCh, 0, f); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	c.Client = MQTT.NewClient(opts)
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}

	// start the connection routine
	fmt.Printf("MQTTConnector.start() Will connect to the broker %v\n", broker)
}
func httpReq(message string) {
	url1 := fmt.Sprintf(baseUrl + message)
	req, _ := http.NewRequest("GET", url1, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}
