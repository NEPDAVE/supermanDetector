package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//Table test inputs and expected outputs
var (
	//"Bob" New York IP Access at 1514764001
	input1 = `{
	"username": "bob",
	"unix_timestamp": 1514764001,
	"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e41",
	"ip_address": "206.81.252.6"
}`
	output1 = `{
	"currentGeo": {
		"lat": 39.211,
		"lon": -76.8362,
		"radius": 5
	},
	"travelToCurrentGeoSuspicious": false,
	"travelFromCurrentGeoSuspicious": false,
	"precedingIpAccess": {
		"ip": "",
		"speed": 0,
		"lat": 0,
		"lon": 0,
		"radius": 0,
		"timestamp": 0
	},
	"subsequentIpAccess": {
		"ip": "",
		"speed": 0,
		"lat": 0,
		"lon": 0,
		"radius": 0,
		"timestamp": 0
	}
}`

	//"Bob" Hong Kong IP Access 1514764000
	input2 = `{
	"username": "bob",
	"unix_timestamp": 1514764000,
	"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e40",
	"ip_address": "119.28.48.231"
}`
	output2 = `{
	"currentGeo": {
		"lat": 22.25,
		"lon": 114.1667,
		"radius": 50
	},
	"travelToCurrentGeoSuspicious": false,
	"travelFromCurrentGeoSuspicious": true,
	"precedingIpAccess": {
		"ip": "",
		"speed": 0,
		"lat": 0,
		"lon": 0,
		"radius": 0,
		"timestamp": 0
	},
	"subsequentIpAccess": {
		"ip": "206.81.252.6",
		"speed": 47108576,
		"lat": 39.211,
		"lon": -76.8362,
		"radius": 5,
		"timestamp": 1514764001
	}
}`

	//"Bob" Moscow IP Access at 1514764002
	input3 = `{
	"username": "bob",
	"unix_timestamp": 1514764002,
	"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42",
	"ip_address": "31.173.221.5"
}`
	output3 = `{
	"currentGeo": {
		"lat": 42.9753,
		"lon": 47.5022,
		"radius": 1000
	},
	"travelToCurrentGeoSuspicious": true,
	"travelFromCurrentGeoSuspicious": false,
	"precedingIpAccess": {
		"ip": "206.81.252.6",
		"speed": 33472130,
		"lat": 39.211,
		"lon": -76.8362,
		"radius": 5,
		"timestamp": 1514764001
	},
	"subsequentIpAccess": {
		"ip": "",
		"speed": 0,
		"lat": 0,
		"lon": 0,
		"radius": 0,
		"timestamp": 0
	}
}`

	//"Bob" Sydney IP Access at 1514764006
	input4 = `{
	"username": "bob",
	"unix_timestamp": 1514764006,
	"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e46",
	"ip_address": "203.2.218.214"
}`

	output4 = `{
	"currentGeo": {
		"lat": -33.8919,
		"lon": 151.1554,
		"radius": 500
	},
	"travelToCurrentGeoSuspicious": true,
	"travelFromCurrentGeoSuspicious": false,
	"precedingIpAccess": {
		"ip": "31.173.221.5",
		"speed": 12165789,
		"lat": 42.9753,
		"lon": 47.5022,
		"radius": 1000,
		"timestamp": 1514764002
	},
	"subsequentIpAccess": {
		"ip": "",
		"speed": 0,
		"lat": 0,
		"lon": 0,
		"radius": 0,
		"timestamp": 0
	}
}`

	//"Bob" New York IP Access at 1514764005
	input5 = `{
	"username": "bob",
	"unix_timestamp": 1514764005,
	"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e45",
	"ip_address": "206.81.252.7"
}`

	output5 = `{
	"currentGeo": {
		"lat": 39.211,
		"lon": -76.8362,
		"radius": 5
	},
	"travelToCurrentGeoSuspicious": true,
	"travelFromCurrentGeoSuspicious": true,
	"precedingIpAccess": {
		"ip": "31.173.221.5",
		"speed": 11157376,
		"lat": 42.9753,
		"lon": 47.5022,
		"radius": 1000,
		"timestamp": 1514764002
	},
	"subsequentIpAccess": {
		"ip": "203.2.218.214",
		"speed": 56655567,
		"lat": -33.8919,
		"lon": 151.1554,
		"radius": 500,
		"timestamp": 1514764006
	}
}`

	//"Bob" Hong Kong IP Access at 1514764003
	input6 = `{
	"username": "bob",
	"unix_timestamp": 1514764003,
	"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e43",
	"ip_address": "119.28.48.231"
  }`

	output6 = `{
	"currentGeo": {
		"lat": 22.25,
		"lon": 114.1667,
		"radius": 50
	},
	"travelToCurrentGeoSuspicious": true,
	"travelFromCurrentGeoSuspicious": true,
	"precedingIpAccess": {
		"ip": "31.173.221.5",
		"speed": 23313771,
		"lat": 42.9753,
		"lon": 47.5022,
		"radius": 1000,
		"timestamp": 1514764002
	},
	"subsequentIpAccess": {
		"ip": "206.81.252.7",
		"speed": 23554288,
		"lat": 39.211,
		"lon": -76.8362,
		"radius": 5,
		"timestamp": 1514764005
	}
}`

	//"Bob" Sydney IP Access at 1514764007
	input7 = `{
	"username": "bob",
	"unix_timestamp": 1514764007,
	"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e47",
	"ip_address": "203.2.218.214"
 }`

	output7 = `{
	"currentGeo": {
		"lat": -33.8919,
		"lon": 151.1554,
		"radius": 500
	},
	"travelToCurrentGeoSuspicious": false,
	"travelFromCurrentGeoSuspicious": false,
	"precedingIpAccess": {
		"ip": "203.2.218.214",
		"speed": 0,
		"lat": -33.8919,
		"lon": 151.1554,
		"radius": 500,
		"timestamp": 1514764006
	},
	"subsequentIpAccess": {
		"ip": "",
		"speed": 0,
		"lat": 0,
		"lon": 0,
		"radius": 0,
		"timestamp": 0
	}
}`
)

//TestAPI tests a series of POST requests in a specific order. The prearranged
//order allows us to test for several different scenarios and get a comprehesive
//total test of the entire API.

//TestAPI does drop the IPAccess table and remigrate it each time the test is run.
//This allows us to control the specific prearranged order of the POST requests.
func TestAPI(t *testing.T) {

	//Droping IPAccess table the remigrating it
	MigrateDB()

	tableTests := []struct {
		in  string
		out string
	}{
		{input1, output1},
		{input2, output2},
		{input3, output3},
		{input4, output4},
		{input5, output5},
		{input6, output6},
		{input7, output7},
	}

	for _, tt := range tableTests {
		t.Run(tt.in, func(t *testing.T) {
			jsonStr := []byte(tt.in)

			req, err := http.NewRequest("POST", "localhost:5000/v1/ipaccess",
				bytes.NewBuffer(jsonStr))

			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(IPAccessHandler)
			handler.ServeHTTP(w, req)

			if status := w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			expected := tt.out

			//the multiline strings in the test have random whitespace, trimming it here
			if strings.Join(strings.Fields(w.Body.String()), "") != strings.Join(strings.Fields(expected), "") {
				t.Errorf("handler returned unexpected body: \n\n got \n\n %v \n want \n %v",
					w.Body.String(), expected)
			}
		})
	}
}
