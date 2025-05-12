package main

import (
	"fmt"
	"net/http"
	"os" // added for exiting the program
)

func main() {
	// Attempt to connect to the database
	db := GetDB()
	if db == nil {
		fmt.Println("Failed to connect to the database. Exiting...")
		os.Exit(1)
	}

	http.HandleFunc("/books/add", func(w http.ResponseWriter, r *http.Request) {
		saveBook(db, w, r)
	})
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		getBooks(db, w, r)
	})

	port := "8080"
	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
