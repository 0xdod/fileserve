package http

import (
	"fmt"
	"net/http"

	"github.com/0xdod/fileserve"
	"github.com/0xdod/fileserve/filestorage"
	"github.com/0xdod/fileserve/sqlite"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Server struct {
	server    *http.Server
	db        *sqlite.DB
	mux       *mux.Router
	filestore filestorage.FileStorage
	fs        fileserve.FileService
}

type NewServerOpts struct {
	DB   *sqlite.DB
	Addr *string
}

func NewServer(opt NewServerOpts) *Server {
	addr := ":8080"

	if opt.Addr != nil {
		addr = *opt.Addr
	}

	s := &Server{
		db:  opt.DB,
		mux: mux.NewRouter(),
		filestore: filestorage.NewS3StorageBackend(&filestorage.S3BackendConfig{
			AccessKeyID:     viper.GetString("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: viper.GetString("AWS_SECRET_ACCESS_KEY"),
			Region:          viper.GetString("AWS_REGION"),
			BucketName:      viper.GetString("AWS_BUCKET_NAME"),
		}),
		fs: sqlite.NewFileService(opt.DB),
	}

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	s.registerRoutes()

	return s
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) registerRoutes() {
	s.mux.Use(loggerMiddleware)
	v1 := s.mux.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "<h1>Hello World</h2>")
	})
	v1.Handle("/files/upload", s.handleUpload()).Methods("POST")
	v1.Handle("/files/download/{fileId}", s.handleDownload()).Methods("GET")
	v1.Handle("/files", s.handleGetFiles()).Methods("GET")
}
