package main

import (
	geoip2 "github.com/oschwald/geoip2-golang"
	"net"
)

// GetIPCoordinates takes a pointer to an IPAccess, passes the IP to the
// geoDB for a search and then sets the Latitude, Longitude, and Radius on
//the IPAccess pointer
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
	//final report expects radius as int NOT uint16, doing conversion here
	ipAccess.Radius = int(record.Location.AccuracyRadius)
}
