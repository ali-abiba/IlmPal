package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func GetDB() *sql.DB {
	once.Do(func() {
		var err error

		connStr := "host=localhost port=5432 user=postgres password=admin dbname=product sslmode=disable"

		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Failed to connect to database: %v\n", err)
			db = nil
			return
		}

		// Test connection
		if err = db.Ping(); err != nil {
			log.Printf("Database is unreachable: %v\n", err)
			db = nil
			return
		}

		fmt.Println("✅ Connected to the database!")
	})

	return db
}
