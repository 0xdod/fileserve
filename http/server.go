package http

import (
	"fmt"
	"net/http"

	"github.com/0xdod/fileserve"
	"github.com/0xdod/fileserve/sqlite"
	"github.com/gorilla/mux"
)

type Server struct {
	server      *http.Server
	db          *sqlite.DB
	mux         *mux.Router
	fileStorage fileserve.FileStorage
	fileService fileserve.FileService
}

type NewServerOpts struct {
	DB          *sqlite.DB
	Addr        *string
	FileStorage fileserve.FileStorage
	FileService fileserve.FileService
}

func NewServer(opt NewServerOpts) *Server {
	addr := ":8080"

	if opt.Addr != nil {
		addr = *opt.Addr
	}

	s := &Server{
		db:          opt.DB,
		mux:         mux.NewRouter(),
		fileStorage: opt.FileStorage,
		fileService: opt.FileService,
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
	s.mux.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "<h1>Hello World</h2>")
	})
	v1 := s.mux.PathPrefix("/api/v1").Subrouter()
	v1.Handle("/files/upload", s.handleUpload()).Methods("POST")
	v1.Handle("/files/download/{fileId}", s.handleDownload()).Methods("GET")
	v1.Handle("/files", s.handleGetFiles()).Methods("GET")
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(nil)
}
