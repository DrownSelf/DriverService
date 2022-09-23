package models

type Rating struct {
	Id          int    `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	Rating      int    `json:"rating"`
}
