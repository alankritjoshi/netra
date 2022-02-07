package main

import (
	"log"

	"github.com/alankritjoshi/netra/api/server"
	"github.com/alankritjoshi/netra/internal/handler"
	"github.com/alankritjoshi/netra/internal/storage"
	"github.com/upper/db/v4/adapter/cockroachdb"
)

// The settings variable stores connection details.
var settings = cockroachdb.ConnectionURL{
	Host:     "localhost",
	Database: "netra",
	User:     "netra",
	Options: map[string]string{
		// Insecure node.
		"sslmode": "disable",
	},
}

func main() {
	sess, err := cockroachdb.Open(settings)
	if err != nil {
		log.Fatal("cockroachdb.Open: ", err)
	}
	defer sess.Close()

	issuesStore := storage.GetIssuesStore(sess)

	issueHandler := handler.NewIssueHandler(issuesStore)

	issueServer := server.NewIssueServer(&server.Config{}, issueHandler)

	if err := issueServer.ListenAndServe(); err != nil {
		log.Fatal("Issue Server: ", err)
	}

}
