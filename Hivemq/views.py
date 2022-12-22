from django.http import HttpResponse
from Hivemq.service import MqttClient;

def index(request):
    return HttpResponse("Hello, world. You're at the polls index.")

def getRegisteredDevices(request):
    client = MqttClient()
    client.listenForDevices()
    
    return HttpResponse("Hello, world")
    