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
	httpReq("www.example.com", message)
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
func sendNotification(msg MQTT.Message) {
	var notif Notification
	err := json.Unmarshal([]byte(msg.Payload()), &notif)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return
	}
	clientUserId := DbMG.getClientUserId(&notif.ClientId)
	if clientUserId != "failed to fetch client User ID" {
		finalurl := strings.Replace(baseUrl, "clientUserIDPlaceholder", clientUserId, -1)
		if notif.Status == "off" {
			httpReq(finalurl, "Water Turned OFF")

		}else if notif.Status == "started"{
			httpReq(finalurl, "Watering Device Started")
		} else {
			httpReq(finalurl, "Water Turned ON")

		}
	} else {
		fmt.Println("Skipping sending notifications")
	}
}

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	if strings.Compare("register-service", msg.Topic()) == 0 {
		log.Println("saving Device Id to DB")
		processNodeRegistration(msg)
	} else if strings.Compare("notification-service", msg.Topic()) == 0 {
		log.Println("got message in notification channel")
		sendNotification(msg)
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
type Notification struct {
	ClientId string `json:clientId`
	Status   string `json:status`
}

const broker string = "tls://fc61e06e9fda466eb883fa570fe337d4.s1.eu.hivemq.cloud"
const port int = 8883
const username string = "QuantumWaterBot"
const password string = "Quantum#123"
const baseUrl string = "https://api.telegram.org/bot6262793721:AAH3Q3dVEXJv2sOHB1b20QxzERiDoZUmsQQ/sendMessage?chat_id=clientUserIDPlaceholder&text="

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
func httpReq(url string, message string) {
	url1 := fmt.Sprintf(url + message)
	req, _ := http.NewRequest("GET", url1, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
}
