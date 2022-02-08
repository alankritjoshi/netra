package storage

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/upper/db/v4"
)

func GetIssuesStore(sess db.Session) *IssuesStore {
	return &IssuesStore{sess.Collection("issues")}
}

type IssuesStore struct {
	db.Collection
}

func (issues *IssuesStore) CreateIssue(issue IssueModel) (string, error) {
	res, err := issues.Insert(issue)
	if err != nil {
		return "", errors.Wrap(err, "Insertion of issue failed")
	}
	return fmt.Sprintf("%#v", res.ID()), nil
}

func (issues *IssuesStore) GetIssues() ([]IssueModel, error) {
	issuesArr := []IssueModel{}
	err := issues.Find().All(&issuesArr)
	if err != nil {
		return nil, errors.Wrap(err, "Insertion of issue failed")
	}
	return issuesArr, nil
}

type IssueModel struct {
	ID          string `db:"id,omitempty"`
	Title       string `db:"title"`
	Description string `db:"description,omitempty"`
}

func (issue *IssueModel) Store(sess db.Session) db.Store {
	return GetIssuesStore(sess)
}

func (issue *IssueModel) BeforeUpdate(sess db.Session) error {
	fmt.Println("**** BeforeUpdate was called ****")
	return nil
}

func (issue *IssueModel) AfterUpdate(sess db.Session) error {
	fmt.Println("**** AfterUpdate was called ****")
	return nil
}

var _ = interface {
	db.Record
	db.BeforeUpdateHook
	db.AfterUpdateHook
}(&IssueModel{})

var _ = interface {
	db.Store
}(&IssuesStore{})
