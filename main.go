package main

import (
	"fmt"
	gorm "github.com/jinzhu/gorm"
	"log"
	"os"
)

var (
	logger *log.Logger //logging object
)

func main() {
	//creating logging file
	logFile, err := os.Create("log.txt")
	defer logFile.Close()

	//panicking if unable to successfully create a logging file
	if err != nil {
		panic(err)
	}

	//creating the logging object
	logger = log.New(logFile, "supermanDetector ", log.LstdFlags|log.Lshortfile)
	logger.Println("Program Start!")

	//NEED THIS
	//https://tutorialedge.net/golang/golang-orm-tutorial/
	//myRouter.HandleFunc("/user/{name}/{email}", newUser).Methods("POST")

	MigrateDB()

	ipAccess := IPAccess{
		Username:      "bob",
		UnixTimestamp: 1514764800,
		EventUUID:     "85ad929a-db03-4bf4-9541-8f728fa12e42",
		IPAddress:     "206.81.252.6",
	}

	ipAccess.SetIPCoordinates()
	ipAccess.WriteIPAccessToDB()

	fmt.Println(ipAccess)

}

//IPAccess represent the JSON data structure posted to the supermanDetector API.
//The Latitude and Longitude fields are populated after receiving the POST
//request by quering the geoDB and adding them
type IPAccess struct {
	gorm.Model
	Username             string
	UnixTimestamp        int
	EventUUID            string
	IPAddress            string
	Latitude             float64
	Longitude            float64
	Radius               uint16
  Speed                int
	TravelToCurrentGeoSuspicious bool
  TravelFromCurrentGeoSuspicious bool

}

//SetCoordinates calls GetIPCoordinates to query the geoDB and set the IPAccess
//Latitude and Longitude struct fields
func (i *IPAccess) SetIPCoordinates() {
	GetIPCoordinates(i)
}

func (i *IPAccess) WriteIPAccessToDB() {
	CreateIPAccess(i)
}

func (i *IPAccess) SetPrecedingIPAccess() {

}

func (i *IPAccess) SetSubsequentIPAccess() {

}

func (i *IPAccess) SetTravelToCurrentGeoSuspicious() {

}

func (i *IPAccess) SetTravelFromCurrentGeoSuspicious() {

}


/*
type Report struct {
	CurrentGeo                     CurrentGeo         `json:"currentGeo"`
	TravelToCurrentGeoSuspicious   bool               `json:"travelToCurrentGeoSuspicious"`
	TravelFromCurrentGeoSuspicious bool               `json:"travelFromCurrentGeoSuspicious"`
	PrecedingIPAccess              PrecedingIPAccess  `json:"precedingIpAccess"`
	SubsequentIPAccess             SubsequentIPAccess `json:"subsequentIpAccess"`
}
type CurrentGeo struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Radius int     `json:"radius"`
}
type PrecedingIPAccess struct {
	IP        string  `json:"ip"`
	Speed     int     `json:"speed"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Radius    int     `json:"radius"`
	Timestamp int     `json:"timestamp"`
}
type SubsequentIPAccess struct {
	IP        string  `json:"ip"`
	Speed     int     `json:"speed"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Radius    int     `json:"radius"`
	Timestamp int     `json:"timestamp"`
}



type AutoGenerated struct {
	CurrentGeo struct {
		Lat    float64 `json:"lat"`
		Lon    float64 `json:"lon"`
		Radius int     `json:"radius"`
	} `json:"currentGeo"`
	TravelToCurrentGeoSuspicious   bool `json:"travelToCurrentGeoSuspicious"`
	TravelFromCurrentGeoSuspicious bool `json:"travelFromCurrentGeoSuspicious"`
	PrecedingIPAccess              struct {
		IP        string  `json:"ip"`
		Speed     int     `json:"speed"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		Radius    int     `json:"radius"`
		Timestamp int     `json:"timestamp"`
	} `json:"precedingIpAccess"`
	SubsequentIPAccess struct {
		IP        string  `json:"ip"`
		Speed     int     `json:"speed"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		Radius    int     `json:"radius"`
		Timestamp int     `json:"timestamp"`
	} `json:"subsequentIpAccess"`
}

{
	"currentGeo": {
		"lat": 39.1702,
		"lon": -76.8538,
		"radius": 20
	},
	"travelToCurrentGeoSuspicious": true,
	"travelFromCurrentGeoSuspicious": false,
	"precedingIpAccess": {
		"ip": "24.242.71.20",
		"speed": 55,
		"lat": 30.3764,
		"lon": -97.7078,
		"radius": 5,
		"timestamp": 1514764800
	},
	"subsequentIpAccess": {
		"ip": "91.207.175.104",
		"speed": 27600,
		"lat": 34.0494,
		"lon": -118.2641,
		"radius": 200,
		"timestamp": 1514851200
	}
}

*/
