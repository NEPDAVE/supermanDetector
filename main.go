package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	logger *log.Logger //global logging object
)

//IPAccessHandler takes a pointer to an HTTP request about a new IP Access and
//returns JSON data about the IP Access and if the IP Access is suspicious.
func IPAccessHandler(w http.ResponseWriter, r *http.Request) {
	//capturing panic from IPAccessHandler
	defer func() {
		if r := recover(); r != nil {
			//TODO figure out why calls to logger here cause more panics
			fmt.Println("Recovered in f", r)
			fmt.Fprintf(w, "%s\n", "Server Error 500 - Please try again")
			return
		}
	}()

	// Double check it's a post request being made
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		logger.Printf("invalid user HTTP method type: %v", r.Method)
		fmt.Fprintf(w, "%s\n", "Invalid HTTP method - Please try again")
		return
	}

	//create new pointer to IPAccess object
	ipAccess := &IPAccess{}

	//Parse json request body and use it to set fields on ipAccess
	err := json.NewDecoder(r.Body).Decode(&ipAccess)

	if err != nil {
		logger.Printf("invalid POST request: %v", err)
		fmt.Fprintf(w, "%s\n", "Invalid POST request - Please try again")
		return
	}

	//passing the IPAccess to NewReport to kick things off
	report := NewReport(ipAccess)

	//checking that the report was populated
	reportCheck := &Report{}

	if report == reportCheck {
		logger.Println("empty report")
		fmt.Fprintf(w, "%s\n", "Server Error 500 - Please try again")
	}

	//Marshal or convert user object back to json and write to response
	reportJson, err := json.Marshal(*report)
	if err != nil {
		logger.Println(err)
	}

	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//Write json response
	fmt.Fprintf(w, "%s\n", reportJson)
}

func main() {
	//setting up logging
	logFile, err := os.Create("log.txt")

	//panicking if unable to successfully create a log file
	if err != nil {
		panic(err)
	}

	defer logFile.Close()

	logger = log.New(logFile, "supermanDetector ", log.LstdFlags|log.Lshortfile)

	//setting up sqlite db and migrating the schema
	MigrateDB()

	//starting the server
	logger.Println("Starting Server!")
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/ipaccess", IPAccessHandler)

	//creating custom server with timeouts and custom handler
	s := &http.Server{
		Addr:         ":5000",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.ListenAndServe()
}
