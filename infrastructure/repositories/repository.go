package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ruspatrick/book-service/infrastructure/errors"

	_ "github.com/lib/pq"
	"github.com/ruspatrick/book-service/domain/models"
	"github.com/ruspatrick/book-service/domain/repositories"
	"github.com/ruspatrick/book-service/presentation/core/config"
)

type postgres struct {
	DB *sql.DB
}

const (
	maxDbConnect     = 20
	ErrDbConnectText = "Ошибка соединения с БД: "
	ErrDbMsg         = "Ошибка базы данных"
)

var (
	booksDB postgres
)

func InitRepo() repositories.BooksRepository {
	dbCfg := config.Get().DB
	db, err := sql.Open("postgres", getPostresConnectString(dbCfg.Host, dbCfg.Name, dbCfg.User, dbCfg.Password))
	if err != nil {
		log.Fatal(ErrDbConnectText, err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(ErrDbConnectText, err)
	}

	db.SetMaxIdleConns(maxDbConnect)
	db.SetMaxOpenConns(maxDbConnect)

	booksDB = postgres{
		DB: db,
	}

	booksDB.createBooksTable()
	booksDB.createUsersTable()

	return &booksDB
}

func (db *postgres) AddBook(book *models.Book) error {
	query := `INSERT INTO Books(title, author, pages, publish_year) VALUES 	($1, $2, $3, $4) RETURNING id`

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

func (db *postgres) ModifyBook(book *models.Book) error {
	query := `UPDATE Books SET author=$2, title=$3, pages=$4, publish_year=$5 WHERE id=$1`
	_, err := db.DB.Exec(query, book.ID, book.Author, book.Title, book.NumberPages, book.PublishYear)
	if err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}
	return nil
}

func (db *postgres) DeleteBook(id int) error {
	query := `DELETE FROM Books WHERE id=$1`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}
	return nil
}

func (db *postgres) GetBooksByFilters(subStr string, minPage int, maxPage int, minYaer int, maxYear int) ([]models.Book, error) {
	query := `SELECT * FROM Books WHERE 
		(%v is null or pages > %v) and 
		(%v is null or pages < %v) and 
		(%v is null or publish_year > %v) and 
		(%v is null or publish_year < %v) and 
		(title ilike '%%%s%%')
	`
	filledQuery := fmt.Sprintf(query, minPage, minPage, maxPage, maxPage, minYaer, minYaer, maxYear, maxYear, subStr)
	rows, err := db.DB.Query(filledQuery)
	if err != nil {
		return nil, errors.CreateDbError(err, ErrDbMsg)
	}

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
	query := `SELECT id, author, title, publish_year, pages FROM Books WHERE id=$1`
	row := db.DB.QueryRow(query, id)
	book := models.Book{}
	if err := row.Scan(&book.ID, &book.Author, &book.Title, &book.PublishYear, &book.NumberPages); err != nil {
		return nil, errors.CreateDbError(err, ErrDbMsg)
	}
	return &book, nil
}

func getPostresConnectString(host, dbname, user, pass string) string {
	return fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable", dbname, host, user, pass)
}

func (db *postgres) createBooksTable() {
	query := `
	CREATE TABLE IF NOT EXISTS Books (
		id serial primary key,
		author varchar(100) not null,
		title varchar(50) not null,
		pages int not null,
		publish_year int not null
	)
	`

	if _, err := db.DB.Exec(query); err != nil {
		log.Println(ErrDbMsg)
	}
}

func (db *postgres) createUsersTable() {

}
