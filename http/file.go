package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/0xdod/fileserve"
	"github.com/0xdod/fileserve/filestorage"
)

func (s *Server) handleUpload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Limit request size to prevent abuse
		r.ParseMultipartForm(32 << 20) // Limit request to 32MB (adjust as needed)

		file, handler, err := r.FormFile("file") // Form field name for the file
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
			return
		}

		defer file.Close()

		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		buffer := make([]byte, handler.Size)
		file.Read(buffer)
		loc, err := s.filestore.Upload(context.Background(), filestorage.UploadParam{
			Name:    handler.Filename,
			Content: buffer,
		})

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error Uploading the File", http.StatusInternalServerError)
			return
		}

		// insert into db
		newFile := fileserve.File{
			Name: handler.Filename,
			URL:  loc,
		}

		if err := s.fs.CreateFile(context.Background(), &newFile); err != nil {
			fmt.Println(err)
			http.Error(w, "Error Saving the File", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newFile)
	})
}
