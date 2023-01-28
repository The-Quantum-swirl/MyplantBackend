package service

import (
	"database/sql"
	"fmt"
	"os"

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
	// fmt.Print(connString)

	var err error

	c.DB, err = sql.Open("postgres", connString)
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

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// func (c *DBConnector) GetUser() {
// 	var users []User
// 	db.Find(&users)

// 	return users
// 	// c.JSON(http.StatusOK, gin.H{"data": users})
// }

// func main() {

// 	// Use Gin to handle HTTP requests
// 	r := gin.Default()
// 	r.GET("/data", func(cgin.Context) {
// 		var users []User
// 		db.Find(&users)
// 		c.JSON(http.StatusOK, gin.H{"data": users})
// 	})

// 	r.GET("/", func(c gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "Hello, World!",
// 		})
// 	})

// 	r.POST("/saveUser", func(cgin.Context) {
// 		us := User{ID: 1, Username: "mukh"}
// 		log(us.ID)
// 		logs(us.Username)
// 		if !db.Migrator().HasTable("users") {
// 			db.Migrator().CreateTable(&User{})
// 		}
// 		res := db.Create(&us)
// 		logs(res.Error.Error())
// 	})
// 	r.Run(":8080")
// }
