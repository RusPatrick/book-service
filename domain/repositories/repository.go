package repositories

import (
	"github.com/ruspatrick/go-toff/domain/models"
)

type BooksRepository interface {
	AddBook(*models.Book) error
	ModifyBook(*models.Book) error
	DeleteBook(int) error
	GetBooksByTitle(string) ([]models.Book, error)
	GetBooksByYear(int) ([]models.Book, error)
	GetBooksByPages(minNumberOfPages, maxNumberOfPages int) ([]models.Book, error)
}
