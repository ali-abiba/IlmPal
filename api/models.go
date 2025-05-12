package main

import (
	"time"
)

type Book struct {
	ID         int            `json:"id"`
	Title      string         `json:"title"`
	Author     string         `json:"author"`
	Content    string         `json:"content"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	Categories []BookCategory `json:"categories"`
}

type BookCategory struct {
	ID       int      `json:"id"`
	Category Category `json:"category"`
	Weight   float64  `json:"weight"`
}

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID            int        `json:"id"`
	Email         string     `json:"email"`
	TopCategories []Category `json:"top_categories"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Reading struct {
	ID           int       `json:"id"`
	User         []User    `json:"User"`
	Book         []Book    `json:"Book"`
	Notes        string    `json:"notes"`
	Started      time.Time `json:"started"`
	LastRead     time.Time `json:"last_read"`
	BookmarkPage int       `json:"Bookmark_page"`
}
