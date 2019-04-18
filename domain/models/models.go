package models

type User struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserDB struct {
	ID       *int64
	Email    *string
	PassHash *string
	PassSalt *string
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
	ID  string
	Exp int64
}
