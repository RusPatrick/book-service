package models

type Config struct {
	App app `json:"application"`
	DB  db  `json:"db"`
}

type app struct {
	Name    string `json:"name"`
	Port    string `json:"port"`
	Version string `json:"version"`
}

type db struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `jspn:"port"`
	User     string `json:"user"`
	Password string `json:"pass"`
}
