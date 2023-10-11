package domain

import (
	"context"
	"errors"
)

type Repository interface {
	Save(context.Context, *File) error
	Get(context.Context, *File) ([]byte, error)
	ListByUser(context.Context, *User) ([]*File, error)
}

var (
	ErrAlreadyExists = errors.New("file already exists")
)
