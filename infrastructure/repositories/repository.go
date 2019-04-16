package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ruspatrick/book-service/infrastructure/errors"

	"github.com/lib/pq"
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

func InitRepo() repositories.BooksRepository {
	var booksDB postgres
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
	query := `INSERT INTO Books(title, author, pages, publish_year) VALUES 	($1, $2, $3, $4) RETURNING id;`

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
	query := `UPDATE Books SET author=$2, title=$3, pages=$4, publish_year=$5 WHERE id=$1;`
	_, err := db.DB.Exec(query, book.ID, book.Author, book.Title, book.NumberPages, book.PublishYear)
	if err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}
	return nil
}

func (db *postgres) DeleteBook(id int) error {
	query := `DELETE FROM Books WHERE id=$1;`
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
		(title ilike '%%%s%%');
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
	query := `SELECT id, author, title, publish_year, pages FROM Books WHERE id=$1;`
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
	);
	`

	if _, err := db.DB.Exec(query); err != nil {
		log.Println(ErrDbMsg)
	}
}

func (db *postgres) createUsersTable() {
	//TODO
	db.DB.Exec(`CREATE TABLE IF NOT EXISTS users(
postgres(# id serial primary key,
postgres(# email varchar(100) not null unique,
postgres(# password_hash text not null,
postgres(# password_salt varchar(20) not null,
postgres(# session_id varchar(250),
postgres(# exp int);`)
}

func (db *postgres) RegisterUser(email, passHash, passSalt string) error {
	query := `INSERT INTO users(email, password_hash, password_salt) VALUES ($1, $2, $3);`

	_, err := db.DB.Exec(query, email, passHash, passSalt)
	dbErr := err.(*pq.Error)
	log.Println(dbErr.Message)
	if err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}

	return nil
}

func (db *postgres) LoginUser(userInfo models.User) (*models.UserDB, error) {
	query := `SELECT id, email, password_hash, password_salt FROM users WHERE email=$1`

	row := db.DB.QueryRow(query, userInfo.Email)
	userDTO := models.UserDTO{}
	if err := row.Scan(&userDTO.ID, &userDTO.Email, &userDTO.PassHash, &userDTO.PassSalt); err != nil {
		return nil, errors.CreateDbError(err, ErrDbMsg)
	}

	return userDTO.ToEntity(), nil
}

func (db *postgres) SetSession(email, session_id string, exp int64) error {
	query := `UPDATE Users SET session_id=$1, exp=$2 WHERE email=$3`

	if _, err := db.DB.Exec(query, session_id, exp, email); err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}

	return nil
}
