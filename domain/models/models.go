package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Book struct {
	Author      string `json:"author"`
	Title       string `json:"title"`
	NumberPages int    `json:"numberPages"`
}
