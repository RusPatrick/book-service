package repositories

import (
	"database/sql"

	"github.com/ruspatrick/book-service/domain/models"
)

type BooksRepository interface {
	GetDbConnect() *sql.DB
	AddBook(*models.Book) error
	UpdateBook(*models.Book) error
	DeleteBook(int) error
	GetBooksByFilters(subStr string, minPage int, maxPage int, minYaer int, maxYear int) ([]models.Book, error)
	GetBookByID(id int) (*models.Book, error)
	RegisterUser(username, passHash string, passSalt string) error
	LoginUser(userInfo models.User) (*models.UserDB, error)
	SetSession(ID int, sessionID string, exp int64) error
	DeleteUser(sessionID string) error
	ClearSessions(exp int64) (int, error)
}
