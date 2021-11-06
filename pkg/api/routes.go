package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	gitutils "github.com/adedayo/git-service-driver/pkg/git"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type GitServiceName int

const (
	GitHub GitServiceName = iota
	GitLab
)

var (
	routes         = mux.NewRouter()
	allowedOrigins = []string{
		"localhost:17285",
		"http://localhost",
	}
	corsOptions = []handlers.CORSOption{
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept", "Accept-Language", "Origin"}),
		handlers.AllowCredentials(),
		handlers.AllowedOriginValidator(allowedOriginValidator),
	}
)

func init() {
	addRoutes()
}

func allowedOriginValidator(origin string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == origin {
			return true
		}
	}
	passCORS := strings.Split(origin, ":")[0] == "localhost" //allow localhost independent of port
	if !passCORS {
		log.Printf("Host %s fails CORS.", origin)
	}
	return passCORS
}

func addRoutes() {
	routes.HandleFunc("/github/clone", cloneFromService(GitHub)).Methods("POST")
	routes.HandleFunc("/gitlab/clone", cloneFromService(GitLab)).Methods("POST")
}

func cloneFromService(service GitServiceName) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var auth *gitutils.GitAuth
		switch service {
		case GitHub:
			auth = gitHubAuth
		case GitLab:
			auth = gitLabAuth
		default:
			auth = &gitutils.GitAuth{}
		}

		var spec gitutils.RepositoryCloneSpec
		if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		spec.Options.Auth = auth
		path, err := gitutils.Clone(r.Context(), spec.Repository, &spec.Options)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(path)
	}
}
