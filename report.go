package main

import (
	"fmt"
	"net/http"
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

func (rp Report) NewReport(r *http.Request) *Report {
	ipAccess := IPAccess{}.NewIPAccess(r)

	report := &Report{
		CurrentGeo: CurrentGeo{
			Lat:    ipAccess.Latitude,
			Lon:    ipAccess.Longitude,
			Radius: ipAccess.Radius,
		},
	}

	report.SetPrecedingIPAccess(ipAccess.UnixTimestamp, ipAccess.Username)
	report.SetSubsequentIPAccess(ipAccess.UnixTimestamp, ipAccess.Username)
	return report
}

/*
because events can arrive out of order calculating speed if the booleans for
shit should be calculated each time on the fly
*/

//TODO have this use pointers like the other one
func (r *Report) SetPrecedingIPAccess(unixTimestamp int, userName string) {
	//this is where speed should be set/calculated
	precedingIPAccess := GetPrecedingIPAccess(unixTimestamp, userName)

	fmt.Println("PRECEDING S Q L MO FOOOKAA")
	fmt.Println(precedingIPAccess)
	fmt.Println("PRECEDING S Q L MO FOOOKAA")
}

//TODO have this use pointers like the other one
func (r *Report) SetSubsequentIPAccess(unixTimestamp int, userName string) {
	//this is where speed should be set/calculated
	subsequentIPAccess := GetSubsequentIPAccess(unixTimestamp, userName)

	fmt.Println("Subsequent S Q L MO FOOOKAA")
	fmt.Println(subsequentIPAccess)
	fmt.Println("Subsequent S Q L MO FOOOKAA")

}

func (r *Report) SetTravelToCurrentGeoSuspicious() {

}

func (r *Report) SetTravelFromCurrentGeoSuspicious() {

}
