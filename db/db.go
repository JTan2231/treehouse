package db

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var dbConn *sql.DB
var initialized = false

func InitDB() {
	config := mysql.Config{
		User:                 os.Getenv("MYSQLUSER"),
		Passwd:               os.Getenv("MYSQLPASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQLHOST") + ":" + os.Getenv("MYSQLPORT"),
		DBName:               os.Getenv("MYSQLDATABASE"),
		AllowNativePasswords: true,
	}

	var err error
	dbConn, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := dbConn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	initialized = true
	fmt.Println("Connected")
}

func GetDB() *sql.DB {
	if initialized {
		return dbConn
	} else {
		InitDB()
		return dbConn
	}
}
