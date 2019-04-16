package repositories

import (
	"github.com/ruspatrick/book-service/domain/models"
)

type BooksRepository interface {
	AddBook(*models.Book) error
	ModifyBook(*models.Book) error
	DeleteBook(int) error
	GetBooksByFilters(subStr string, minPage int, maxPage int, minYaer int, maxYear int) ([]models.Book, error)
	GetBookByID(id int) (*models.Book, error)
	RegisterUser(username, passHash, passSalt string) error
	LoginUser(userInfo models.User) (*models.UserDB, error)
	SetSession(email, session_id string, exp int64) error
}
