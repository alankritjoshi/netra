package storage

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/upper/db/v4"
)

type Store interface {
	db.Collection
	Create(IssueModel) (string, error)
	GetByID(string) (*IssueModel, error)
	Delete(*IssueModel) error
	GetAll() ([]IssueModel, error)
	Search(string, string, int, int) ([]*IssueModel, error)
}

func GetIssuesStore(sess db.Session) Store {
	return &IssuesStore{sess.Collection("issues")}
}

type IssuesStore struct {
	db.Collection
}

func (issues *IssuesStore) Create(issue IssueModel) (string, error) {
	res, err := issues.Insert(issue)
	if err != nil {
		return "", errors.Wrap(err, "Insertion of issue failed")
	}
	return fmt.Sprintf("%#v", res.ID()), nil
}

func (issues *IssuesStore) GetByID(id string) (*IssueModel, error) {
	var issue IssueModel
	err := issues.Find(db.Cond{"id": id}).One(&issue)
	if err != nil {
		return nil, errors.Wrap(err, "Get issue by ID failed")
	}
	return &issue, nil
}

func (issues *IssuesStore) Delete(issue *IssueModel) error {
	err := issues.Session().Delete(issue)
	if err != nil {
		return errors.Wrap(err, "Deletion of issue by ID failed")
	}
	return nil
}

func (issues *IssuesStore) GetAll() ([]IssueModel, error) {
	issuesArr := []IssueModel{}
	err := issues.Find().All(&issuesArr)
	if err != nil {
		return nil, errors.Wrap(err, "Retrieval of issues failed")
	}
	return issuesArr, nil
}

func (issues *IssuesStore) Search(titleKey, descKey string, priorityLow, priorityHigh int) ([]*IssueModel, error) {
	issuesArr := []*IssueModel{}
	selector := issues.Session().SQL().SelectFrom("issues")
	if len(titleKey) > 0 {
		selector = selector.And("title LIKE ?", fmt.Sprintf("%%%s%%", titleKey))
	}
	if len(descKey) > 0 {
		selector = selector.And("description LIKE ?", fmt.Sprintf("%%%s%%", descKey))
	}
	if priorityLow > -1 {
		selector = selector.And("priority >=", priorityLow)
	}
	if priorityHigh > -1 {
		selector = selector.And("priority <=", priorityHigh)
	}
	err := selector.All(&issuesArr)
	if err != nil {
		return nil, errors.Wrap(err, "Search failed")
	}
	return issuesArr, nil
}

type IssueModel struct {
	ID          string `db:"id,omitempty"`
	Title       string `db:"title"`
	Description string `db:"description,omitempty"`
	Priority    uint   `db:"priority,omitempty"`
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
