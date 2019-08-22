package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
  //"Bob" New York IP Access at 1514764001
	inputA = `{"username": "bob", "unix_timestamp": 1514764001,
    "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e41",
    "ip_address": "206.81.252.6"}`
  outputA = `{"currentGeo":{"lat":39.211,"lon":-76.8362,"radius":5},
  "travelToCurrentGeoSuspicious":false,"travelFromCurrentGeoSuspicious":false,
  "precedingIpAccess":{"ip":"","speed":0,"lat":0,"lon":0,"radius":0,"timestamp":0},
  "subsequentIpAccess":{"ip":"","speed":0,"lat":0,"lon":0,"radius":0,"timestamp":0}}`

  //"Bob" Hong Kong IP Access 1514764000
	inputB = `{"username": "bob", "unix_timestamp": 1514764000,
    "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e40",
    "ip_address": "119.28.48.231"}`
	outputB = `{"currentGeo":{"lat":22.25,"lon":114.1667,"radius":50},
  "travelToCurrentGeoSuspicious":false,"travelFromCurrentGeoSuspicious":true,
  "precedingIpAccess":{"ip":"","speed":0,"lat":0,"lon":0,"radius":0,"timestamp":0},
  "subsequentIpAccess":{"ip":"206.81.252.6","speed":47108576,"lat":39.211,
    "lon":-76.8362,"radius":5,"timestamp":1514764001}}`

  //"Bob" Moscow IP Access at 1514764002
  inputC = `{"username": "bob", "unix_timestamp": 1514764002,
    "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42",
    "ip_address": "31.173.221.5"}`
  outputC = `{"currentGeo":{"lat":42.9753,"lon":47.5022,"radius":1000},
  "travelToCurrentGeoSuspicious":true,"travelFromCurrentGeoSuspicious":false,
  "precedingIpAccess":{"ip":"206.81.252.6","speed":33472130,"lat":39.211,
    "lon":-76.8362,"radius":5,"timestamp":1514764001},
    "subsequentIpAccess":{"ip":"","speed":0,"lat":0,"lon":0,"radius":0,"timestamp":0}}`

	//inputD =

)

// tableTests := []struct {
// 	in  string
// 	out string
// }{
// 	{inputA, outputA},
//   {inputB, outputB},
//   {inputC, outputC},
// }
// {"%-a", "[%-a]"},
// {"%+a", "[%+a]"},

func TestAPI(t *testing.T) {
	tableTests := []struct {
		in  string
		out string
	}{
		{inputA, outputA},
		{inputB, outputB},
		//{inputC, outputC},
	}
	//var flagprinter flagPrinter
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

      if strings.Join(strings.Fields(w.Body.String()), "") != strings.Join(strings.Fields(expected), "") {
        t.Errorf("handler returned unexpected body: \n\n got \n\n %v \n want \n %v",
          w.Body.String(), expected)
      }
		})
	}
}
