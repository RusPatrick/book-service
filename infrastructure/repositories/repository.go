package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
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
	ErrNoSuchBook          = "Книги с таким id нет в хранилище"
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

	return &postgres{
		DB: db,
	}
}

func (db *postgres) GetDbConnect() *sql.DB {
	return db.DB
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

func (db *postgres) UpdateBook(book *models.Book) error {
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
	query := `SELECT id, author, title, publish_year, pages FROM Books WHERE id=$1;`
	row := db.DB.QueryRow(query, id)
	bookDTO := models.BookDTO{}
	if err := row.Scan(&bookDTO.ID, &bookDTO.Author, &bookDTO.Title, &bookDTO.PublishYear, &bookDTO.NumberPages); err != nil {
		return nil, errors.CreateDbError(err, ErrDbMsg)
	}
	if !bookDTO.ID.Valid {
		return nil, errors.CreateDbError(nil, ErrNoSuchBook)
	}
	book := bookDTO.ToEntity()
	return &book, nil
}

func getPostresConnectString(host, dbname, user, pass string) string {
	return fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable", dbname, host, user, pass)
}

func (db *postgres) RegisterUser(email, passHash string, passSalt string) error {
	query := `INSERT INTO users(email, password_hash, password_salt) VALUES ($1, $2, $3);`

	_, err := db.DB.Exec(query, email, passHash, passSalt)
	if err != nil {
		dbErr := err.(*pq.Error)
		if dbErr.Code == "23505" {
			return errors.CreateDbError(err, fmt.Sprintf(ErrThereAreAlreadyUser, email))
		}
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

func (db *postgres) SetSession(ID int, session_id string, exp int64) error {
	query := `INSERT INTO sessions VALUES($1, $2, $3)`

	if _, err := db.DB.Exec(query, ID, session_id, exp); err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}

	return nil
}

func (db *postgres) DeleteUser(sessionID string) error {
	selectQuery := `SELECT DISTINCT user_id FROM sessions WHERE session_id=$1`

	row := db.DB.QueryRow(selectQuery, sessionID)
	var userID int
	if err := row.Scan(&userID); err != nil {
		return err
	}

	deleteQuery := `DELETE FROM users WHERE id=$1 returning id`
	row = db.DB.QueryRow(deleteQuery, userID)

	if err := row.Scan(&userID); err != nil {
		return err
	}
	return nil
}

func (db *postgres) ClearSessions(exp int64) (int, error) {
	query := `DELETE FROM sessions WHERE exp<$1`

	result, err := db.DB.Exec(query, exp)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}
