package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	once sync.Once
)

func GetDB() *sql.DB {
	once.Do(func() {
		var err error

		connStr := "host=localhost port=5432 user=postgres password=admin dbname=postgres sslmode=disable"

		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Test connection
		if err = db.Ping(); err != nil {
			log.Fatalf("Database is unreachable: %v", err)
		}

		fmt.Println("✅ Connected to the database!")
	})
}
