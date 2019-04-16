package services

import (
	"crypto/sha512"
	"net/http"
	"reflect"
	"time"

	"github.com/ruspatrick/book-service/domain/models"
	repoI "github.com/ruspatrick/book-service/domain/repositories"
	"github.com/ruspatrick/book-service/infrastructure/errors"
	"github.com/ruspatrick/book-service/infrastructure/repositories"
)

const (
	ErrWrongPassword = "Неверный пароль"
)

var (
	repository repoI.BooksRepository
)

func Init() {
	repository = repositories.InitRepo()
}

func AddBook(book *models.Book) error {
	if err := validateBook(book); err != nil {
		return err
	}
	return repository.AddBook(book)
}

func ModifyBook(book *models.Book) error {
	if err := validateBook(book); err != nil {
		return err
	}
	return repository.ModifyBook(book)
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

	return repository.GetBooksByFilters(booksQuery.SubStr, booksQuery.MinPages, booksQuery.MaxPages, booksQuery.MinYear, booksQuery.MaxYear)
}

func GetBookByID(bookID int) (*models.Book, error) {
	return repository.GetBookByID(bookID)
}

func DeleteBook(bookID int) error {
	return repository.DeleteBook(bookID)
}

func Signup(userInfo models.User, salt string) error {
	passHash := hash(userInfo.Password + salt)
	return repository.RegisterUser(userInfo.Email, string(passHash), salt)
}

func Login(userInfo models.User) (*http.Cookie, error) {
	userDB, err := repository.LoginUser(userInfo)
	if err != nil {

	}
	if !reflect.DeepEqual(hash(userInfo.Password+*userDB.PassSalt), userDB.PassHash) {
		return nil, errors.CreateBusinessError(err, ErrWrongPassword)
	}

	session_id := hash(userInfo.Email + *userDB.PassSalt + time.Now().String())
	exp := time.Now().Add(60 * time.Minute)
	if err := repository.SetSession(userInfo.Email, session_id, exp.Unix()); err != nil {
		return nil, err
	}
	return &http.Cookie{Name: "session_id", Value: session_id, Expires: exp}, nil
}

func hash(str string) string {
	hasher := sha512.New()
	hasher.Write([]byte(str))
	return string(hasher.Sum(nil))
}
