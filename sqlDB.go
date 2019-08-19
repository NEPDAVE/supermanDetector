package main

import (
	"fmt"
	gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//MigrateDB maps the Login struct fields to a newly created IPAccess SQL table
//where we can store IPAccesss
func MigrateDB() {
	db, err := gorm.Open("sqlite3", "supermanDetector.db")
	defer db.Close()

	//panicking if unable to successfully connect to the sqliteDB
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&IPAccess{})
}

//CreateIPAccess creates a new IPAccess entry in the IPAccess db table
func CreateIPAccess(ipAccess *IPAccess) {
	db, err := gorm.Open("sqlite3", "supermanDetector.db")
	defer db.Close()

	//panicking if unable to successfully connect to the sqliteDB
	if err != nil {
		panic(err)
	}

	db.Create(&ipAccess)
}

func QueryIPAccess(ipAccess *IPAccess) {
	db, err := gorm.Open("sqlite3", "supermanDetector.db")
	defer db.Close()

	//panicking if unable to successfully connect to the sqliteDB
	if err != nil {
		panic(err)
	}

	db.First(&ipAccess, "unix_timestamp = ?", 1514764800) // find product with code l1212

	fmt.Println(ipAccess)
}
