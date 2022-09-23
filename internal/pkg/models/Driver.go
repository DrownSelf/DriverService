package models

type Driver struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	TaxiType    string `json:"taxiType"`
	Status      bool   `json:"status"`
}
