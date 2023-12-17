package filestorage

import "context"

type UploadParam struct {
	Name    string
	Content []byte
}

type FileStorage interface {
	Upload(ctx context.Context, param UploadParam) (location string, err error)
	Download(ctx context.Context, name string) ([]byte, error)
}
