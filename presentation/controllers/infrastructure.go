package controllers

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"strconv"

	"github.com/ruspatrick/book-service/domain/models"

	"github.com/ruspatrick/book-service/infrastructure/errors"
)

func marshalJSON(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, errors.CreateServerError(err, errInternalServerErrorMsg)
	}
	return b, nil
}

func unmarshalJSON(data []byte, source interface{}) error {
	if err := json.Unmarshal(data, source); err != nil {
		return errors.CreateServerError(err, errInternalServerErrorMsg)
	}
	return nil
}

func parseQuery(values url.Values) models.BooksQuery {
	return models.BooksQuery{
		MinPages: parseIntOrSetDefault(values.Get("minPages"), 0),
		MaxPages: parseIntOrSetDefault(values.Get("maxPages"), 10000),
		MinYear:  parseIntOrSetDefault(values.Get("minYear"), -1900),
		MaxYear:  parseIntOrSetDefault(values.Get("maxYear"), 2030),
		SubStr:   values.Get("term"),
	}
}

func parseIntOrSetDefault(str string, defaultValue int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return value
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
