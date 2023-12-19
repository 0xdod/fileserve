package http_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/0xdod/fileserve"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUpload(t *testing.T) {
	s := startServer(t)
	defer s.Shutdown()

	mockFileStorage.On("Upload", mock.Anything, mock.Anything).Return("test-location", nil)
	mockFileService.On("CreateFile", mock.Anything, mock.Anything).Return(nil)

	url := "http://localhost:8080/api/v1/files/upload"

	resp, err := makeMultipartUploadRequest(t, url)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func makeMultipartUploadRequest(t testing.TB, url string) (*http.Response, error) {
	t.Helper()
	// Create a new buffer to write the multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Add a file to the form data
	file, err := os.Open("../testdata/file.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", "file.txt")

	if err != nil {
		fmt.Println("Error creating form file:", err)
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return http.DefaultClient.Do(req)
}

func TestHandleDownload(t *testing.T) {
	s := startServer(t)
	defer s.Shutdown()

	sampleFileData := []byte("test-data")
	mockFileStorage.On("Download", mock.Anything, mock.Anything).Return(sampleFileData, nil)
	mockFileService.On("GetFile", mock.Anything, mock.Anything).Return(&fileserve.File{
		ID:   "test-id",
		URL:  "test-url",
		Name: "test-name",
	}, nil)

	url := "http://localhost:8080/api/v1/files/download" + "/test-id"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, len(sampleFileData))
	resp.Body.Read(buffer)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/octet-stream", resp.Header.Get("Content-Type"))
	assert.Equal(t, sampleFileData, buffer)
}

func TestHandleGetFiles(t *testing.T) {
	s := startServer(t)

	defer s.Shutdown()

	now := time.Now()
	testFiles := []*fileserve.File{
		{
			ID:        "test-id",
			URL:       "test-url",
			Name:      "test-name",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "test-id",
			URL:       "test-url",
			Name:      "test-name",
			CreatedAt: now.AddDate(0, 0, -1),
			UpdatedAt: now.AddDate(0, 0, -1),
		},
		{
			ID:        "test-id",
			URL:       "test-url",
			Name:      "test-name",
			CreatedAt: now.AddDate(0, 0, -2),
			UpdatedAt: now.AddDate(0, 0, -2),
		},
	}

	mockFileService.On("GetFiles", mock.Anything, mock.Anything).Return(testFiles, nil)

	url := "http://localhost:8080/api/v1/files"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}
