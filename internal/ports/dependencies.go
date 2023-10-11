package ports

import (
	"context"

	"github.com/gerich/s3test/internal/app"
)

var _ App = (*app.Service)(nil)

type App interface {
	SaveFile(context.Context, string, string, []byte) (string, error)
	GetFile(context.Context, string, string) ([]byte, error)
	ListFilesByUser(context.Context, string) (*app.FilesList, error)
}
