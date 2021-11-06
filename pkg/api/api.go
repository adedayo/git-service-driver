package api

import (
	"fmt"
	"log"
	"net/http"

	gitutils "github.com/adedayo/git-service-driver/pkg/git"
	"github.com/gorilla/handlers"
)

var (
	gitHubAuth *gitutils.GitAuth
	gitLabAuth *gitutils.GitAuth
)

//ServeAPI serves the analysis service on the specified port
func ServeAPI(config Config) {
	hostPort := "localhost:%d"
	if !config.Local {
		// not localhost electron app
		hostPort = ":%d"
	}
	gitHubAuth = &config.GitHubAuth
	gitLabAuth = &config.GitLabAuth
	corsOptions = append(corsOptions, handlers.AllowedOrigins(allowedOrigins))
	log.Printf("Running Git Service API on %s", fmt.Sprintf(hostPort, config.ApiPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(hostPort, config.ApiPort), handlers.CORS(corsOptions...)(routes)))
}
