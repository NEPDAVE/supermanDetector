package main

import (
//"fmt"
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
	report.SetPrecedingIPAccess(ipAccess.UnixTimestamp, ipAccess.Username)
	report.SetSubsequentIPAccess(ipAccess.UnixTimestamp, ipAccess.Username)

	return report
}

/*
because events can arrive out of order calculating speed if the booleans for
shit should be calculated each time on the fly
*/

func (rp *Report) SetCurrentGeo(latitude float64, longitude float64, radius int) {
	rp.CurrentGeo = CurrentGeo{
		Lat:    latitude,
		Lon:    longitude,
		Radius: radius,
	}
}

func (rp *Report) SetPrecedingIPAccess(unixTimestamp int, userName string) {
	precedingIPAccess := GetPrecedingIPAccess(unixTimestamp, userName)

	rp.PrecedingIPAccess = PrecedingIPAccess{
		IP:        precedingIPAccess.IPAddress,
		Lat:       precedingIPAccess.Latitude,
		Lon:       precedingIPAccess.Longitude,
		Radius:    precedingIPAccess.Radius,
		Timestamp: precedingIPAccess.UnixTimestamp,
	}

	//TODO need to calculate speed for PrecedingIPAccess
}

func (rp *Report) SetSubsequentIPAccess(unixTimestamp int, userName string) {
	subsequentIPAccess := GetSubsequentIPAccess(unixTimestamp, userName)

	rp.SubsequentIPAccess = SubsequentIPAccess{
		IP:        subsequentIPAccess.IPAddress,
		Lat:       subsequentIPAccess.Latitude,
		Lon:       subsequentIPAccess.Longitude,
		Radius:    subsequentIPAccess.Radius,
		Timestamp: subsequentIPAccess.UnixTimestamp,
	}

	//TODO need to calculate speed for PrecedingIPAccess
}

func (rp *Report) SetTravelToCurrentGeoSuspicious() {

}

func (rp *Report) SetTravelFromCurrentGeoSuspicious() {

}
