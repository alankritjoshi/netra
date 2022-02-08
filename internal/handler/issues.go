package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	api "github.com/alankritjoshi/netra/api/v1"
	"github.com/alankritjoshi/netra/internal/storage"
)

func NewIssuesHandler(issuesStore storage.Store) *IssuesHandler {
	return &IssuesHandler{
		issuesStore: issuesStore,
	}
}

type IssuesHandler struct {
	issuesStore storage.Store
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
		r.Get("/", h.GetAll)
		r.Post("/", h.Create)
		r.Get("/search", h.Search)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(h.IssueCtx)
			r.Get("/", h.GetOne)
			r.Delete("/", h.Delete)
		})
	})

	return r
}

func (h *IssuesHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &api.IssueRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	issue := storage.IssueModel{
		Title:       data.Issue.Title,
		Description: data.Issue.Description,
		Priority:    data.Issue.Priority,
	}
	id, err := h.issuesStore.Create(issue)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	data.Issue.ID = id
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &api.IssueResponse{Issue: data.Issue})
}

func (h *IssuesHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	issueModel := r.Context().Value("issue").(*storage.IssueModel)
	issue := &api.Issue{
		ID:          issueModel.ID,
		Title:       issueModel.Title,
		Description: issueModel.Description,
		Priority:    issueModel.Priority,
	}
	if err := render.Render(w, r, &api.IssueResponse{Issue: issue}); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func (h *IssuesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	issueModel := r.Context().Value("issue").(*storage.IssueModel)
	err := h.issuesStore.Delete(issueModel)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}
	issue := &api.Issue{
		ID:          issueModel.ID,
		Title:       issueModel.Title,
		Description: issueModel.Description,
		Priority:    issueModel.Priority,
	}
	render.Render(w, r, &api.IssueResponse{Issue: issue})
}

func (h *IssuesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list := []render.Renderer{}
	issuesModels, err := h.issuesStore.GetAll()
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	for i := range issuesModels {
		list = append(list, &api.IssueResponse{
			Issue: &api.Issue{
				ID:          issuesModels[i].ID,
				Title:       issuesModels[i].Title,
				Description: issuesModels[i].Description,
				Priority:    issuesModels[i].Priority,
			}})
	}
	if err := render.RenderList(w, r, list); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
}

func (h *IssuesHandler) Search(w http.ResponseWriter, r *http.Request) {
	list := []render.Renderer{}
	titleKey := r.URL.Query().Get("title")
	descKey := r.URL.Query().Get("description")
	var priorityLow, priorityHigh int
	var err error
	priorityLowStr := r.URL.Query().Get("priority_low")
	if len(priorityLowStr) > 0 {
		priorityLow, err = strconv.Atoi(priorityLowStr)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
	} else {
		priorityLow = -1
	}
	priorityHighStr := r.URL.Query().Get("priority_high")
	if len(priorityHighStr) > 0 {
		priorityHigh, err = strconv.Atoi(priorityHighStr)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
	} else {
		priorityHigh = -1
	}
	issuesModels, _ := h.issuesStore.Search(titleKey, descKey, priorityLow, priorityHigh)
	for i := range issuesModels {
		list = append(list, &api.IssueResponse{
			Issue: &api.Issue{
				ID:          issuesModels[i].ID,
				Title:       issuesModels[i].Title,
				Description: issuesModels[i].Description,
				Priority:    issuesModels[i].Priority,
			}})
	}
	if err := render.RenderList(w, r, list); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
}

func (h *IssuesHandler) IssueCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var issue *storage.IssueModel
		var err error

		if id := chi.URLParam(r, "id"); id != "" {
			issue, err = h.issuesStore.GetByID(id)
		} else {
			render.Render(w, r, ErrNotFound(err))
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound(err))
			return
		}

		ctx := context.WithValue(r.Context(), "issue", issue)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error.",
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

func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 404,
		StatusText:     "Resource not found.",
		ErrorText:      err.Error(),
	}
}
