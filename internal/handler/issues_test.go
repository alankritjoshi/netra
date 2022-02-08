package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/alankritjoshi/netra/internal/handler"
	"github.com/alankritjoshi/netra/internal/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/upper/db/v4"
)

type MockIssuesCollection struct {
	db.Collection
	mock.Mock
}

func (m *MockIssuesCollection) Create(issue storage.IssueModel) (string, error) {
	args := m.Called(issue)
	return args.String(0), args.Error(1)
}

func (m *MockIssuesCollection) GetByID(id string) (*storage.IssueModel, error) {
	args := m.Called(id)
	return args.Get(0).(*storage.IssueModel), args.Error(1)
}

func (m *MockIssuesCollection) Delete(issue *storage.IssueModel) error {
	args := m.Called(issue)
	return args.Error(0)
}

func (m *MockIssuesCollection) GetAll() ([]storage.IssueModel, error) {
	args := m.Called()
	return args.Get(0).([]storage.IssueModel), args.Error(1)
}

func (m *MockIssuesCollection) Search(titleKey, descKey string, priorityLow, priorityHigh int) ([]*storage.IssueModel, error) {
	args := m.Called()
	return args.Get(0).([]*storage.IssueModel), args.Error(1)
}

func TestCreateIssueSuccess(t *testing.T) {
	reader := strings.NewReader(`{"title": "hi", "description": "lol"}`)
	req := httptest.NewRequest(http.MethodPost, "/issues", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockIssuesStore := new(MockIssuesCollection)
	mockIssuesStore.On("Create", mock.Anything).Return("123", nil)
	issuesHandler := handler.NewIssuesHandler(mockIssuesStore)
	issuesHandler.Create(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateIssueBadInput(t *testing.T) {
	reader := strings.NewReader(`{"title": "hi", "description": "lol"`)
	req := httptest.NewRequest(http.MethodPost, "/issues", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockIssuesStore := new(MockIssuesCollection)
	mockIssuesStore.On("Create", mock.Anything).Return("123", nil)
	issuesHandler := handler.NewIssuesHandler(mockIssuesStore)
	issuesHandler.Create(w, req)
	require.Equal(t, w.Code, http.StatusBadRequest)
}

func TestCreateIssueMissingInput(t *testing.T) {
	reader := strings.NewReader(`{"description": "lol"}`)
	req := httptest.NewRequest(http.MethodPost, "/issues", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockIssuesStore := new(MockIssuesCollection)
	mockIssuesStore.On("Create", mock.Anything).Return("", errors.New("Insertion of issue failed"))
	issuesHandler := handler.NewIssuesHandler(mockIssuesStore)
	issuesHandler.Create(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetOneIssueSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/issues/1", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockIssuesStore := new(MockIssuesCollection)
	issueModel := &storage.IssueModel{
		ID:          "1",
		Title:       fmt.Sprintf("Title %d", 1),
		Description: fmt.Sprintf("Description %d", 1),
	}
	issuesHandler := handler.NewIssuesHandler(mockIssuesStore)
	issuesHandler.GetOne(w, req.WithContext(context.WithValue(req.Context(), "issue", issueModel)))
	require.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteIssueSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/issues/1", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockIssuesStore := new(MockIssuesCollection)
	issueModel := &storage.IssueModel{
		ID:          "1",
		Title:       fmt.Sprintf("Title %d", 1),
		Description: fmt.Sprintf("Description %d", 1),
	}
	issuesHandler := handler.NewIssuesHandler(mockIssuesStore)
	mockIssuesStore.On("Delete", mock.Anything).Return(nil)
	issuesHandler.Delete(w, req.WithContext(context.WithValue(req.Context(), "issue", issueModel)))
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetIssuesSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/issues", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockIssuesStore := new(MockIssuesCollection)
	issuesModels := make([]storage.IssueModel, 3)
	for i := 0; i < 3; i++ {
		issuesModels[i] = storage.IssueModel{
			ID:          strconv.FormatInt(int64(i), 10),
			Title:       fmt.Sprintf("Title %d", i),
			Description: fmt.Sprintf("Description %d", i),
			Priority:    uint(i),
		}
	}
	mockIssuesStore.On("GetAll").Return(issuesModels, nil)
	issuesHandler := handler.NewIssuesHandler(mockIssuesStore)
	issuesHandler.GetAll(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetIssuesStoreError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/issues", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockIssuesStore := new(MockIssuesCollection)
	mockIssuesStore.On("GetAll").Return(make([]storage.IssueModel, 0), errors.New("Insertion of issue failed"))
	issuesHandler := handler.NewIssuesHandler(mockIssuesStore)
	issuesHandler.GetAll(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}
