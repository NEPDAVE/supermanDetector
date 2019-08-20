package main

import (
	"encoding/json"
	gorm "github.com/jinzhu/gorm"
	"net/http"
)

//IPAccess represents the JSON data structure posted to the supermanDetector API.
//The Latitude, Longitude and Radius fields are populated after receiving the POST
//request by quering the geoDB. Once the additional fields are added, the
//the IPAccess is written to the sqlDB.
type IPAccess struct {
	gorm.Model
	Username      string  `json:"username"`
	UnixTimestamp int     `json:"unix_timestamp"`
	EventUUID     string  `json:"event_uuid"`
	IPAddress     string  `json:"ip_address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Radius        int     `json:"radius"`
}

//New IPAccess takes an *http.Request, decodes it into an IPAccess struct,
//populates Latitude, Longitude and Radius fields and then writes the struct
//to the sqlDB.
func (i IPAccess) NewIPAccess(r *http.Request) *IPAccess {
	ipAccess := &IPAccess{}

	//Parse json request body and use it to set fields on user
	//Note that user is passed as a pointer variable so that it's fields can be modified
	err := json.NewDecoder(r.Body).Decode(&ipAccess)
	if err != nil {
		panic(err)
	}

	if err != nil {
		logger.Println(err)
	}

	//Setting the IP Latitude, Longitude and Radius fields
	ipAccess.SetIPCoordinates()
	//Writing the fully populated IPAccess struct to the sqlDB
	ipAccess.WriteIPAccessToDB()

	return ipAccess
}

//SetCoordinates calls GetIPCoordinates to query the geoDB and set the IPAccess
//Latitude, Longitude and Radius fields
func (i *IPAccess) SetIPCoordinates() {
	GetIPCoordinates(i)
}

//WriteIPAccessToDB writes an IPAccess struct to the sqlDB
func (i *IPAccess) WriteIPAccessToDB() {
	CreateIPAccess(i)
}
