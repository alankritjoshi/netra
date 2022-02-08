package server

import (
	"net/http"
	"time"

	"github.com/alankritjoshi/netra/internal/handler"
	"github.com/alankritjoshi/netra/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Config struct {
	Logging bool
	Version string
}

type Server interface {
	ListenAndServe() error
}

func NewIssuesServer(cfg *Config, store *storage.IssuesStore) *IssuesServer {
	return &IssuesServer{
		cfg:   cfg,
		store: store,
	}
}

type IssuesServer struct {
	cfg   *Config
	store *storage.IssuesStore
}

func (server *IssuesServer) ListenAndServe() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	r.Mount("/issues", handler.NewIssuesHandler(server.store).Routes())

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		return errors.Wrap(err, "ListenAndServe failed.")
	}
	return nil
}
