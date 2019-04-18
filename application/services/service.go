package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/ruspatrick/book-service/domain/models"
	repoI "github.com/ruspatrick/book-service/domain/repositories"
	"github.com/ruspatrick/book-service/infrastructure/errors"
	"github.com/ruspatrick/book-service/infrastructure/repositories"
)

var (
	repository       repoI.BooksRepository
	ErrWrongUserData = fmt.Errorf("Неверный email или пароль")
)

func Init() {
	repository = repositories.InitRepo()
}

func GetDB() *sql.DB {
	return repository.GetDbConnect()
}

func AddBook(book *models.Book) error {
	if err := validateBook(book); err != nil {
		return err
	}
	return repository.AddBook(book)
}

func UpdateBook(book *models.Book) error {
	if err := validateBook(book); err != nil {
		return err
	}
	return repository.UpdateBook(book)
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
	return repository.RegisterUser(userInfo.Email, passHash, salt)
}

func Login(userInfo models.User) (*http.Cookie, error) {
	userDB, err := repository.LoginUser(userInfo)
	if err != nil {
		return nil, err
	}
	curPassHash := hash(userInfo.Password + *userDB.PassSalt)
	dbPassHash := *userDB.PassHash
	if !reflect.DeepEqual(curPassHash, dbPassHash) {
		return nil, errors.CreateBusinessError(ErrWrongUserData, ErrWrongUserData.Error())
	}

	session_id := hash(userInfo.Email + *userDB.PassSalt + time.Now().String())
	exp := time.Now().Add(60 * time.Minute)
	if err := repository.SetSession(int(*userDB.ID), session_id, exp.Unix()); err != nil {
		return nil, err
	}
	return &http.Cookie{Name: "session_id", Value: session_id, Expires: exp}, nil
}

func DeleteUser(sessionID string) error {
	return repository.DeleteUser(sessionID)
}

func PeriodicalClearSessions() {
	ticker := time.NewTicker(time.Hour)
	for _ = range ticker.C {
		rowsDeleted, err := clearSessions()
		if err != nil {
			log.Println("Ошбка при чистке истекших сессий")
			return
		}
		log.Printf("Sessions deleted: %#+v\n", rowsDeleted)
	}
}

func clearSessions() (int, error) {
	currentTime := time.Now().Unix()
	return repository.ClearSessions(currentTime)
}

func hash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
