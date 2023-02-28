package model

import (
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	DeviceId        string    `json:"deviceId"`
	DeviceType      string    `json:"deviceType"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	UpdatedAt       time.Time `json:"updatedAt"`
	ProfilePhotoUrl string    `json:"profilePhotoUrl"`
	Registered      bool      `json:"registered"`
	MobileNumber    string    `json:"mobileNumber"`
	ClientUserId    string    `json:"clientUserId"`
}

type UserImpl interface {
	GetId() uuid.UUID
	SetId(id uuid.UUID)
	GetEmail() string
	SetEmail(Email string)
	GetName() string
	GetFirstName() string
	GetLastName() string
	SetName(firstName string, lastName string)
	IsRegistered() bool
	RegisterIt()
	SetDevice(deviceId string, deviceType string)
	GetDeviceId() string
	GetDeviceType()
	SetProfilePhoto(url string)
	GetProfilePhoto() string
	SetMobileNumber(no string)
	GetMobileNumber() string
	SetClientUserId(clientUserId string)
	GetClientUserId() string
}

func NewUser(Email string) *User {
	u := &User{
		ID:              uuid.New(),
		Email:           strings.ToLower(Email),
		DeviceId:        "default",
		DeviceType:      "default",
		FirstName:       "",
		LastName:        "",
		UpdatedAt:       time.Now(),
		ProfilePhotoUrl: "",
		MobileNumber:    "",
		Registered:      false,
		ClientUserId:    "null",
	}
	log.Output(1, u.GetId().String())
	return u
}

func (u *User) SetId(id uuid.UUID) {
	u.ID = id
}

func (u *User) GetId() uuid.UUID {
	return u.ID
}

func (u *User) SetEmail(Email string) {
	u.Email = Email
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) SetName(firstName string, lastName string) {
	u.FirstName = firstName
	u.LastName = lastName
}

func (u *User) GetName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) GetFirstName() string {
	return u.FirstName
}

func (u *User) GetLastName() string {
	return u.LastName
}

func (u *User) RegisterIt() {
	u.Registered = true
}

func (u *User) IsRegistered() bool {
	return u.Registered
}

func (u *User) SetDevice(deviceId string, deviceType string) {
	u.DeviceId = deviceId
	u.DeviceType = deviceType
}

func (u *User) GetDeviceId() string {
	return u.DeviceId
}

func (u *User) GetDeviceType() string {
	return u.DeviceType
}

func (u *User) SetProfilePhoto(url string) {
	u.ProfilePhotoUrl = url
}

func (u *User) GetProfilePhoto() string {
	return u.ProfilePhotoUrl
}

func (u *User) SetMobileNumber(no string) {
	u.MobileNumber = no
}

func (u *User) GetMobileNumber() string {
	return u.MobileNumber
}

func (u *User) SetClientUserId(clientUserId string) {
	u.MobileNumber = clientUserId
}

func (u *User) GetClientUserId() string {
	return u.ClientUserId
}

func (u *User) SetUpdatedAt(uat time.Time) {
	u.UpdatedAt = uat
}
