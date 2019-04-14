package services

import (
	"github.com/ruspatrick/go-toff/domain/models"
	"github.com/ruspatrick/go-toff/infrastructure/repositories"
)

var (
	booksRepo = repositories.InitRepo()
)

func AddBook(book *models.Book) error {
	return booksRepo.AddBook(book)
}

func ModifyBook(book *models.Book) error {
	return booksRepo.ModifyBook(book)
}
