package filesystem

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gerich/s3test/internal/domain"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

// Repository TODO: unit tests
type Repository struct {
	root string
	cfg  *Config
	fs   afero.Fs
	log  *zap.Logger
}

type Config struct {
	Buckets []string
}

var _ domain.Repository = (*Repository)(nil)

func New(cfg *Config, log *zap.Logger) *Repository {
	return &Repository{cfg: cfg, fs: afero.NewOsFs(), log: log.Named("filesystem")}
}

func (r *Repository) Init() error {
	for _, b := range r.cfg.Buckets {
		log := r.log.With(zap.String("bucket", b))
		fi, err := r.fs.Stat(b)
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			return fmt.Errorf("%s is not a directory", b)
		}

		if fi.Mode().Perm()&0644 < 0644 {
			return fmt.Errorf("%s is not writable", b)
		}
		log.Debug("bucked checked")
	}

	return nil
}

// Save Сохранение файла на сервера
func (r *Repository) Save(_ context.Context, file *domain.File) error {
	for idx := range r.cfg.Buckets {
		userDir := filepath.Join(r.cfg.Buckets[idx], file.UserName())
		err := r.fs.Mkdir(userDir, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}

		filePath := filepath.Join(userDir, file.Name())
		// Save file part
		if err := r.saveFile(filePath, file.Parts(len(r.cfg.Buckets), idx)); err != nil {
			return err
		}

		// Save file part hash
		if err := r.saveFile(filePath+".hash", []byte(file.ID())); err != nil {
			return err
		}
	}

	return nil
}

// Get Чтение файла
func (r *Repository) Get(_ context.Context, file *domain.File) ([]byte, error) {
	exists := false
	buf := bytes.NewBuffer(nil)
	var lastHash []byte
	for idx := range r.cfg.Buckets {
		filePath := filepath.Join(r.cfg.Buckets[idx], file.Path())
		n, err := r.readFilePart(filePath, buf)
		if err != nil {
			return nil, err
		}

		if n == 0 {
			continue
		}
		exists = true

		currHash, err := r.readFileHash(filePath)
		if err != nil {
			return nil, err
		}

		if len(lastHash) > 0 && !bytes.Equal(currHash, lastHash) {
			return nil, domain.ErrFileCorrupted
		}
		currHash = lastHash
	}

	if !exists {
		return nil, nil
	}
	// 18399883
	return io.ReadAll(buf)
}

// ListByUser получение инфы по всем файлам пользователя
func (r *Repository) ListByUser(_ context.Context, user *domain.User) ([]*domain.File, error) {
	files := make(map[string]*domain.File)
	for idx := range r.cfg.Buckets {
		dirPath := filepath.Join(r.cfg.Buckets[idx], user.ID())
		dir, err := r.fs.Open(dirPath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		fi, err := dir.Stat()
		if err != nil {
			return nil, err
		}

		if !fi.IsDir() {
			return nil, fmt.Errorf("%s is not dir", fi.Name())
		}

		fsi, err := dir.Readdir(0)
		if err != nil {
			return nil, err
		}

		for _, fi := range fsi {
			if strings.HasSuffix(fi.Name(), ".hash") {
				continue
			}

			if _, ok := files[fi.Name()]; !ok {
				files[fi.Name()] = domain.NewFileFormName(user.ID(), fi.Name())
			}
		}
	}

	res := make([]*domain.File, 0, len(files))
	for _, f := range files {
		res = append(res, f)
	}

	return res, nil
}

func (r *Repository) saveFile(filePath string, data []byte) error {
	f, err := r.fs.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil && !os.IsNotExist(err) && !os.IsExist(err) {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	return f.Close()
}

func (r *Repository) readFilePart(filePath string, buf io.ReaderFrom) (int64, error) {
	f, err := r.fs.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}

	if os.IsNotExist(err) {
		return 0, nil
	}

	count, err := buf.ReadFrom(f)
	if err != nil {
		return 0, err
	}

	return count, f.Close()
}

func (r *Repository) readFileHash(filePath string) ([]byte, error) {
	f, err := r.fs.OpenFile(filePath+".hash", os.O_RDONLY, 0644)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if os.IsNotExist(err) {
		return nil, nil
	}

	res, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return res, f.Close()
}
