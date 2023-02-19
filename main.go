package main

import (
	"MYPLANTBACKEND/common"
	"MYPLANTBACKEND/model"
	"MYPLANTBACKEND/service"
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

func getTodos(c *gin.Context, dbService *service.DBConnector) {
	res := dbService.GetAllUser()
	log.Output(2, "get all user call executed")

	if res == nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Unable to fetch Users",
		})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": res,
		})
	}
}

func saveUserDetails(c *gin.Context, dbService *service.DBConnector) {
	var requestBody UserRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		log.Output(1, err.Error())
	}

	userToBeSaved := model.NewUser(requestBody.Email)
	userToBeSaved.SetName(requestBody.FirstName, requestBody.LastName)
	userToBeSaved.SetProfilePhoto(requestBody.PhotoUrl)
	userToBeSaved.RegisterIt()

	if dbService.SaveNewUser(userToBeSaved) != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"userId":  userToBeSaved.GetId(),
			"message": "Saved",
		})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Unable to update",
		})
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

func findUserById(c *gin.Context, dbService *service.DBConnector, ID string) {
	// validation
	if !common.IsValidUUID(ID) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadGateway,
			"message": ID + " is not a valid Id",
		})
		return
	}
	// search for user
	res := dbService.GetUserByID(&ID)

	if res == nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "An error occured in Finding User for ID : " + ID,
		})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"user": res,
		})
	}
}

func main() {

	// gin.SetMode(gin.ReleaseMode)

	// setting db connection
	DbCon := &service.DBConnector{DB: nil}
	DbCon.Start()

	// setting mqtt connection
	MqttCon := &service.MQTTConnector{Client: nil, SubCh: "register-service", DBCon: DbCon}
	MqttCon.Start()

	// setting router
	router := gin.Default()

	//paths
	router.GET("user/getAll", func(context *gin.Context) {
		getTodos(context, DbCon)
	})

	router.GET("/publishTest", func(context *gin.Context) {
		message := "Hello, World!"
		MqttCon.Client.Publish("publish-service", 0, false, message)
		context.IndentedJSON(http.StatusOK, "published")
	})

	router.GET("user/getByEmail/:email", func(context *gin.Context) {
		email := context.Param("email")
		findUserByEmailId(context, DbCon, email)
	})

	router.GET("user/getById/:id", func(context *gin.Context) {
		ID := context.Param("id")
		findUserById(context, DbCon, ID)
	})

	router.POST("user/save", func(context *gin.Context) {
		saveUserDetails(context, DbCon)
	})

	router.Run(":8080")
}
