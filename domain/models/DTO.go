package models

import "database/sql"

type UserDTO struct {
	ID       sql.NullInt64
	Email    sql.NullString
	PassHash sql.NullString
	PassSalt sql.NullString
}

func (dto *UserDTO) ToEntity() *UserDB {
	return &UserDB{
		ID:       &dto.ID.Int64,
		Email:    &dto.Email.String,
		PassHash: &dto.PassHash.String,
		PassSalt: &dto.PassSalt.String,
	}
}
