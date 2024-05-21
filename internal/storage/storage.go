package storage

import "errors"

var (
	ErrGenreNotFound = errors.New("genre not found")
	ErrPictureExists   = errors.New("picture already exists")
)