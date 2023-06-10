package main

import (
	"MYPLANTBACKEND/common"
	"MYPLANTBACKEND/model"
	"MYPLANTBACKEND/service"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"log"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"
)

type UserRequestBody struct {
	Email     string
	PhotoUrl  string
	FirstName string
	LastName  string
	Token     string
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
	userToBeSaved.SetNotificationToken(requestBody.Token)

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

	c.IndentedJSON(http.StatusOK, gin.H{
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

type demo struct {
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle

	w   io.Writer
	ctx context.Context
	// cleanUp is a list of filenames that need cleaning up at the end of the demo.
	// failed indicates that one or more of the demo steps failed.
	failed bool
}

// readFile reads the named file in Google Cloud Storage.
func (d *demo) readFile(fileName string) {
	io.WriteString(d.w, "\nAbbreviated file content (first line and last 1K):\n")

	rc, err := d.bucket.Object(fileName).NewReader(d.ctx)
	if err != nil {
		log.Printf("readFile: unable to open file from bucket %q, file %q: %v", d.bucketName, fileName, err)
		return
	}
	defer rc.Close()
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Printf("readFile: unable to read data from bucket %q, file %q: %v", d.bucketName, fileName, err)
		return
	}

	fmt.Fprintf(d.w, "%s\n", bytes.SplitN(slurp, []byte("\n"), 2)[0])
	if len(slurp) > 1024 {
		fmt.Fprintf(d.w, "...%s\n", slurp[len(slurp)-1024:])
	} else {
		fmt.Fprintf(d.w, "%s\n", slurp)
	}
}

// readFile reads the named file in Google Cloud Storage.
func readFile(c *gin.Context, fileName string) {

	ctx := appengine.NewContext(c.Request)

	bucket := "myplantapk"

	log.Println(bucket)

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Println(ctx, "failed to create client: %v", err)
		return
	}
	defer client.Close()

	c.Writer.Header().Set("Content-Type", "application/vnd.android.package-archive; charset=utf-8")
	// log.Println("Demo GCS Application running from Version: %v\n", appengine.VersionID(ctx))
	// log.Println("Using bucket name: %v\n\n", bucket)

	buf := &bytes.Buffer{}
	d := &demo{
		w:          buf,
		ctx:        ctx,
		client:     client,
		bucket:     client.Bucket(bucket),
		bucketName: bucket,
	}

	d.readFile(fileName)
	log.Println("reada file")

	if d.failed {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		buf.WriteTo(c.Writer)
		log.Println("Demo failed.")
	} else {
		c.Writer.WriteHeader(http.StatusOK)
		buf.WriteTo(c.Writer)
		fmt.Println("Demo succeeded.")
	}

}

var logsFilePath = os.Getenv("LOG_FILE_PATH")

// Replace with the path where you want to store the logs file

func uploadLogsHandler(w http.ResponseWriter, r *http.Request) {
	var LogTime string = time.Now().Format("3:4:05 PM")

	// Read the logs from the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	fmt.Println("Logs File Path : ", logsFilePath)

	file, err := os.OpenFile(logsFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// Add a newline character before appending the logs

	_, err = file.WriteString("\n" + LogTime + "\n")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Append the logs to the file
	_, err = file.Write(body)
	if err != nil {
		http.Error(w, "Failed to write logs to file", http.StatusInternalServerError)
		return
	}

	// Return success response
	fmt.Fprintf(w, "Logs uploaded successfully")
}

func main() {
	// to:= "mukhar.jain2009@gmail.com"
	// err := service.SendEmail(to)
	// if err != nil {
	// 	fmt.Println("error in sending email")
	// }else{
	// 	fmt.Println("sent")
	// }

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

	router.GET("download/apk/:name", func(context *gin.Context) {
		fileName := context.Param("name")
		readFile(context, fileName)
	})

	router.GET("notification/send/:id", func(context *gin.Context) {
		deviceId := context.Param("id")
		var notificationTime string = time.Now().Format("3:4:05 PM")
		service.PushNotification(deviceId, "Water Turned On", "At "+notificationTime)
		context.IndentedJSON(http.StatusOK, "pushed")
	})

	router.POST("upload_logs", func(context *gin.Context) {
		uploadLogsHandler(context.Writer, context.Request)
		context.IndentedJSON(http.StatusOK, "logs saved Succcessfully")

	})

	router.Run(":8080")
}
