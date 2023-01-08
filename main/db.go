package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/go-sql-driver/mysql"
    "golang.org/x/crypto/bcrypt"

    schema "treehouse/schema"
)

var db *sql.DB

func initDB() {
    config := mysql.Config {
        User: os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net: "tcp",
        Addr: "127.0.0.1:3306",
        DBName: "treehouse",
        AllowNativePasswords: true,
    }

    var err error
    db, err = sql.Open("mysql", config.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }

    fmt.Println("Connected")
}

// TODO
func verifyUser(user schema.User) (schema.User, error) {
    return user, nil;
}

func addUser(user schema.User) (int64, error) {
    newUser, err := verifyUser(user)

    if err != nil {
        return 0, fmt.Errorf("addUser: %v", err)
    }

    pass := []byte(newUser.Password)
    hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)

    result, err := db.Exec(
        `insert into User (
            Username,
            Email,
            Password
        ) values (?, ?, ?, ?, ?)`,
        newUser.Username,
        newUser.Email,
        hashed,
    )

    if err != nil {
        return 0, fmt.Errorf("addUser: %v", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("addUser: %v", err)
    }

    fmt.Printf("NEW USER ADDED ID: %#v", id)

    return id, nil
}
