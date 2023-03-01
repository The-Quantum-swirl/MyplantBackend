package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var knt int
var DbMG *DBConnector
var handleMQTTMessage MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Output(1, "MSG: "+string(msg.Payload()))
	log.Println("this is result msg" + strconv.Itoa(knt))
	message := fmt.Sprintf("Message: %s", msg.Payload())
	httpReq(message)
	knt++
}

func processNodeRegistration(msg MQTT.Message) {
	var det Detail
	err := json.Unmarshal([]byte(msg.Payload()), &det)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return
	}
	det.Email = strings.ToLower(det.Email)
	DbMG.HandleRegisterFromNodeDb(det.Email, det.ClientId)
}

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	if strings.Compare("register-service", msg.Topic()) == 0 {
		log.Println("saving Device Id to DB")
		processNodeRegistration(msg)
	} else if strings.Compare("notification-service", msg.Topic()) == 0 {
		log.Println("got message in notification channel")
	}
}

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	log.Println("Connected")
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	log.Println("Connect lost: ", err)
	client.Connect()
}

type MQTTConnector struct {
	Client  MQTT.Client
	SubCh   string
	NotifCh string
	DBCon   *DBConnector
}
type Detail struct {
	Email    string `json:"email"`
	ClientId string `json:"clientId"`
}

const broker string = "tls://fc61e06e9fda466eb883fa570fe337d4.s1.eu.hivemq.cloud"
const port int = 8883
const username string = "QuantumWaterBot"
const password string = "Quantum#123"
const baseUrl string = "https://api.telegram.org/bot1638003720:AAG1JD9I4XjQYEkYiUTa7An3rOGiVk9sq4M/sendMessage?chat_id=-568647766&text="

var connection string = "backend1"

func (c *MQTTConnector) Start() {

	time.Sleep(3 * time.Second)
	log.Println("---------- MQTT started ----------")
	connectionId := os.Getenv("DB_ID")

	knt = 0
	// configure the mqtt client
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", broker, port))
	opts.SetClientID(connectionId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnectionLost = connectLostHandler
	opts.OnConnect = func(cl MQTT.Client) {
		// on connect will subscribe to default topic
		if token := cl.Subscribe(c.SubCh, 0, messagePubHandler); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		// on connect will subscribe to notification topic
		if token := cl.Subscribe(c.NotifCh, 0, messagePubHandler); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	c.Client = MQTT.NewClient(opts)
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		log.Printf("Connected to server\n")
	}
	DbMG = c.DBCon
	// start the connection routine
	log.Printf("MQTT Will connect to the broker %v\n", broker)
}
func httpReq(message string) {
	url1 := fmt.Sprintf(baseUrl + message)
	req, _ := http.NewRequest("GET", url1, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
}
