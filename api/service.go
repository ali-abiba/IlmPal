package main

import (
	"net/http"
	"encoding/json"
)

func saveBook(w http.ResponseWriter, r *httpRequest) {
	db := GetDB()
	var newBook Book

	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO book (title, content, author) VAULES ($1, $2, $3)", newBook.title, newBook.content, newBook.author);

	if err != nil {
		http.Error(w, "Failed to insert book", http.StatusInternalServerError);
		return
	}

	query := "INSERT INTO category (name) VAULES "
	values := []interface{}{}
	for i, category := range newBook.categories {
		query += fmt.Sprintf("(%d), ", i)

		values = append(values, category)
	}

	query = query[:len(query)-1]

	_, err := db.Exec(query, values...)
	if err != nil {
		http.Error(w, "Failed to insert categories", http.StatusInternalServerError);
		return
	}
}