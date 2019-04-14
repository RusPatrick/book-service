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

type Error struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`s
}
