package handler

import (
	"net/http"
	"time"

	"github.com/alankritjoshi/netra/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	api "github.com/alankritjoshi/netra/api/v1"
)

func NewIssuesHandler(issuesStore *storage.IssuesStore) *IssuesHandler {
	return &IssuesHandler{
		issuesStore: issuesStore,
	}
}

type IssuesHandler struct {
	issuesStore *storage.IssuesStore
}

func (h *IssuesHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	r.Route("/issues", func(r chi.Router) {
		r.Get("/", h.GetIssues)
		r.Post("/", h.CreateIssue)
	})

	return r
}

func (h *IssuesHandler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	data := &api.IssueRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	issue := storage.IssueModel{
		Title:       data.Issue.Title,
		Description: data.Issue.Description,
	}
	id, err := h.issuesStore.CreateIssue(issue)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	data.Issue.ID = id
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &api.IssueResponse{Issue: data.Issue})
}

func (h *IssuesHandler) GetIssues(w http.ResponseWriter, r *http.Request) {
	list := []render.Renderer{}
	issuesModels, _ := h.issuesStore.GetIssues()
	for i := range issuesModels {
		list = append(list, &api.IssueResponse{
			Issue: &api.Issue{
				ID:          issuesModels[i].ID,
				Title:       issuesModels[i].Title,
				Description: issuesModels[i].Description,
			}})
	}
	if err := render.RenderList(w, r, list); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
}

func (h *IssuesHandler) SearchIssue(w http.ResponseWriter, r *http.Request) {
	list := []render.Renderer{}
	issuesModels, _ := h.issuesStore.GetIssues()
	for i := range issuesModels {
		list = append(list, &api.IssueResponse{
			Issue: &api.Issue{
				ID:          issuesModels[i].ID,
				Title:       issuesModels[i].Title,
				Description: issuesModels[i].Description,
			}})
	}
	if err := render.RenderList(w, r, list); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
