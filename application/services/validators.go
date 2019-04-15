package services

import (
	"fmt"

	"github.com/ruspatrick/book-service/domain/models"
	"github.com/ruspatrick/book-service/infrastructure/errors"
)

const (
	minPages       = 0
	maxPages       = 10000
	lengthTitle    = 50
	minPublishyear = -1900
)

var (
	ErrIncorrectTitle       = fmt.Errorf("максимально допустимое название 50 символов")
	ErrIncorrectNumberPages = fmt.Errorf("количество страниц должно быть в диапазоне 0 - 10000 страниц")
	ErrInvorrectPublishYear = fmt.Errorf("год издания не может быть меньше 1900г до н.э")
)

func validateBook(book *models.Book) error {
	if err := validateTitle(book.Title); err != nil {
		return err
	}
	if err := validatePages(book.NumberPages); err != nil {
		return err
	}
	if err := validatePublishYear(book.PublishYear); err != nil {
		return err
	}
	return nil
}

func validateTitle(title string) error {
	if len(title) > lengthTitle {
		return errors.CreateBusinessError(ErrIncorrectTitle, ErrIncorrectTitle.Error())
	}
	return nil
}

func validatePages(numberPages int) error {
	if numberPages < minPages || numberPages > maxPages {
		return errors.CreateBusinessError(ErrIncorrectNumberPages, ErrIncorrectNumberPages.Error())
	}
	return nil
}

func validatePublishYear(publishYear int) error {
	if publishYear < minPublishyear {
		return errors.CreateBusinessError(ErrInvorrectPublishYear, ErrInvorrectPublishYear.Error())
	}
	return nil
}
