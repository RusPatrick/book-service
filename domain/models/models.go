package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Book struct {
	ID          int    `json:"id"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	NumberPages int    `json:"number_pages"`
	PublishYear int    `json:"publish_year"`
}

type BooksQuery struct {
	MinYear  int
	MaxYear  int
	SubStr   string
	MinPages int
	MaxPages int
}

type Session struct {
	ID string `json:"session_id"`
}