package fileserve

import (
	"context"
	"time"
)

type File struct {
	ID        string
	Name      string
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FileUploadParam struct {
	Name    string
	Content []byte
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
	UploadFile(ctx context.Context, param FileUploadParam) (*File, error)
	DownloadFile(ctx context.Context, file File) ([]byte, error)
	GetFiles(ctx context.Context, param GetFilesParam) ([]*File, error)
	GetFile(ctx context.Context, id string) (*File, error)
	UpdateFile(ctx context.Context, id string, param UpdateFileParam) error
	// DeleteFile(ctx context.Context, file File) error
}
