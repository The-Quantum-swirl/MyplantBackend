package service

import (
	"MYPLANTBACKEND/model"
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
	log.Output(2, "Init main Called ")
}

func (c *DBConnector) Start() {
	connectionName = os.Getenv("INSTANCE_CONNECTION_NAME")
	dbuser = os.Getenv("DB_USER")
	dbpassword = os.Getenv("DB_PASS")
	dbname = os.Getenv("DB_NAME")
	dbport = os.Getenv("DB_PORT")
	// sslmode=disable
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", connectionName, dbport, dbuser, dbpassword, dbname)
	// connString := "postgresql://postgres:quicuxeo@localhost/core-service?sslmode=disable"
	log.Print(connString)
	var err error

	c.DB, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("error initializing the database: " + err.Error())
	}

	err = c.DB.Ping()
	if err != nil {
		log.Fatal("error connecting to the database: " + err.Error())
	}
	if err != nil {
		log.Println(err)
	}
	// max timeout for db connection
	c.DB.SetConnMaxLifetime(1800 * time.Second)

	// Set maximum number of connections in idle connection pool.
	c.DB.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	c.DB.SetMaxOpenConns(7)
}

func (c *DBConnector) HandleRegisterFromNodeDb(email, clientId string) error {
	log.Println("----- registering device -----")
	log.Printf("Email: %s || Client ID : %s", email, clientId)

	res, err := c.DB.Exec("UPDATE users SET device_id = $1 WHERE email = $2", clientId, email)
	if err != nil {
		return fmt.Errorf("error inserting data into the database: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("Number of Rows affected are :", rowsAffected)

	// if no rows are affected then add new user
	if rowsAffected == 0 {
		log.Printf("Adding new user as Email not found: %s", email)
		newUser := model.NewUser(email)
		newUser.SetDevice(clientId, "wp")
		c.SaveNewUser(newUser)
	}

	defer c.DB.Close()
	return nil
}

// fetch user from db
func (c *DBConnector) GetAllUser() []*model.User {
	query := `SELECT * FROM users`
	rows, err := c.DB.Query(query)
	if err != nil {
		log.Println(err)
		return nil
	}

	var res model.User
	defer rows.Close()

	var UserList []*model.User
	for rows.Next() {
		rows.Scan(&res.ID, &res.Email, &res.DeviceId, &res.DeviceType, &res.FirstName, &res.LastName, &res.UpdatedAt, &res.ProfilePhotoUrl, &res.Registered, &res.MobileNumber)
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
	query := `SELECT * FROM users WHERE email = $1`
	rows, err := c.DB.Query(query, email)
	if err != nil {
		log.Println(err)
		return nil
	}

	var res model.User
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&res.ID, &res.Email, &res.DeviceId, &res.DeviceType, &res.FirstName, &res.LastName, &res.UpdatedAt, &res.ProfilePhotoUrl, &res.Registered, &res.MobileNumber)
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
	query := `INSERT INTO users (email, device_id, device_type, first_name, last_name, updated_at, profile_photo_url, registered, mobile_number)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id`
	_, err := c.DB.Exec(query, (*u).GetEmail(), (*u).GetDeviceId(), (*u).GetDeviceId(), (*u).GetFirstName(), (*u).GetLastName(), time.Now(), (*u).GetProfilePhoto(), (*u).IsRegistered(), (*u).GetMobileNumber())
	if err != nil {
		log.Println(err)

		query := `UPDATE users SET email = $1 first_name = $2 last_name = $3 profile_photo_url = $4 updated_at = $5 registered = $6`
		_, err := c.DB.Exec(query, (*u).GetEmail(), (*u).GetFirstName(), (*u).GetLastName(), (*u).GetProfilePhoto(), time.Now(), true)
		if err != nil {
			log.Println(err)
			return nil
		}
		return u
	}
	log.Output(1, "User Added Sucessfully")
	return u
}
