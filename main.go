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
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764800, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42", "ip_address": "206.81.252.6"}' http://localhost:5000/v1/ipaccess
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
	fmt.Fprintf(w, "%s\n", string(reportJson))

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
