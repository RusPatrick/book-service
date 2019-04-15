package services

import (
	"errors"

	"github.com/ruspatrick/book-service/domain/models"
	repoI "github.com/ruspatrick/book-service/domain/repositories"
	"github.com/ruspatrick/book-service/infrastructure/repositories"
)

var (
	booksRepo      repoI.BooksRepository
	errInvalidBook = errors.New("Invalid ")
)

func Init() {
	booksRepo = repositories.InitRepo()
}

func AddBook(book *models.Book) error {
	if err := validateBook(book); err != nil {
		return err
	}
	return booksRepo.AddBook(book)
}

func ModifyBook(book *models.Book) error {
	if err := validateBook(book); err != nil {
		return err
	}
	return booksRepo.ModifyBook(book)
}

func GetBooks(booksQuery models.BooksQuery) ([]models.Book, error) {
	if err := validateTitle(booksQuery.SubStr); err != nil {
		return nil, err
	}
	if err := validatePublishYear(booksQuery.MinYear); err != nil {
		return nil, err
	}
	if err := validatePublishYear(booksQuery.MaxYear); err != nil {
		return nil, err
	}
	if err := validatePages(booksQuery.MinPages); err != nil {
		return nil, err
	}
	if err := validatePages(booksQuery.MaxPages); err != nil {
		return nil, err
	}

	return booksRepo.GetBooksByFilters(booksQuery.SubStr, booksQuery.MinPages, booksQuery.MaxPages, booksQuery.MinYear, booksQuery.MaxYear)
}

func GetBookByID(bookID int) (*models.Book, error) {
	return booksRepo.GetBookByID(bookID)
}

func DeleteBook(bookID int) error {
	return booksRepo.DeleteBook(bookID)
}
