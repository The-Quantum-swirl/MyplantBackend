package main

import (
	"MYPLANTBACKEND/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	MqttCon := &service.MQTTConnector{Client: nil, SubCh: "register-service"}
	MqttCon.Start()

	router := gin.Default()
	router.GET("/todos", getTodos)
	router.Run("localhost:8080")
}
