package netra_v1

import (
	"net/http"

	"github.com/pkg/errors"
)

type Issue struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateIssueRequest struct {
	*Issue
}

func (request *CreateIssueRequest) Bind(r *http.Request) error {
	if request.Issue == nil {
		return errors.New("Missing required Issue fields.")
	}
	return nil
}

type CreateIssueResponse struct {
	*Issue
}

func (response *CreateIssueResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetIssuesRequest struct {
}

func (request *GetIssuesRequest) Bind(r *http.Request) error {
	return nil
}

type GetIssuesResponse struct {
	Issues []*Issue
}

func (response *GetIssuesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
