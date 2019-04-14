package services

import (
	"github.com/ruspatrick/go-toff/domain/models"
)

const (
	minPages       = 0
	maxPages       = 10000
	lengthTitle    = 50
	minPublishyear = -1900
)

func validateBook(book *models.Book) int {
	var failedValidation int
	failedValidation += validateTitle(book.Title)
	failedValidation += validatePages(book.NumberPages)
	failedValidation += validatePublishYear(book.PublishYear)
	return failedValidation
}

func validateTitle(title string) int {
	if len(title) > lengthTitle {
		return 1
	}
	return 0
}

func validatePages(numberPages int) int {
	if numberPages < minPages || numberPages > maxPages {
		return 1
	}
	return 0
}

func validatePublishYear(publishYear int) int {
	if publishYear < minPublishyear {
		return 1
	}
	return 0
}
