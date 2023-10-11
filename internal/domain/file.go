package domain

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path/filepath"
)

const minFileSizeBytes = 20

type File struct {
	hash   string
	name   string
	data   []byte
	reader *bytes.Reader
	owner  *User
}

func NewFileFromBody(owner string, name string, data []byte) (*File, error) {
	hasher := md5.New()
	if _, err := hasher.Write(data); err != nil {
		return nil, err
	}

	if len(data) < minFileSizeBytes {
		return nil, ErrFileTooSmall
	}

	return &File{
			hash:   hex.EncodeToString(hasher.Sum(nil)),
			name:   name,
			data:   data,
			reader: bytes.NewReader(data),
			owner:  NewUser(owner)},
		nil
}

func NewFileFormHashAndName(owner, name, hash string) *File {
	return &File{hash: hash, owner: NewUser(owner), name: name}
}

func NewFileFormName(owner string, name string) *File {
	return &File{owner: NewUser(owner), name: name}
}

func (f *File) ID() string {
	return f.hash
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Data() []byte {
	return f.data
}

func (f *File) UserName() string {
	return f.owner.ID()
}

func (f *File) Size() int {
	return len(f.data)
}

func (f *File) Path() string {
	return filepath.Join(f.owner.ID(), f.name)
}

// Parts разбиение файла на части
func (f *File) Parts(count, number int) []byte {
	mod := len(f.data) % count
	partSize := (len(f.data) - mod) / count
	// Если последняя часть и на равные части не делиться
	tailBytes := 0
	if count == number+1 && mod != 0 {
		tailBytes = mod
	}
	if mod > count/2 {
		partSize += 1
		if tailBytes > 0 {
			tailBytes = mod - count
		}
	}

	res := make([]byte, partSize+tailBytes)

	n, err := f.reader.ReadAt(res, int64(partSize*number))
	if err == io.EOF {
		res = res[:n]
	}

	return res
}

var (
	ErrFileTooSmall  = fmt.Errorf("file too small, allowed size is %d bytes", minFileSizeBytes)
	ErrFileCorrupted = errors.New("file corrupted, please re-upload file")
)
