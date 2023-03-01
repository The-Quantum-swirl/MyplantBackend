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
	// connString := "postgresql://postgres:postgres@localhost/postgres?sslmode=disable"
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
	c.DB.SetConnMaxLifetime(1200 * time.Second)

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
		newUser.RegisterIt()
		c.RegisterNewUser(newUser)
	}
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
		rows.Scan(&res.ID, &res.Email, &res.DeviceId, &res.DeviceType, &res.FirstName, &res.LastName, &res.UpdatedAt, &res.ProfilePhotoUrl, &res.Registered, &res.MobileNumber, &res.ClientUserId)
		// initialize new user
		newUser := model.NewUser(res.Email)
		newUser.SetId(res.ID)
		newUser.SetDevice(res.DeviceId, res.DeviceType)
		newUser.SetMobileNumber(res.MobileNumber)
		newUser.SetName(res.FirstName, res.LastName)
		newUser.SetProfilePhoto(res.ProfilePhotoUrl)
		log.Println(res.UpdatedAt.GoString())
		newUser.SetUpdatedAt(res.UpdatedAt)
		newUser.SetClientUserId(res.ClientUserId)
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
		newUser.SetId(res.ID)
		newUser.SetDevice(res.DeviceId, res.DeviceType)
		newUser.SetMobileNumber(res.MobileNumber)
		newUser.SetName(res.FirstName, res.LastName)
		newUser.SetProfilePhoto(res.ProfilePhotoUrl)
		newUser.SetUpdatedAt(res.UpdatedAt)
		// return the first user found
		return newUser
	}
	log.Output(1, "user not found")
	return nil
}

// fetch user from db
func (c *DBConnector) GetUserByID(ID *string) *model.User {
	query := `SELECT * FROM users WHERE id = $1`
	rows, err := c.DB.Query(query, ID)
	if err != nil {
		log.Println(err)
		return nil
	}

	var res model.User
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&res.ID, &res.Email, &res.DeviceId, &res.DeviceType, &res.FirstName, &res.LastName, &res.UpdatedAt, &res.ProfilePhotoUrl, &res.Registered, &res.MobileNumber, &res.ClientUserId)
		// initialize new user
		newUser := model.NewUser(res.Email)
		newUser.SetId(res.ID)
		newUser.SetDevice(res.DeviceId, res.DeviceType)
		newUser.SetMobileNumber(res.MobileNumber)
		newUser.SetName(res.FirstName, res.LastName)
		newUser.SetProfilePhoto(res.ProfilePhotoUrl)
		newUser.SetUpdatedAt(res.UpdatedAt)
		// return the first user found
		return newUser
	}
	log.Output(1, "user not found")
	return nil
}

func (c *DBConnector) InsertNewUser(u *model.User) *model.User {
	log.Output(1, "Inserting New User")
	InsertUserQuery := `INSERT INTO users (id, email, first_name, last_name, updated_at, profile_photo_url, device_id, device_type, registered, mobile_number)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`
	res, err := c.DB.Exec(InsertUserQuery, (*u).GetId(), (*u).GetEmail(), (*u).GetFirstName(), (*u).GetLastName(), time.Now(), (*u).GetProfilePhoto(), (*u).GetDeviceId(), (*u).GetDeviceType(), (*u).IsRegistered(), (*u).GetMobileNumber())
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println(res)
	return u
}

// save user in db
func (c *DBConnector) SaveNewUser(u *model.User) *model.User {
	UpdateUserQuery := `UPDATE users SET
	first_name = $2, last_name = $3, updated_at = $4, profile_photo_url = $5
	WHERE email = $1
	RETURNING id`
	err := c.DB.QueryRow(UpdateUserQuery, (*u).GetEmail(), (*u).GetFirstName(), (*u).GetLastName(), time.Now(), (*u).GetProfilePhoto()).Scan(&(*u).ID)
	if err != nil {
		log.Println(err)
		return c.InsertNewUser(u)
	}
	return u
}

func (c *DBConnector) RegisterNewUser(u *model.User) *model.User {
	UpdateUserQuery := `UPDATE users set
	updated_at = $2, registered = $3, device_id = $4, device_type = $5
	WHERE email = $1`
	res, err := c.DB.Exec(UpdateUserQuery, (*u).GetEmail(), time.Now(), (*u).Registered, (*u).GetDeviceId(), (*u).GetDeviceType())
	if err != nil {
		log.Println(err)
		return nil
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	// unable to update so inserting new record
	if count == 0 {
		return c.InsertNewUser(u)
	} else {
		log.Output(1, "User Updated Sucessfully")
	}
	return u
}
func (c *DBConnector) UpdateClientUserId(ID *string, clientUserId *string) string {
	UpdateUserQuery := `UPDATE users SET client_user_id = $2 WHERE id = $1`
	rows, err := c.DB.Query(UpdateUserQuery, ID, clientUserId)
	if err != nil {
		log.Println(err)
		return "failed"
	}
	defer rows.Close()
	return "updated successfully"
}

func (c *DBConnector) getClientUserId(deviceId *string) string {
	var client_user_id string
	UpdateUserQuery := `SELECT client_user_id from users WHERE device_id = $1`
	err := c.DB.QueryRow(UpdateUserQuery, deviceId).Scan(&client_user_id)
	if err != nil {
		log.Println(err)
		return "failed to fetch client User ID"
	}
	fmt.Printf("The client Id is:")
	log.Output(1, client_user_id)
	return client_user_id
}
