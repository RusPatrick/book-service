package repositories

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/ruspatrick/book-service/domain/models"
	"github.com/ruspatrick/book-service/domain/repositories"
	"github.com/ruspatrick/book-service/infrastructure/errors"
	"github.com/ruspatrick/book-service/presentation/core/config"
)

type postgres struct {
	DB *sql.DB
}

const (
	maxDbConnect           = 20
	ErrDbConnectText       = "Ошибка соединения с БД: "
	ErrDbMsg               = "Ошибка базы данных"
	ErrThereAreAlreadyUser = "Пользователь с email: '%s' уже зарегестрирован"
	ErrNoSuchData          = "Данные в БД отсутствуют"
)

var (
	ErrNoSuchUser = fmt.Errorf("Пользователь отсутствует в БД")
	ErrNoSuchBook = fmt.Errorf("Книга отсутствует в БД")
)

func InitRepo() repositories.BooksRepository {
	dbCfg := config.Get().DB
	connStr := getPostresConnectString(dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.User, dbCfg.Password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(ErrDbConnectText, err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(ErrDbConnectText, err)
	}

	db.SetMaxIdleConns(maxDbConnect)
	db.SetMaxOpenConns(maxDbConnect)

	return &postgres{
		DB: db,
	}
}

func (db *postgres) GetDbConnect() *sql.DB {
	return db.DB
}

func (db *postgres) AddBook(book *models.Book) error {
	query := `INSERT INTO books.books(title, author, pages, publish_year) VALUES 	($1, $2, $3, $4) RETURNING id;`

	row := db.DB.QueryRow(query,
		book.Title,
		book.Author,
		book.NumberPages,
		book.PublishYear,
	)

	if err := row.Scan(&book.ID); err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}

	return nil
}

func (db *postgres) UpdateBook(book *models.Book) error {
	query := `UPDATE books.books SET author=$2, title=$3, pages=$4, publish_year=$5 WHERE id=$1;`
	_, err := db.DB.Exec(query, book.ID, book.Author, book.Title, book.NumberPages, book.PublishYear)
	if err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}
	return nil
}

func (db *postgres) DeleteBook(id int) error {
	query := `DELETE FROM books.books WHERE id=$1;`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}
	return nil
}

func (db *postgres) GetBooksByFilters(subStr string, minPage int, maxPage int, minYaer int, maxYear int) ([]models.Book, error) {
	query := `SELECT * FROM books.books WHERE 
		(%v is null or pages > %v) and 
		(%v is null or pages < %v) and 
		(%v is null or publish_year > %v) and 
		(%v is null or publish_year < %v) and 
		(title ilike '%%%s%%');
	`
	filledQuery := fmt.Sprintf(query, minPage, minPage, maxPage, maxPage, minYaer, minYaer, maxYear, maxYear, subStr)
	rows, err := db.DB.Query(filledQuery)
	if err != nil {
		return nil, errors.CreateDbError(err, ErrDbMsg)
	}
	defer rows.Close()

	books := make([]models.Book, 0)

	for rows.Next() {
		book := models.Book{}
		if err := rows.Scan(&book.ID, &book.Author, &book.Title, &book.NumberPages, &book.PublishYear); err != nil {
			return nil, errors.CreateDbError(err, ErrDbMsg)
		}

		books = append(books, book)
	}

	return books, nil
}

func (db *postgres) GetBookByID(id int) (*models.Book, error) {
	query := `SELECT id, author, title, publish_year, pages FROM books.books WHERE id=$1;`
	row := db.DB.QueryRow(query, id)
	bookDTO := models.BookDTO{}
	if err := row.Scan(&bookDTO.ID, &bookDTO.Author, &bookDTO.Title, &bookDTO.PublishYear, &bookDTO.NumberPages); err != nil {
		return nil, errors.CreateDbError(err, ErrDbMsg)
	}
	if !bookDTO.ID.Valid {
		return nil, errors.CreateDbError(ErrNoSuchBook, ErrNoSuchData)
	}
	book := bookDTO.ToEntity()
	return &book, nil
}

func getPostresConnectString(host, port, dbname, user, pass string) string {
	return fmt.Sprintf("dbname=%s host=%s port=%s user=%s password=%s sslmode=disable", dbname, host, port, user, pass)
}
