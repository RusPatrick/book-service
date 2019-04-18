package models

import "database/sql"

type UserDTO struct {
	ID       sql.NullInt64
	Email    sql.NullString
	PassHash sql.NullString
	PassSalt sql.NullString
}

type BookDTO struct {
	ID          sql.NullInt64
	Author      sql.NullString
	Title       sql.NullString
	PublishYear sql.NullInt64
	NumberPages sql.NullInt64
}

func (dto *UserDTO) ToEntity() *UserDB {
	return &UserDB{
		ID:       &dto.ID.Int64,
		Email:    &dto.Email.String,
		PassHash: &dto.PassHash.String,
		PassSalt: &dto.PassSalt.String,
	}
}
func (dto *BookDTO) ToEntity() Book {
	return Book{
		ID:          int(dto.ID.Int64),
		Author:      dto.Author.String,
		Title:       dto.Title.String,
		PublishYear: int(dto.PublishYear.Int64),
		NumberPages: int(dto.NumberPages.Int64),
	}
}
