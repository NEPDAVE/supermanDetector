package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	logger *log.Logger //global logging object
)

/*
There is a bug in the challenge.

In the Expected Request there is a timestamp of:                                   1514764800
In the Expected Response for the proceeding ip access there is the same timestamp: 1514764800

Hong_Kong - 0,2
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764000, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e40", "ip_address": "119.28.48.231"}' http://localhost:5000/v1/ipaccess

New York - 1,1
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764001, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e41", "ip_address": "206.81.252.6"}' http://localhost:5000/v1/ipaccess

Moscow - 2,3
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764002, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42", "ip_address": "31.173.221.5"}' http://localhost:5000/v1/ipaccess

Hong_Kong - 3
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764003, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e43", "ip_address": "119.28.48.231"}' http://localhost:5000/v1/ipaccess

New York - 4
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764004, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e44", "ip_address": "206.81.252.7"}' http://localhost:5000/v1/ipaccess

New York - 5
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764005, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e45", "ip_address": "206.81.252.7"}' http://localhost:5000/v1/ipaccess

Sydney - 6,4
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764006, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e46", "ip_address": "203.2.218.214"}' http://localhost:5000/v1/ipaccess

*/

func IPAccessHandler(w http.ResponseWriter, r *http.Request) {
	// Double check it's a post request being made
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid http method - POST requests only\n")
		return
	}

	ipAccess := &IPAccess{}

	//Parse json request body and use it to set fields on user
	//Note that user is passed as a pointer variable so that it's fields can be modified
	err := json.NewDecoder(r.Body).Decode(&ipAccess)
	if err != nil {
		panic(err)
	}

	if err != nil {
		logger.Println(err)
	}

	//passing the IPAccess to NewReport to kick things off
	report := NewReport(ipAccess)

	//TODO make this not suck
	//checking that the report was populated
	reportCheck := &Report{}

	if report == reportCheck {
		fmt.Fprintf(w, "%s\n", "Server Error 500 - Please Try Again")
	}

	//Marshal or convert user object back to json and write to response
	reportJson, err := json.Marshal(*report)
	if err != nil {
		panic(err)
	}

	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//Write json response
	//json.NewEncoder(w).Encode(reportJson)
	//fmt.Fprintf(w, "%s\n", string(reportJson))
	fmt.Fprintf(w, "%s\n", reportJson)

}

func main() {
	//setting up logging
	logFile, err := os.Create("log.txt")
	defer logFile.Close()

	//panicking if unable to successfully create a log file
	if err != nil {
		panic(err)
	}

	logger = log.New(logFile, "supermanDetector ", log.LstdFlags|log.Lshortfile)

	//setting up sqlite db and migrating the schema
	MigrateDB()

	//starting the server
	logger.Println("Starting Server!")

	//TODO make sure to set timeouts on server
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/ipaccess", IPAccessHandler)
	http.ListenAndServe(":5000", mux)
}
