package main

import (
	"MYPLANTBACKEND/common"
	"MYPLANTBACKEND/model"
	"MYPLANTBACKEND/service"
	"log"
	"net/http"
	"strings"

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

	userToBeSaved := model.NewUser(strings.ToLower(requestBody.Email))
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

func findDeviceStatus(c *gin.Context, mqttCon *service.MQTTConnector, dbService *service.DBConnector, ID string) {
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

	result := mqttCon.CheckStatus(res.GetDeviceId())

	c.IndentedJSON(http.StatusBadGateway, gin.H{
		"code":   http.StatusOK,
		"online": result,
	})
}

func saveClientSecret(c *gin.Context, dbService *service.DBConnector, ID string, ClientUserId string) {
	if !common.IsValidUUID(ID) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadGateway,
			"message": ID + " is not a valid Id",
		})
		return
	}

	res := dbService.UpdateClientUserId(&ID, &ClientUserId)
	log.Output(1, res)
	if res == "updated successfully" {
		reusableResponse(c, "Client User ID updated Successfully", 1)
	} else {
		reusableResponse(c, "Failed to update Client USer Id", 0)
	}
}

func reusableResponse(c *gin.Context, message string, selector int) {
	if selector == 0 {
		c.IndentedJSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": message,
		})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"user": message,
		})
	}

}

func main() {

	// gin.SetMode(gin.ReleaseMode)

	// setting db connection
	DbCon := &service.DBConnector{DB: nil}
	DbCon.Start()

	// setting mqtt connection
	MqttCon := &service.MQTTConnector{Client: nil, SubCh: "register-service", NotifCh: "notification-service", DBCon: DbCon}
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
		email := strings.ToLower(context.Param("email"))
		findUserByEmailId(context, DbCon, email)
	})

	router.GET("user/getById/:id", func(context *gin.Context) {
		ID := context.Param("id")
		findUserById(context, DbCon, ID)
	})

	router.GET("user/online/:id", func(context *gin.Context) {
		ID := context.Param("id")
		findDeviceStatus(context, MqttCon, DbCon, ID)
	})

	router.POST("user/save", func(context *gin.Context) {
		saveUserDetails(context, DbCon)
	})

	router.GET("user/notification/:ClientUserId/:id", func(context *gin.Context) {
		ID := context.Param("id")
		ClientUserId := context.Param("ClientUserId")
		saveClientSecret(context, DbCon, ID, ClientUserId)
	})

	router.Run(":8080")
}
