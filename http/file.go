package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/0xdod/fileserve"
	"github.com/0xdod/fileserve/filestorage"
	"github.com/gorilla/mux"
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
		loc, err := s.fileStorage.Upload(context.Background(), filestorage.UploadParam{
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

		if err := s.fileService.CreateFile(context.Background(), &newFile); err != nil {
			fmt.Println(err)
			http.Error(w, "Error Saving the File", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newFile)
	})
}

func (s *Server) handleDownload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileId := vars["fileId"]

		file, err := s.fileService.GetFile(context.Background(), fileId)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error Retrieving the File", http.StatusInternalServerError)
			return
		}

		content, err := s.fileStorage.Download(context.Background(), file.Name)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error Retrieving the File", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
		w.Write(content)
	})
}

func (s *Server) handleGetFiles() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		files, err := s.fileService.GetFiles(context.Background(), fileserve.GetFilesParam{})
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error Retrieving the Files", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(files)
	})
}
