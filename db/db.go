package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB
var initialized = false

func InitDB() {
    config := mysql.Config {
        User: os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net: "tcp",
        Addr: "127.0.0.1:3306",
        DBName: "treehouse",
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
