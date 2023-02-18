package model

import "time"

type User struct {
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
}

type UserImpl interface {
	GetId() int
	SetId(id int)
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
}

func NewUser(Email string) *User {
	u := &User{
		Email:           "",
		DeviceId:        "",
		DeviceType:      "",
		FirstName:       "",
		LastName:        "",
		ProfilePhotoUrl: "",
		MobileNumber:    "",
	}
	u.SetEmail(Email)
	return u
}

func (u *User) SetId(id int) {
	u.ID = id
}

func (u *User) GetId() int {
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

func (u *User) SetUpdatedAt(uat time.Time) {
	u.UpdatedAt = uat
}
