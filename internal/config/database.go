package config

import (
	"database/sql"
	"log"
	"sync"
    "time"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DB  *sql.DB
	once sync.Once
)

func GetDbInstance() (*sql.DB,error) {
	var err error
	once.Do(func() {
		dsn := DSN
		var db *sql.DB
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}

		db.SetMaxOpenConns(50)                 // Max 50 connections open at once
		db.SetMaxIdleConns(25)                 // Keep 25 idle (standby) connections
		db.SetConnMaxLifetime(5 * time.Minute) // Refresh connections every 5 minutes
        
		// Test the connection
		err = db.Ping()
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}
		DB=db
	})
	return DB,err
}
