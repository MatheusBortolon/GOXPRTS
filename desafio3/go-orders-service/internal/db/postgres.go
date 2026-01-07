package db

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

var db *sql.DB

func Connect(connStr string) {
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }

    fmt.Println("Successfully connected to the database")
}

func Disconnect() {
    if err := db.Close(); err != nil {
        log.Fatalf("Error closing the database: %v", err)
    }
    fmt.Println("Database connection closed")
}

func GetDB() *sql.DB {
    return db
}