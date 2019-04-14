package services

import (
	"github.com/ruspatrick/go-toff/domain/models"
	repoI "github.com/ruspatrick/go-toff/domain/repositories"
	"github.com/ruspatrick/go-toff/infrastructure/repositories"
)

var (
	booksRepo repoI.BooksRepository
	errInvalidBook = errors.New("Invalid ")
)

func Init() {
	booksRepo = repositories.InitRepo()
}

func AddBook(book *models.Book) error {+
	if validateBook(book) != 0 {
		return 
	}
	return booksRepo.AddBook(book)
}

func ModifyBook(book *models.Book) error {

	return booksRepo.ModifyBook(book)
}
