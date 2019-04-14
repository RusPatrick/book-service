package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ruspatrick/go-toff/presentation/core/config"

	_ "github.com/lib/pq"
	"github.com/ruspatrick/go-toff/domain/models"
	"github.com/ruspatrick/go-toff/domain/repositories"
)

type postgres struct {
	DB *sql.DB
}

const (
	MAX_DB_CONNECT   = 20
	ErrDbConnectText = "Failed set db connect: "
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

	db.SetMaxIdleConns(MAX_DB_CONNECT)
	db.SetMaxOpenConns(MAX_DB_CONNECT)

	booksDB = postgres{
		DB: db,
	}

	return &booksDB
}

func (db *postgres) AddBook(book *models.Book) error {
	rows, err := db.DB.Query(
		`INSERT INTO 
			Books(title, author, number_pages, publish_year) 
		VALUES 
			($1, $2, $3, $4) 
		RETURNING id`,
		book.Title,
		book.Author,
		book.NumberPages,
		book.PublishYear,
	)

	if err != nil {
		return err
	}

	if err := rows.Scan(book.ID); err != nil {
		return err
	}

	return nil
}

func (db *postgres) ModifyBook(book *models.Book) error {
	return nil
}

func (db *postgres) DeleteBook(id int) error {
	return nil
}

func (db *postgres) GetBooksByTitle(title string) ([]models.Book, error) {
	return nil, nil
}

func (db *postgres) GetBooksByYear(year int) ([]models.Book, error) {
	return nil, nil
}

func (db *postgres) GetBooksByPages(minNumberOfPages, maxNumberOfPages int) ([]models.Book, error) {
	return nil, nil
}

func getPostresConnectString(host, dbname, user, pass string) string {
	return fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable", dbname, host, user, pass)
}
