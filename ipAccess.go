package main

import (
	gorm "github.com/jinzhu/gorm"
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

//New IPAccess takes an *IPAccess populates Latitude, Longitude and Radius
//fields and then writes the struct to the sqlDB.
func NewIPAccess(ipAccess *IPAccess) *IPAccess {
	//Setting the IP Latitude, Longitude and Radius fields
	ipAccess.SetIPCoordinates()

	//Writing the fully populated IPAccess struct to the sqlDB
	CreateIPAccess(ipAccess)

	return ipAccess
}

//SetCoordinates calls GetIPCoordinates to query the geoDB and sets the IPAccess
//Latitude, Longitude and Radius fields
func (i *IPAccess) SetIPCoordinates() {
	coordinates := GetIPCoordinates(i.IPAddress)

	i.Latitude = coordinates.Latitude
	i.Longitude = coordinates.Longitude
	i.Radius = coordinates.Radius
}
