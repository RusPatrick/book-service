package errors

import (
	"fmt"
	"net/http"
	"path"
	"runtime"
)

type Error struct {
	Type           string `json:"type"`
	Title          string `json:"title"`
	Details        string `json:"detail"`
	ErrorType      string `json:"errorType"`
	StatusCode     int    `json:"status"`
	Source         string `json:"-"`
	ErrDescription error  `json:"-"`
}

const (
	dbErrorTitle = "Ошибка базы данных"
	errorTitle   = "Ошибка"
)

func (err Error) Error() string {
	return err.Details
}

func CreateDbError(err error, details string) error {
	return Error{
		Type:           "about:blank",
		Title:          dbErrorTitle,
		Details:        details,
		ErrorType:      "server",
		StatusCode:     http.StatusInternalServerError,
		Source:         Where(),
		ErrDescription: err,
	}
}

func CreateServerError(err error, details string) error {
	return Error{
		Type:           "about:blank",
		Title:          errorTitle,
		Details:        details,
		ErrorType:      "server",
		StatusCode:     http.StatusInternalServerError,
		Source:         Where(),
		ErrDescription: err,
	}
}

func CreateBusinessError(err error, details string) error {
	return Error{
		Type:           "about:blank",
		Title:          errorTitle,
		Details:        details,
		ErrorType:      "business",
		StatusCode:     http.StatusBadRequest,
		Source:         Where(),
		ErrDescription: err,
	}
}

func Where() string {
	depth := 2
	function, file, line, _ := runtime.Caller(depth)
	_, fileName := path.Split(file)
	return fmt.Sprintf("File: %s Function: %s Line: %v", fileName, runtime.FuncForPC(function).Name(), line)
}
