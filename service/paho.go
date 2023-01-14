package service

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var knt int

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())
	fmt.Printf("this is result msg #%d!", knt)
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
}

type MQTTConnector struct {
	Client MQTT.Client
	SubCh  string
}

const broker string = "tls://cf4585ba36124b2cbc748a34dce5431d.s1.eu.hivemq.cloud"
const port int = 8883
const username string = "hivemq.webclient.1668925004357"
const password string = "0A:1yRZa2V?xu&7JlKv*"

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
