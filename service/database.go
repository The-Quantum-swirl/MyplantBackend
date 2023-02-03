package service

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DBConnector struct {
	// db gorm.DB
	DB *sql.DB
}

var connectionName string
var dbuser string
var dbpassword string
var dbname string
var dbport string

func (c *DBConnector) Start() {
	connectionName = os.Getenv("INSTANCE_CONNECTION_NAME")
	dbuser = os.Getenv("DB_USER")
	dbpassword = os.Getenv("DB_PASS")
	dbname = os.Getenv("DB_NAME")
	dbport = os.Getenv("DB_PORT")

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", connectionName, dbport, dbuser, dbpassword, dbname)
	// connString := "postgresql://postgres:postgres@localhost/core-service?sslmode=disable"
	fmt.Print(connString)

	var err error

	c.DB, err = sql.Open("postgres", connString)
	if err != nil {
		fmt.Println("Error in Initialising DB")
	}
	err = c.DB.Ping()
	if err != nil {
		fmt.Println("Error in connecting to DB") // do something here
	} else {
		fmt.Println("Connected Successfully to DB")
	}

	// c.db, err = gorm.Open(postgres.New(postgres.Config{
	// 	DriverName: "cloudsqlpostgres",
	// 	DSN:        connString,
	// }))

	// Connect to the Postgres database on Google Cloud SQL

	if err != nil {
		fmt.Println(err)
	}

}

func (c *DBConnector) HandleRegisterFromNodeDb(email, clientId string) error {
	fmt.Println("ionside db handle rgister")
	fmt.Printf("Email: %s || Client ID : %s", email, clientId)
	res, err := c.DB.Exec("UPDATE public.user SET device_id = $1 WHERE email = $2", clientId, email)

	if err != nil {
		return fmt.Errorf("error inserting data into the database: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("Number of Rows affected are :", rowsAffected)

	if rowsAffected == 0 {
		fmt.Printf("Addign new user as Email not found: %s", email)
		c.saveNewUser(clientId, email)
		return fmt.Errorf("email not found: %s", email)
	}

	return nil
}

func (c *DBConnector) saveNewUser(clientId, email string) error {
	dummyUser := user{
		Email:           email,
		DeviceId:        clientId,
		DeviceName:      "",
		DeviceType:      "",
		FirstName:       "",
		LastName:        "",
		UpdatedAt:       time.Now(),
		ProfilePhotoUrl: "https://example.com/john.doe.jpg",
		Registered:      false,
		MobileNumber:    "555-555-5555",
	}
	saveUser, err := c.DB.Prepare(`INSERT INTO public.user (email, device_id, device_type, first_name, last_name, updated_at, profile_photo_url, registered, mobile_number,device_name)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING id`)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer saveUser.Close()
	_, err1 := saveUser.Exec(dummyUser.Email, dummyUser.DeviceId, dummyUser.DeviceType, dummyUser.FirstName, dummyUser.LastName, dummyUser.UpdatedAt, dummyUser.ProfilePhotoUrl, dummyUser.Registered, dummyUser.MobileNumber, dummyUser.DeviceName)
	if err1 != nil {
		fmt.Println(err)
	}
	fmt.Print("User Saved Successfully")
	return nil
}

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
