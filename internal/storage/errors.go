package storage

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrNoKeyFound     = errors.New("no keys available")
	ErrCardIsInvalid  = errors.New("card information is invalid")
)
