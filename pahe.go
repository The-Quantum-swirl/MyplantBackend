// package main

// import (
// 	"fmt"
// 	"time"

// 	MQTT "github.com/eclipse/paho.mqtt.golang"
// )

// var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
// 	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
// }

// var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
// 	fmt.Println("Connected")
// }

// var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
// 	fmt.Printf("Connect lost: %v", err)
// }

// type MQTTConnector struct {
// 	client MQTT.Client
// 	pubCh  string
// 	subCh  string
// }

// const defaultQoS = 1
// const broker string = "tls://cf4585ba36124b2cbc748a34dce5431d.s1.eu.hivemq.cloud"
// const port int = 8883

// func (c *MQTTConnector) start() {
// 	fmt.Println("MQTTConnector.start()")

// 	// configure the mqtt client
// 	c.configureMqttConnection()

// 	// start the connection routine
// 	fmt.Printf("MQTTConnector.start() Will connect to the broker %v\n", broker)
// 	go c.connect(0)

// }

// // func (c *MQTTConnector) stop() {
// // 	fmt.Println("MQTTConnector.stop()")
// // 	if c.client != nil && c.client.IsConnected() {
// // 		c.client.Disconnect(500)
// // 	}
// // }

// func (c *MQTTConnector) connect(backOff int) {
// 	if c.client == nil {
// 		fmt.Printf("MQTTConnector.connect() client is not configured")
// 		return
// 	}
// 	for {
// 		fmt.Printf("MQTTConnector.connect() connecting to the broker %v, backOff: %v sec\n", broker, backOff)
// 		time.Sleep(time.Duration(backOff) * time.Second)
// 		if c.client.IsConnected() {
// 			break
// 		}
// 		token := c.client.Connect()
// 		token.Wait()
// 		if token.Error() == nil {
// 			break
// 		}
// 		fmt.Printf("MQTTConnector.connect() failed to connect: %v\n", token.Error().Error())
// 		if backOff == 0 {
// 			backOff = 10
// 		} else if backOff <= 600 {
// 			backOff *= 2
// 		}
// 	}

// 	fmt.Printf("MQTTConnector.connect() connected to the broker %v", broker)
// 	return
// }

// // func (c *MQTTConnector) onConnected(client *MQTT.Client) {
// // 	// subscribe if there is at least one resource with SUB in MQTT protocol is configured
// // 	if len(c.subCh) > 0 {
// // 		fmt.Println("MQTTPulbisher.onConnected() will (re-)subscribe to all configured SUB topics")

// // 		// topicFilters := make(map[string]byte)
// // 		// for topic, _ := range c.subCh {
// // 		// 	fmt.Printf("MQTTPulbisher.onConnected() will subscribe to topic %s", topic)
// // 		// 	topicFilters[topic] = defaultQoS
// // 		// }
// // 		// client.SubscribeMultiple(topicFilters, c.messageHandler)
// // 	} else {
// // 		fmt.Println("MQTTPulbisher.onConnected() no resources with SUB configured")
// // 	}
// // }

// // func (c *MQTTConnector) onConnectionLost(client *MQTT.Client, reason error) {
// // 	fmt.Println("MQTTPulbisher.onConnectionLost() lost connection to the broker: ", reason.Error())

// // 	// Initialize a new client and reconnect
// // 	c.configureMqttConnection()
// // 	go c.connect(0)
// // }

// func (c *MQTTConnector) configureMqttConnection() {
// 	opts := MQTT.NewClientOptions()
// 	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
// 	opts.SetClientID("backend")
// 	opts.SetUsername("hivemq.webclient.1668925004357")
// 	opts.SetPassword("0A:1yRZa2V?xu&7JlKv*")
// 	opts.SetDefaultPublishHandler(messagePubHandler)
// 	opts.OnConnect = connectHandler
// 	opts.OnConnectionLost = connectLostHandler

// 	c.client = MQTT.NewClient(opts)
// }

// // func main() {

// // 	var broker = "tls://cf4585ba36124b2cbc748a34dce5431d.s1.eu.hivemq.cloud" // find the host name in the Overview of your cluster (see readme)
// // 	var port = 8883                                                          // find the port right under the host name, standard is 8883
// // 	opts := mqtt.NewClientOptions()
// // 	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
// // 	opts.SetClientID("backend")
// // 	opts.SetUsername("hivemq.webclient.1668925004357")
// // 	opts.SetPassword("0A:1yRZa2V?xu&7JlKv*")
// // 	opts.SetDefaultPublishHandler(messagePubHandler)
// // 	opts.OnConnect = connectHandler
// // 	opts.OnConnectionLost = connectLostHandler
// // 	client := mqtt.NewClient(opts)
// // 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// // 		panic(token.Error())
// // 	}

// // 	sub(client)
// // 	publish(client)

// // 	client.Disconnect(250)
// // }
