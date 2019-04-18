package repositories

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/ruspatrick/book-service/domain/models"
	"github.com/ruspatrick/book-service/infrastructure/errors"
)

func (db *postgres) RegisterUser(email, passHash string, passSalt string) error {
	query := `INSERT INTO books.users(email, password_hash, password_salt) VALUES ($1, $2, $3);`

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
	query := `SELECT id, email, password_hash, password_salt FROM books.users WHERE email=$1`

	row := db.DB.QueryRow(query, userInfo.Email)
	userDTO := models.UserDTO{}

	if err := row.Scan(&userDTO.ID, &userDTO.Email, &userDTO.PassHash, &userDTO.PassSalt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.CreateBusinessError(ErrNoSuchUser, "Неверный email или пароль")
		}
		return nil, errors.CreateDbError(err, ErrDbMsg)
	}
	if !userDTO.ID.Valid {
		return nil, errors.CreateBusinessError(ErrNoSuchUser, ErrNoSuchData)
	}

	return userDTO.ToEntity(), nil
}

func (db *postgres) SetSession(userID int, sessionID string, exp int64) error {
	query := `INSERT INTO books.sessions(user_id, session_id, exp) VALUES($1, $2, $3)`

	if _, err := db.DB.Exec(query, userID, sessionID, exp); err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}

	return nil
}

func (db *postgres) DeleteUser(sessionID string) error {
	deleteQuery := `DELETE FROM books.users WHERE id=(SELECT id FROM books.sessions WHERE session_id=$1)`
	if _, err := db.DB.Exec(deleteQuery, sessionID); err != nil {
		return errors.CreateDbError(err, ErrDbMsg)
	}
	return nil
}

func (db *postgres) ClearSessions(exp int64) (int, error) {
	query := `DELETE FROM books.sessions WHERE exp<$1`

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
