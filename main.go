package main

import (
	"MYPLANTBACKEND/service"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	Email           string `json:"email"`
	DeviceId        string `json:"deviceId"`
	DeviceType      string `json:"deviceType"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	UpdatedAt       string `json:"updatedAt"`
	ProfilePhotoUrl string `json:"profilePhotoUrl"`
	Registered      bool   `json:"registered"`
	MobileNumber    string `json:"mobileNumber"`
}

// var todos = []todo{
// 	{ID: "1", Item: "Clean Room", Completed: false},
// 	{ID: "2", Item: "Read Book", Completed: true},
// }

func getTodos(context *gin.Context, DB *sql.DB) {
	var res user
	var todos []user
	rows, err := DB.Query("SELECT * FROM public.user")
	if err != nil {
		fmt.Println(err)
		context.IndentedJSON(http.StatusBadGateway, "An error occured")
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&res)
		todos = append(todos, res)
	}
	context.IndentedJSON(http.StatusOK, todos)
}

func main() {

	MqttCon := &service.MQTTConnector{Client: nil, SubCh: "register-service"}
	MqttCon.Start()

	// DbCon := &service.DBConnector{DB: nil}
	// DbCon.Start()

	router := gin.Default()
	// router.GET("/todos", func(context *gin.Context) {
	// 	getTodos(context, DbCon.DB)
	// })
	router.GET("/", func(context *gin.Context) {
		message := "Hello, World!"
		MqttCon.Client.Publish("publish-service", 0, false, message)
	})
	router.Run(":8080")
}
