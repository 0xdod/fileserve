package http_test

import (
	"context"
	"testing"

	"github.com/0xdod/fileserve"
	fshttp "github.com/0xdod/fileserve/http"
	"github.com/0xdod/fileserve/sqlite"
	"github.com/stretchr/testify/mock"
)

type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) Upload(ctx context.Context, param fileserve.UploadParam) (string, error) {
	args := m.Called(ctx, param)
	return args.String(0), args.Error(1)
}

func (m *MockFileStorage) Download(ctx context.Context, name string) ([]byte, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]byte), args.Error(1)
}

type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) CreateFile(ctx context.Context, file *fileserve.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileService) GetFile(ctx context.Context, fileId string) (*fileserve.File, error) {
	args := m.Called(ctx, fileId)
	return args.Get(0).(*fileserve.File), args.Error(1)
}

func (m *MockFileService) GetFiles(ctx context.Context, param fileserve.GetFilesParam) ([]*fileserve.File, error) {
	args := m.Called(ctx, param)
	return args.Get(0).([]*fileserve.File), args.Error(1)
}

func (m *MockFileService) DeleteFile(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileService) UpdateFile(ctx context.Context, id string, param fileserve.UpdateFileParam) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

var mockFileStorage = new(MockFileStorage)
var mockFileService = new(MockFileService)

func startServer(tb testing.TB) *fshttp.Server {
	tb.Helper()

	s := fshttp.NewServer(fshttp.NewServerOpts{
		DB: &sqlite.DB{
			DSN: ":memory:",
		},
		FileStorage: mockFileStorage,
		FileService: mockFileService,
	})
	go s.Run()
	return s
}
