package main

import (
	"MYPLANTBACKEND/model"
	"MYPLANTBACKEND/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRequestBody struct {
	Email     string
	PhotoUrl  string
	FirstName string
	LastName  string
}

func getTodos(context *gin.Context, dbService *service.DBConnector) {
	res := dbService.GetAllUser()
	log.Output(2, "get all user call executed")

	if res == nil {
		context.IndentedJSON(http.StatusBadGateway, "An error occured in Fetching User")
	} else {
		context.IndentedJSON(http.StatusOK, res)
	}
}

func saveUserDetails(context *gin.Context, dbService *service.DBConnector) {
	var requestBody UserRequestBody
	if err := context.BindJSON(&requestBody); err != nil {
		fmt.Println(err)
	}

	userToBeSaved := model.NewUser(requestBody.Email)
	userToBeSaved.SetName(requestBody.FirstName, requestBody.LastName)
	userToBeSaved.SetProfilePhoto(requestBody.PhotoUrl)
	userToBeSaved.RegisterIt()

	if dbService.SaveNewUser(userToBeSaved) != nil {
		context.IndentedJSON(http.StatusOK, "Saved")
	} else {
		context.IndentedJSON(http.StatusAlreadyReported, "User already there")
	}
}

func findUserByEmailId(context *gin.Context, dbService *service.DBConnector, email string) {
	res := dbService.GetUser(&email)

	if res == nil {
		context.IndentedJSON(http.StatusBadGateway, "An error occured in Finding User : "+email)
	} else {
		context.IndentedJSON(http.StatusOK, res)
	}
}

func main() {
	// setting db connection
	DbCon := &service.DBConnector{DB: nil}
	DbCon.Start()

	// setting mqtt connection
	MqttCon := &service.MQTTConnector{Client: nil, SubCh: "register-service1", DBCon: DbCon}
	MqttCon.Start()

	// setting router
	router := gin.Default()

	//paths
	router.GET("/todos", func(context *gin.Context) {
		getTodos(context, DbCon)
	})

	router.GET("/publishTest", func(context *gin.Context) {
		message := "Hello, World!"
		MqttCon.Client.Publish("publish-service", 0, false, message)
	})

	router.GET("/fetchUser/:email", func(context *gin.Context) {
		email := context.Param("email")
		findUserByEmailId(context, DbCon, email)
	})

	router.POST("/saveUserDetails", func(context *gin.Context) {
		saveUserDetails(context, DbCon)
	})

	router.Run(":8080")
}
