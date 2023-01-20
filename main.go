package main

import (
	"MYPLANTBACKEND/service"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type user struct {
	email           string `gorm:"primaryKey"`
	DeviceId        string `json:"deviceId"`
	DeviceType      string `json:"deviceType"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName" gorm:"default:galeone"`
	ProfilePhotoUrl string `json:"profilePhotoUrl"`
	Registered      bool   `json:"registered"`
	MobileNumber    string `json:"mobileNumber"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// var todos = []todo{
// 	{ID: "1", Item: "Clean Room", Completed: false},
// 	{ID: "2", Item: "Read Book", Completed: true},
// }

// func getTodos(context *gin.Context, DB *sql.DB) {
func getTodos(context *gin.Context, DB *gorm.DB) {
	var res user
	var todos []user
	// rows, err := DB.Query("SELECT * FROM public.user")
	result := DB.First(&res)

	// result.RowsAffected // returns count of records found

	if result.Error != nil {
		fmt.Println(result.Error)
		context.IndentedJSON(http.StatusBadGateway, "An error occured")
	}

	// defer rows.Close()
	// for rows.Next() {
	// 	rows.Scan(&res)
	// 	todos = append(todos, res)
	// }
	todos = append(todos, res)
	context.IndentedJSON(http.StatusOK, todos)
}

func adduser(context *gin.Context, DB *gorm.DB) {
	dummy_user := user{
		email:           "example@email.com",
		DeviceId:        "12345",
		DeviceType:      "ios",
		FirstName:       "John",
		LastName:        "Doe",
		ProfilePhotoUrl: "https://example.com/image.jpg",
		Registered:      true,
		MobileNumber:    "1234567890",
	}
	if !DB.Migrator().HasTable("users") {
		fmt.Print("Creating table")
		DB.Migrator().CreateTable(&user{})
	}
	result := DB.Create(&dummy_user)
	errors.Is(result.Error, gorm.ErrInvalidData)
}
func main() {

	MqttCon := &service.MQTTConnector{Client: nil, SubCh: "register-service"}
	MqttCon.Start()

	DbCon := &service.DBConnector{DB: nil}
	DbCon.Start()

	router := gin.Default()
	router.GET("/todos", func(context *gin.Context) {
		getTodos(context, DbCon.DB)
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	router.GET("/adduser", func(context *gin.Context) {
		adduser(context, DbCon.DB)
		context.JSON(http.StatusOK, gin.H{
			"message": "User Craeted Successfully!",
		})
	})

	router.GET("/adduser", func(context *gin.Context) {
			rows := DbCon.DB.Table("information_schema.tables")
			var tables []string
			var name string
			for rows.Next() {
				row.Scan(&name)
				tables = append(tables, name)
			}
context.JSON(http.StatusOK, gin.H{
			"message": "User Craeted Successfully!",
		})
	})

	router.Run("localhost:8080")

}
