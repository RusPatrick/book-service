package repositories

import (
	"database/sql"

	"github.com/ruspatrick/go-toff/domain/models"
	"github.com/ruspatrick/go-toff/domain/repositories"
)

type postgres struct {
	DB *sql.DB
}

var (
	booksDB postgres
)

func InitRepo() repositories.BooksRepository {
	return &booksDB
}

func (db *postgres) AddBook(*models.Book) error {
	return nil
}

func (db *postgres) ModifyBook(*models.Book) error {
	return nil
}

func (db *postgres) DeleteBook(int) error {
	return nil
}

func (db *postgres) GetBooksByTitle(string) ([]models.Book, error) {
	return nil, nil
}

func (db *postgres) GetBooksByYear(int) ([]models.Book, error) {
	return nil, nil
}

func (db *postgres) GetBooksByPages(minNumberOfPages, maxNumberOfPages int) ([]models.Book, error) {
	return nil, nil
}
