package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"MYPLANTBACKEND/model"
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
	// connString := "postgresql://postgres:quicuxeo@localhost/core-service?sslmode=disable"
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
		newUser := model.NewUser(email)
		newUser.SetDevice(clientId,"Android")
		c.SaveNewUser(newUser)
		return fmt.Errorf("email not found: %s", email)
	}

		defer c.DB.Close()
		return nil
	}

// fetch user from db
func (c *DBConnector) GetAllUser() []*model.User {
	query := `SELECT * FROM public.user`
	rows, err := c.DB.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var res user
	defer rows.Close()

	var UserList []*model.User
	for rows.Next() {
		rows.Scan(&res.ID, &res.Email, &res.DeviceId, &res.DeviceType, &res.FirstName, &res.LastName, &res.UpdatedAt, &res.ProfilePhotoUrl, &res.Registered, &res.MobileNumber, &res.DeviceName)
		// initialize new user
		newUser := model.NewUser(res.Email)
		newUser.SetDevice(res.DeviceId, "wp")
		newUser.SetMobileNumber(res.MobileNumber)
		newUser.SetName(res.FirstName, res.LastName)
		newUser.SetProfilePhoto(res.ProfilePhotoUrl)
		// return the first user found
		UserList = append(UserList, newUser)
	}
	return UserList
}

// fetch user from db
func (c *DBConnector) GetUser(email *string) *model.User {
	query := `SELECT * FROM public.user WHERE email = $1`
	rows, err := c.DB.Query(query, email)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var res user
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&res.ID, &res.Email, &res.DeviceId, &res.DeviceType, &res.FirstName, &res.LastName, &res.UpdatedAt, &res.ProfilePhotoUrl, &res.Registered, &res.MobileNumber, &res.DeviceName)
		// initialize new user
		newUser := model.NewUser(res.Email)
		newUser.SetDevice(res.DeviceId, "wp")
		newUser.SetMobileNumber(res.MobileNumber)
		newUser.SetName(res.FirstName, res.LastName)
		newUser.SetProfilePhoto(res.ProfilePhotoUrl)
		// return the first user found
		return newUser
	}
	log.Output(1, "user not found")
	return nil
}

// save user in db
func (c *DBConnector) SaveNewUser(u *model.User) *model.User {
	query := `INSERT INTO public.user (email, device_id, device_type, first_name, last_name, updated_at, profile_photo_url, registered, mobile_number)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id`
	_, err := c.DB.Exec(query, (*u).GetEmail(), (*u).GetDeviceId(), (*u).GetDeviceId(), (*u).GetFirstName(), (*u).GetLastName(), time.Now(), (*u).GetProfilePhoto(), (*u).IsRegistered(), (*u).GetMobileNumber())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	log.Output(1, "User Added Sucessfully")
	return u
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
