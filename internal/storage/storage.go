package storage

import "errors"

var (
	ErrUrlExist    = errors.New("url already exist")
	ErrUrlNotFound = errors.New("url not found")
)
