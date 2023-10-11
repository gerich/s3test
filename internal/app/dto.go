package app

import (
	"errors"
	"net/http"
)

// File и FilesList по хорошему не место здесь
type (
	File struct {
		Name string `json:"name"`
	}
	FilesList struct {
		Files []*File `json:"files"`
		User  string  `json:"user"`
	}
	Error struct {
		error
		Status int
	}
)

var (
	ErrFileNotFound = &Error{error: errors.New("file not found"), Status: http.StatusNotFound}
	ErrFileExists   = &Error{error: errors.New("file exists"), Status: http.StatusConflict}
)
