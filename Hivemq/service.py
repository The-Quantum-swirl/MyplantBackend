import time
import paho.mqtt.client as paho
from paho import mqtt
from decouple import config
import json

from Hivemq.DBservice import DBService

dbService = DBService()


# setting callbacks for different events to see if it works, print the message etc.
def on_connect(client, userdata, flags, rc, properties=None):
    print("CONNACK received with code %s." % rc)

# with this callback you can see if your publish was successful
def on_publish(client, userdata, mid, properties=None):
    print("mid: " + str(mid))

# print which topic was subscribed to
def on_subscribe(client, userdata, mid, granted_qos, properties=None):
    print("Subscribed: " + str(mid) + " " + str(granted_qos))

# print message, useful for checking if it was successful
def on_message(client, userdata, msg):
    print(msg.topic + " " + str(msg.qos) + " " + str(msg.payload))

    if(msg.topic == config('REGISTER')):
        print("new device incoming")

        json_str = msg.payload.decode('utf-8')
        json_data = json.loads(json_str)

        if("clientId" in json_data.keys() and "email" in json_data.keys()):
            dbService.insert_data(json_data["email"], json_data["clientId"])

        print(str(json_data))


class MqttClient():
    def __init__(self):
        # using MQTT version 5 here, for 3.1.1: MQTTv311, 3.1: MQTTv31
        # userdata is user defined data of any type, updated by user_data_set()
        # client_id is the given name of the client
        self.client = paho.Client(client_id="", userdata=None, protocol=paho.MQTTv5)
        self.client.on_connect = on_connect

        # enable TLS for secure connection
        self.client.tls_set(tls_version=mqtt.client.ssl.PROTOCOL_TLS)
        # set username and password
        # print("username : " + config('MQTT_USER_NAME'))
        # print("pass : " + config('MQTT_PASSWORD'))

        self.client.username_pw_set(config('MQTT_USER_NAME'),config('MQTT_PASSWORD'))
        # connect to HiveMQ Cloud on port 8883 (default for MQTT)
        
        print("attempting to connect to server url : " + config('MQTT_SERVERURL'))
        self.client.connect(config('MQTT_SERVERURL'), 8883)

        # setting callbacks, use separate functions like above for better visibility
        self.client.on_subscribe = on_subscribe
        self.client.on_message = on_message
        self.client.on_publish = on_publish
    
    def listenForDevices(self):
        print("subscribing to topic: " + config('REGISTER'))
        # subscribe to all topics of encyclopedia by using the wildcard "#"
        self.client.subscribe(config('REGISTER'), qos=1)
        

        # loop_forever for simplicity, here you need to stop the loop manually
        # you can also use loop_start and loop_stop
        # self.client.loop_forever()
        self.client.loop_start();
    