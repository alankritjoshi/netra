package server

import (
	"net/http"
	"time"

	"github.com/alankritjoshi/netra/internal/handler"
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

func NewIssueServer(cfg *Config, handler *handler.IssueHandler) *IssueServer {
	return &IssueServer{
		cfg:     cfg,
		handler: handler,
	}
}

type IssueServer struct {
	cfg     *Config
	handler *handler.IssueHandler
}

func (server *IssueServer) ListenAndServe() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	r.Route("/issues", func(r chi.Router) {
		r.Get("/", server.handler.GetIssues)
		r.Post("/", server.handler.CreateIssue)
	})

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		return errors.Wrap(err, "ListenAndServe failed.")
	}
	return nil
}
