package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DBConnector struct {
	DB *sql.DB
}

var connectionName string
var dbuser string
var dbpassword string
var dbname string
var dbport string

var saveUserStmt *sql.Stmt

func init() {
	log.Output(1, "Init main Called ")
}
func (c *DBConnector) init() {
	log.Output(1, "Init called")
	var err error
	saveUserStmt, err = c.prepareSaveNewUserStmt()
	if err != nil {
		fmt.Println(err)
	}
}
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
		fmt.Errorf("error initializing the database: %v", err)
	}

	err = c.DB.Ping()
	if err != nil {
		fmt.Errorf("error connecting to the database: %v", err)
	}
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

	defer c.DB.Close()
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
	query := `INSERT INTO public.user (email, device_id, device_type, first_name, last_name, updated_at, profile_photo_url, registered, mobile_number,device_name)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING id`
	_, err := c.DB.Exec(query, dummyUser.Email, dummyUser.DeviceId, dummyUser.DeviceType, dummyUser.FirstName, dummyUser.LastName, dummyUser.UpdatedAt, dummyUser.ProfilePhotoUrl, dummyUser.Registered, dummyUser.MobileNumber, dummyUser.DeviceName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("User Added Sucessfully")
	return nil
}
func (c *DBConnector) prepareSaveNewUserStmt() (*sql.Stmt, error) {
	saveUser, err := c.DB.Prepare(`INSERT INTO public.user (email, device_id, device_type, first_name, last_name, updated_at, profile_photo_url, registered, mobile_number,device_name)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING id`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return saveUser, nil
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
