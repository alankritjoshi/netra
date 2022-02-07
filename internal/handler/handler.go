package handler

import (
	"log"
	"net/http"

	"github.com/alankritjoshi/netra/internal/storage"
	"github.com/go-chi/render"

	api "github.com/alankritjoshi/netra/api/v1"
)

func NewIssueHandler(issuesStore *storage.IssuesStore) *IssueHandler {
	return &IssueHandler{
		issuesStore: issuesStore,
	}
}

type IssueHandler struct {
	issuesStore *storage.IssuesStore
}

func (handler *IssueHandler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	data := &api.CreateIssueRequest{}
	log.Print("reached handler1")
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	log.Print("reached handler")
	issue := storage.IssueModel{
		Title:       data.Issue.Title,
		Description: data.Issue.Description,
	}
	log.Printf("reached handler %s", issue.Title)
	id, err := handler.issuesStore.CreateIssue(issue)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	log.Printf("Id %s", id)
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &api.CreateIssueResponse{Issue: data.Issue})
}

func (handler *IssueHandler) GetIssues(w http.ResponseWriter, r *http.Request) {
	log.Print("Inside handler!")
	list := []render.Renderer{}
	issuesModels, _ := handler.issuesStore.GetIssues()
	for i := range issuesModels {
		list = append(list, &api.CreateIssueResponse{
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
