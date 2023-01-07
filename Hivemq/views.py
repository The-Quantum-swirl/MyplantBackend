from django.http import HttpResponse
from django.http import JsonResponse
import json

from Hivemq.service import MqttClient;
from Hivemq.DBservice import DBService;


client = MqttClient()
client.listenForDevices()

dbService = DBService()

def index(request):
    return HttpResponse("Hello, world. You're at index page.")

def getDevices(request):

    body = json.loads(request.body)
    email = body['email']
    print(email);
    
    res = dbService.fetchDeviceId(email);

    print(res[1]);
    
    return JsonResponse({'deviceId':res[1]})
    