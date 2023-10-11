package app

import (
	"context"
	"fmt"

	"github.com/gerich/s3test/internal/domain"
	"go.uber.org/zap"
)

type Service struct {
	r   domain.Repository
	log *zap.Logger
}

func NewService(r domain.Repository, log *zap.Logger) *Service {
	return &Service{r: r, log: log.Named("service")}
}

func (s *Service) SaveFile(ctx context.Context, user string, fileName string, data []byte) (string, error) {
	log := s.log.With(zap.String("user", user), zap.String("file_name", fileName))
	s.log.Debug("save file")

	file, err := domain.NewFileFromBody(user, fileName, data)
	if err != nil {
		log.Error("cant prepare file for saving", zap.Error(err))
		return "", err
	}

	if err := s.r.Save(ctx, file); err != nil {
		log.Error("cant save file", zap.String("file_hash", file.ID()), zap.Error(err))
		return "", fmt.Errorf("cant save file: %w", err)
	}

	s.log.Info("file saved", zap.String("file_hash", file.ID()))
	return file.ID(), nil
}

func (s *Service) GetFile(ctx context.Context, user string, fileName string) ([]byte, error) {
	log := s.log.With(zap.String("user", user), zap.String("file_name", fileName))
	log.Debug("get file")

	params := domain.NewFileFormName(user, fileName)
	data, err := s.r.Get(ctx, params)
	if err != nil {
		s.log.Error("cant get file", zap.Error(err))
		return nil, fmt.Errorf("cant download file: %w", err)
	}

	return data, nil
}

func (s *Service) ListFilesByUser(ctx context.Context, userName string) (*FilesList, error) {
	log := s.log.With(zap.String("user", userName))
	log.Debug("list files by user")

	user := domain.NewUser(userName)
	files, err := s.r.ListByUser(ctx, user)
	if err != nil {
		s.log.Error("cant list files", zap.Error(err))
		return nil, fmt.Errorf("cant list files: %w", err)
	}

	res := &FilesList{Files: make([]*File, 0, len(files)), User: user.ID()}
	for _, f := range files {
		res.Files = append(res.Files, &File{Name: f.Name()})
	}

	return res, nil
}
