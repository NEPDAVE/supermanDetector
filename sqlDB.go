package main

import (
	//"fmt"
	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//MigrateDB maps the Login struct fields to a newly created IPAccess SQL table
//where we can store IPAccesss
func MigrateDB() {
	db, err := gorm.Open("sqlite3", "supermanDetector.db")

	//panicking if unable to successfully connect to the sqliteDB
	if err != nil {
		logger.Println("unable to migrate database")
		logger.Fatalln(err)
	}
	defer db.Close()

	// Migrate the schemaer
	db.AutoMigrate(&IPAccess{})
}

//CreateIPAccess creates a new IPAccess entry in the IPAccess db table
func CreateIPAccess(ipAccess *IPAccess) {
	db, err := gorm.Open("sqlite3", "supermanDetector.db")

	if err != nil {
		logger.Println("unable to create IPAccess in database")
		logger.Fatalln(err)
	}
	defer db.Close()

	db.Create(&ipAccess)
}

//DropIPAccessTable drops the IP Access table to allow for a clean table for testing
func DropIPAccessTable() {
	db, err := gorm.Open("sqlite3", "supermanDetector.db")

	if err != nil {
		logger.Printf("unable to drop IPAccess table in database: %v", err)
	}
	defer db.Close()

	// Drop model IPAccess table
	db.DropTable(&IPAccess{})
}

func GetPrecedingIPAccess(unixTimestamp int, eventUUID string,
	userName string) *IPAccess {
	db, err := gorm.Open("sqlite3", "supermanDetector.db")

	if err != nil {
		logger.Printf("unable to get preceding IPAccess from database: %v", err)
	}
	defer db.Close()

	ipAccess := &IPAccess{}

	db.Raw(`SELECT * FROM ip_accesses WHERE ip_accesses.unix_timestamp < ?
		AND ip_accesses.event_uuid != ? AND ip_accesses.username = ?
		ORDER BY unix_timestamp DESC LIMIT 1`,
		unixTimestamp, eventUUID, userName).Scan(&ipAccess)

	return ipAccess
}

func GetSubsequentIPAccess(unixTimestamp int, eventUUID string,
	userName string) *IPAccess {
	db, err := gorm.Open("sqlite3", "supermanDetector.db")

	if err != nil {
		logger.Printf("unable to get subsequent IPAccess from database: %v", err)
	}
	defer db.Close()

	ipAccess := &IPAccess{}

	db.Raw(`SELECT * FROM ip_accesses WHERE ip_accesses.unix_timestamp > ?
		AND ip_accesses.event_uuid != ? AND ip_accesses.username = ?
		ORDER BY unix_timestamp ASC LIMIT 1`,
		unixTimestamp, eventUUID, userName).Scan(&ipAccess)

	return ipAccess
}
