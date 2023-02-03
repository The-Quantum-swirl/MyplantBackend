package main

import (
	"MYPLANTBACKEND/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID              int       `json:"id"`
	Email           string    `json:"email"`
	DeviceId        string    `json:"deviceId"`
	DeviceType      string    `json:"deviceType"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	UpdatedAt       time.Time `json:"updatedAt"`
	ProfilePhotoUrl string    `json:"profilePhotoUrl"`
	Registered      bool      `json:"registered"`
	MobileNumber    string    `json:"mobileNumber"`
	DeviceName      string    `json:"deviceName"`
}

// var todos = []todo{
// 	{ID: "1", Item: "Clean Room", Completed: false},
// 	{ID: "2", Item: "Read Book", Completed: true},
// }

func getTodos(context *gin.Context, DB *sql.DB) {
	var res user
	var todos []user
	rows, err := DB.Query(`SELECT * FROM public.user`)
	if err != nil {
		fmt.Println(err)
		context.IndentedJSON(http.StatusBadGateway, "An error occured")
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&res.ID, &res.Email, &res.DeviceId, &res.DeviceType, &res.FirstName, &res.LastName, &res.UpdatedAt, &res.ProfilePhotoUrl, &res.Registered, &res.MobileNumber, &res.DeviceName)
		fmt.Println(res)
		todos = append(todos, res)
	}
	context.IndentedJSON(http.StatusOK, todos)
}
func findUserByEmailId(context *gin.Context, stmt *sql.Stmt, email string) {
	// email := context.Param("email")
	var res user
	//one liner//
	err := stmt.QueryRow(email).Scan(&res.ID, &res.Email)
	// var todos []user

	// rows, err := stmt.Query(email)
	if err != nil {
		fmt.Println(err)
		context.IndentedJSON(http.StatusBadGateway, "An error occured in Finding User : "+email)
	} else {
		context.IndentedJSON(http.StatusOK, res)
	}
}
func saveUserFunc(context *gin.Context, Db *sql.DB, stmt *sql.Stmt, dummyUser user) {
	_, err := stmt.Exec(dummyUser.Email, dummyUser.DeviceId, dummyUser.DeviceType, dummyUser.FirstName, dummyUser.LastName, dummyUser.UpdatedAt, dummyUser.ProfilePhotoUrl, dummyUser.Registered, dummyUser.MobileNumber, dummyUser.DeviceName)
	if err != nil {
		fmt.Println(err)
		context.IndentedJSON(http.StatusBadGateway, "An error occured in saving User")
	}
	var id int
	if err := Db.QueryRow("SELECT public.user.id FROM user WHERE email = $1", dummyUser.Email).Scan(&id); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Dummy user inserted successfully with id:", id)

	context.IndentedJSON(http.StatusOK, "User Saved Successfully")
}
func runQuery(context *gin.Context, Db *sql.DB) {
	var res user

	stmt, err := Db.Prepare(`SELECT * FROM public.user where email = 'john.doe@example.com'`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&res.ID, &res.Email, &res.DeviceId, &res.DeviceType, &res.FirstName, &res.LastName, &res.UpdatedAt, &res.ProfilePhotoUrl, &res.Registered, &res.MobileNumber, &res.DeviceName)
		fmt.Println(res)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
func HandleRegisterFromNode() {

}

func main() {

	DbCon := &service.DBConnector{DB: nil}
	DbCon.Start()

	MqttCon := &service.MQTTConnector{Client: nil, SubCh: "register-service1", DBCon: DbCon}
	MqttCon.Start()
	Db := DbCon.DB
	fetchUserByEmail, err := Db.Prepare(`SELECT id, email FROM public.user where email = $1;`)
	if err != nil {
		log.Fatal(err)
	}
	defer fetchUserByEmail.Close()

	saveUser, err := Db.Prepare(`INSERT INTO public.user (email, device_id, device_type, first_name, last_name, updated_at, profile_photo_url, registered, mobile_number,device_name)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING id`)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer saveUser.Close()

	router := gin.Default()
	router.GET("/todos", func(context *gin.Context) {
		getTodos(context, DbCon.DB)
	})

	router.GET("/publishTest", func(context *gin.Context) {
		message := "Hello, World!"
		MqttCon.Client.Publish("publish-service", 0, false, message)
	})

	router.POST("/saveUser", func(context *gin.Context) {
		dummyUser := user{
			Email:           "john.doe@example.com",
			DeviceId:        "1234567890",
			DeviceName:      "Spitzy",
			DeviceType:      "Android",
			FirstName:       "John",
			LastName:        "Doe",
			UpdatedAt:       time.Now(),
			ProfilePhotoUrl: "https://example.com/john.doe.jpg",
			Registered:      true,
			MobileNumber:    "555-555-5555",
		}
		saveUserFunc(context, Db, saveUser, dummyUser)
	})

	router.GET("/fetchUser/:email", func(context *gin.Context) {
		email := context.Param("email")
		fmt.Println("searching for : " + email)
		findUserByEmailId(context, fetchUserByEmail, email)
	})
	router.GET("/query", func(context *gin.Context) {
		runQuery(context, Db)
	})
	router.Run(":8080")
}
