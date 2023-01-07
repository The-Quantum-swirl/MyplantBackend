from django.http import HttpResponse
from Hivemq.service import MqttClient;

client = MqttClient()
client.listenForDevices()

def index(request):
    return HttpResponse("Hello, world. You're at the polls index.")

def getRegisteredDevices(request):
    return HttpResponse("Hello, world")
    