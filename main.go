package main

import (
	"log"

	"github.com/alankritjoshi/netra/api/server"
	"github.com/alankritjoshi/netra/internal/storage"
	"github.com/upper/db/v4/adapter/cockroachdb"
)

// The settings variable stores connection details.
var settings = cockroachdb.ConnectionURL{
	Host:     "cockroachdb",
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

	issueServer := server.NewIssuesServer(
		&server.Config{},
		storage.GetIssuesStore(sess),
	)

	log.Println("Server started...")

	if err := issueServer.ListenAndServe(); err != nil {
		log.Fatal("Issue Server: ", err)
	}

}
