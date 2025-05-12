package main

import (
	"database/sql" // added import for sql package
	"encoding/json"
	"fmt"
	"log" // added import for logging
	"net/http"
	"strconv"
	"strings" // added import for string operations
)

type getParameters struct {
	Sort       sort
	PageSize   int
	PageNumber int
	Filters    []filter
}

type sort struct {
	Field string
	Order string
}

type filter struct {
	Field    string
	Operator string
	Value    string
}

var whiteListedColumns = map[string]string{
	// Book table columns
	"id":         "id",
	"title":      "title",
	"content":    "content",
	"author":     "author",
	"created_at": "created_at",
	"updated_at": "updated_at",

	// Category table columns
	"name": "name",

	// Book_category table columns
	"book_id":     "book_id",
	"category_id": "category_id",
	"weight":      "weight",

	// User table columns
	"email": "email",

	// Top_categories table columns
	"user_id": "user_id",

	// Reading table columns
	"notes":           "notes",
	"started_reading": "started_reading",
	"last_reading":    "last_reading",
	"bookmark_page":   "bookmark_page",
}

func parseGetParameters(r *http.Request) getParameters {
	log.Println("Parsing GET parameters")
	sortField := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")
	if sortField == "" {
		sortField = "id"
	}

	if sortOrder == "" {
		sortOrder = "asc"
	}

	sort := sort{
		Field: sortField,
		Order: sortOrder,
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 10 // Default page size if parsing fails
	}
	if pageSize <= 0 {
		pageSize = 10 // Default page size if value is invalid
	}

	pageNumber, err := strconv.Atoi(r.URL.Query().Get("pageNumber"))

	filters := []filter{}
	for key, value := range r.URL.Query() {
		if key != "sort" && key != "order" && key != "pageSize" && key != "pageNumber" {
			switch key {
			case "[gt]":
				filters = append(filters, filter{Field: key, Operator: ">", Value: value[0]})
			case "[lt]":
				filters = append(filters, filter{Field: key, Operator: "<", Value: value[0]})
			case "[gte]":
				filters = append(filters, filter{Field: key, Operator: ">=", Value: value[0]})
			case "[lte]":
				filters = append(filters, filter{Field: key, Operator: "<=", Value: value[0]})
			case "[eq]":
				filters = append(filters, filter{Field: key, Operator: "=", Value: value[0]})
			case "[has]":
				filters = append(filters, filter{Field: key, Operator: "LIKE", Value: "%" + value[0] + "%"})
			}
		}
	}

	log.Printf("Parsed parameters: %+v\n", getParameters{Sort: sort, PageSize: pageSize, Filters: filters, PageNumber: pageNumber})
	return getParameters{Sort: sort, PageSize: pageSize, Filters: filters, PageNumber: pageNumber}
}

func parseFiltersToQuery(filters []filter) string {
	log.Println("Parsing filters to query")
	allowedFilters := []string{}
	for _, f := range filters {
		if col, ok := whiteListedColumns[f.Field]; ok {
			allowedFilters = append(allowedFilters, fmt.Sprintf(" %s %s %s ", col, f.Operator, f.Value))
		}
	}
	query := ""
	if len(allowedFilters) > 0 {
		query = " WHERE " + strings.Join(allowedFilters, " AND ")
	}
	log.Printf("Generated query: %s\n", query)
	return query
}

func saveBook(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("Saving a new book")
	var newBook Book

	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v\n", err)
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Insert book and retrieve its ID
	err = tx.QueryRow("INSERT INTO book (title, content, author) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id", newBook.Title, newBook.Content, newBook.Author).Scan(&newBook.ID)
	if err != nil {
		log.Printf("Error inserting book: %v\n", err)
		tx.Rollback()
		http.Error(w, "Failed to insert book", http.StatusInternalServerError)
		return
	}

	// Insert categories
	log.Println("Inserting categories")
	query := "INSERT INTO category (name) VALUES "
	values := []interface{}{}
	for i, category := range newBook.Categories {
		query += fmt.Sprintf("($%d), ", i+1)
		values = append(values, category.Category.Name)
	}
	query = query[:len(query)-2] // Remove the trailing comma and space
	query += "ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id, name"

	var new_categories []Category
	rows, err := tx.Query(query, values...)
	if err != nil {
		log.Printf("Error inserting categories: %v\n", err)
		tx.Rollback()
		http.Error(w, "Failed to insert categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			log.Printf("Error scanning category: %v\n", err)
			tx.Rollback()
			http.Error(w, "Failed to scan category", http.StatusInternalServerError)
			return
		}
		new_categories = append(new_categories, cat)
	}

	for i, bookCategory := range newBook.Categories {
		for _, category := range new_categories {
			if bookCategory.Category.Name == category.Name {
				newBook.Categories[i].Category = category
				break
			}
		}
	}

	// Insert book_category
	log.Println("Inserting book categories")
	query = "INSERT INTO book_category (book_id, category_id, weight) VALUES "
	values = []interface{}{}
	for i, category := range newBook.Categories {
		query += fmt.Sprintf("($%d, $%d, $%d), ", i*3+1, i*3+2, i*3+3)
		values = append(values, newBook.ID, category.Category.ID, category.Weight)
	}
	query = query[:len(query)-2]
	query += "ON CONFLICT DO NOTHING RETURNING id"

	_, err = tx.Exec(query, values...)
	if err != nil {
		log.Printf("Error inserting book categories: %v\n", err)
		tx.Rollback()
		http.Error(w, "Failed to insert book categories", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	log.Println("Book saved successfully")
	w.WriteHeader(http.StatusCreated)
}

func getBooks(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching books")
	params := parseGetParameters(r)

	// Base query to fetch books and their categories
	query := `
		SELECT 
			b.id, b.title, b.content, b.author, 
			COALESCE(json_agg(json_build_object('id', c.id, 'category', json_build_object('id', c.id, 'name', c."name", 'created_at', c.created_at), 'weight', bc.weight)) FILTER (WHERE c.id IS NOT NULL), '[]') AS categories
		FROM book b
		LEFT JOIN book_category bc ON b.id = bc.book_id
		LEFT JOIN category c ON bc.category_id = c.id
		GROUP BY b.id, b.title, b.content, b.author
	`

	// Apply filters if any
	if len(params.Filters) > 0 {
		log.Println("Applying filters to query")
		query += parseFiltersToQuery(params.Filters)
	}

	// Add sorting
	log.Printf("Sorting by %s %s\n", params.Sort.Field, params.Sort.Order)
	query += fmt.Sprintf(" ORDER BY %s %s", params.Sort.Field, params.Sort.Order)

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching books: %v\n", err)
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parse results
	var books []Book
	for rows.Next() {
		var book Book
		var categoriesJSON string
		if err := rows.Scan(&book.ID, &book.Title, &book.Content, &book.Author, &categoriesJSON); err != nil {
			log.Printf("Error scanning book: %v\n", err)
			http.Error(w, "Failed to scan book", http.StatusInternalServerError)
			return
		}

		// Decode categories JSON into the Book struct
		if err := json.Unmarshal([]byte(categoriesJSON), &book.Categories); err != nil {
			log.Printf("Error parsing categories JSON: %v\n", err)
			http.Error(w, "Failed to parse categories", http.StatusInternalServerError)
			return
		}

		books = append(books, book)
	}

	log.Printf("Fetched %d books\n", len(books))
	// Respond with books
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
