package main

import (
    "fmt"
    "net/http"
)


type book struct {
    ID int `json:"id"`
    Title string `json:"title"`
    Author string `json:"author"`
    Content string `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Categories []category `json:"categories"`
}

type category struct {
    ID int `json:"id"`
    Name string `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

type user struct {
    ID int `json:"id"`
    Email string `json:"email"`
    TopCategories []category `json:"top_categories"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type reading struct {
    ID int `json:"id"`
    User []user `json:"user"`
    Book []book `json:"book"`
    Notes string `json:"notes`
    Started time.Time `json:"started"`
    LastRead time.Time `json:"last_read`
    BookmarkPage int `json:"bookmark_page"`
}
    
func main() {
    http.HandleFunc();
    http.HandleFunc();

    port := "8080"
    fmt.Printf("Server is running on port %s\n", port)
    http.ListenAndServe(":"+port, nil)
}