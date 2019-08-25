package main

import (
	geoip2 "github.com/oschwald/geoip2-golang"
	"net"
	"strings"
)

type Coordinates struct {
	Latitude  float64
	Longitude float64
	Radius    int
}

// GetIPCoordinates takes an IPAddress, passes the IP to the
// geoDB for a search and returns a Coordinates struct
func GetIPCoordinates(ipStr string) *Coordinates {
	//TODO figure out if it's better to create a persistent db connection to the
	//GeoLite2 db instead of opening and closing the connection for each query
	db, err := geoip2.Open("GeoLite2-City.mmdb")

	//panicking if unable to successfully connect to the geoDB
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()

	//TODO think of better way to validate string here
	ip := net.ParseIP(strings.TrimSpace(ipStr))
	record, err := db.City(ip)

	if err != nil {
		logger.Println(err)
		return &Coordinates{}
	}

	coordinates := &Coordinates{
		Latitude:  record.Location.Latitude,
		Longitude: record.Location.Longitude,
		Radius:    int(record.Location.AccuracyRadius), //final report expects int
	}

	return coordinates
}
