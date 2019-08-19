package main

import (
	geoip2 "github.com/oschwald/geoip2-golang"
	"net"
)

// GetIPCoordinates converts a string to an IP object, passes the IP to the
// geoDB for a search and returns the IP's coordinates and an error
func GetIPCoordinates(ipAccess *IPAccess) {
	//TODO figure out if it's better to create a persistent db connection to the
	//GeoLite2 db instead of opening and closing the connection for each query
	db, err := geoip2.Open("GeoLite2-City_20190813/GeoLite2-City.mmdb")
	defer db.Close()

	//panicking if unable to successfully connect to the geoDB
	if err != nil {
		panic(err)
	}

	//TODO think of better way to validate strings here
	ip := net.ParseIP(ipAccess.IPAddress)
	record, err := db.City(ip)

	if err != nil {
		logger.Println(err)
	}

	ipAccess.Latitude = record.Location.Latitude
	ipAccess.Longitude = record.Location.Longitude
	ipAccess.Radius = record.Location.AccuracyRadius
}
