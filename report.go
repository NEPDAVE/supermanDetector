package main

import (
	"fmt"
	"github.com/umahmood/haversine"
)

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

func NewReport(ipAccess *IPAccess) *Report {
	ipAccess = NewIPAccess(ipAccess)
	report := &Report{}

	report.SetCurrentGeo(ipAccess.Latitude, ipAccess.Longitude, ipAccess.Radius)

	report.SetPrecedingIPAccess(ipAccess.UnixTimestamp, ipAccess.EventUUID,
		ipAccess.Username)

	report.SetSubsequentIPAccess(ipAccess.UnixTimestamp, ipAccess.EventUUID,
		ipAccess.Username)

	report.SetTravelToCurrentGeoSuspicious()

	report.SetTravelFromCurrentGeoSuspicious()

	return report
}

/*
because events can arrive out of order calculating speed if the booleans for
shit should be calculated each time on the fly
*/

func (r *Report) SetCurrentGeo(latitude float64, longitude float64, radius int) {
	r.CurrentGeo = CurrentGeo{
		Lat:    latitude,
		Lon:    longitude,
		Radius: radius,
	}
}

func (r *Report) SetPrecedingIPAccess(unixTimestamp int, eventUUID string,
	userName string) {
	precedingIPAccess := GetPrecedingIPAccess(unixTimestamp, eventUUID, userName)

	//checking for an ip of "" to prevent a false positive
	if precedingIPAccess.IPAddress == "" {
		return
	}

	r.PrecedingIPAccess = PrecedingIPAccess{
		IP:        precedingIPAccess.IPAddress,
		Lat:       precedingIPAccess.Latitude,
		Lon:       precedingIPAccess.Longitude,
		Radius:    precedingIPAccess.Radius,
		Timestamp: precedingIPAccess.UnixTimestamp,
	}

	//calculating speed for PrecedingIPAccess
	currentCoord := haversine.Coord{Lat: r.CurrentGeo.Lat, Lon: r.CurrentGeo.Lon}
	precedingCoord := haversine.Coord{Lat: r.PrecedingIPAccess.Lat,
		Lon: r.PrecedingIPAccess.Lon}

	_, km := haversine.Distance(currentCoord, precedingCoord)

	time := unixTimestamp - precedingIPAccess.UnixTimestamp
	r.PrecedingIPAccess.Speed = CalculateSpeed(time, km)
}

func (r *Report) SetSubsequentIPAccess(unixTimestamp int, eventUUID string,
	userName string) {
	subsequentIPAccess := GetSubsequentIPAccess(unixTimestamp, eventUUID, userName)

	//checking for an ip of "" to prevent a false positive
	if subsequentIPAccess.IPAddress == "" {
		fmt.Println("nothing here")
		return
	}

	r.SubsequentIPAccess = SubsequentIPAccess{
		IP:        subsequentIPAccess.IPAddress,
		Lat:       subsequentIPAccess.Latitude,
		Lon:       subsequentIPAccess.Longitude,
		Radius:    subsequentIPAccess.Radius,
		Timestamp: subsequentIPAccess.UnixTimestamp,
	}

	//calculating speed for PrecedingIPAccess
	currentCoord := haversine.Coord{Lat: r.CurrentGeo.Lat, Lon: r.CurrentGeo.Lon}
	precedingCoord := haversine.Coord{Lat: r.SubsequentIPAccess.Lat,
		Lon: r.SubsequentIPAccess.Lon}

	_, km := haversine.Distance(currentCoord, precedingCoord)

	time := subsequentIPAccess.UnixTimestamp - unixTimestamp
	speed := CalculateSpeed(time, km)
	r.SubsequentIPAccess.Speed = speed

	fmt.Println("SubsequentIPAccess shit")
	fmt.Println(km)
	fmt.Println(speed)
}

func (r *Report) SetTravelToCurrentGeoSuspicious() {
	if r.PrecedingIPAccess.Speed >= 500 {
		r.TravelToCurrentGeoSuspicious = true
	}
	fmt.Println(r.TravelToCurrentGeoSuspicious)
}

func (r *Report) SetTravelFromCurrentGeoSuspicious() {
	if r.SubsequentIPAccess.Speed >= 500 {
		r.TravelFromCurrentGeoSuspicious = true
	}

	fmt.Println(r.TravelFromCurrentGeoSuspicious)
}

//CalculateSpeed takes an amount of time in seconds and a distance in Kilometers
//and returns the Kilometers per hour
func CalculateSpeed(time int, km float64) int {
	floatTime := float64(time)
	//calculates KM per second
	kmPerSecond := km / floatTime
	//converts the KM per second to KM per hour
	kmPerHour := kmPerSecond * 60 * 60
	return int(kmPerHour)
}
