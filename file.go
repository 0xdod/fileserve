package fileserve

import (
	"context"
	"time"
)

type File struct {
	ID        string
	Name      string
	URL       string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type GetFilesParam struct {
	Name   *string
	Limit  *int
	Offset *int
}

type UpdateFileParam struct {
	Name *string
}

type FileService interface {
	CreateFile(ctx context.Context, file *File) error
	GetFiles(ctx context.Context, param GetFilesParam) ([]*File, error)
	GetFile(ctx context.Context, id string) (*File, error)
	UpdateFile(ctx context.Context, id string, param UpdateFileParam) error
	DeleteFile(ctx context.Context, id string) error
}
