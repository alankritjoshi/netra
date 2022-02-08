package netra_v1

import (
	"net/http"

	"github.com/pkg/errors"
)

type Issue struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type IssueRequest struct {
	*Issue
}

func (request *IssueRequest) Bind(r *http.Request) error {
	if request.Issue == nil {
		return errors.New("Missing required Issue fields.")
	}
	return nil
}

type IssueResponse struct {
	*Issue
}

func (response *IssueResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
