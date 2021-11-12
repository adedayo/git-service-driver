package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	model "github.com/adedayo/git-service-driver/pkg"
	gitutils "github.com/adedayo/git-service-driver/pkg/git"
	"github.com/adedayo/git-service-driver/pkg/gitlab"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
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
	for _, rs := range GetRoutes() {
		routes.HandleFunc(rs.Path, rs.Handler).Methods(rs.Methods...)
	}
}

func GetRoutes() []RouteSpec {
	routeSpecs := []RouteSpec{
		{
			Path:    "/github/clone",
			Handler: cloneFromService(GitHub),
			Methods: []string{"POST"},
		},
		{
			Path:    "/gitlab/clone",
			Handler: cloneFromService(GitLab),
			Methods: []string{"POST"},
		},
		{
			Path:    "/gitlab/discover",
			Handler: discoverGitLab,
			Methods: []string{"GET"},
		},
	}
	return routeSpecs
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

func discoverGitLab(w http.ResponseWriter, r *http.Request) {
	proj, err := gitlab.GetRepositories(r.Context(), model.GitService{
		GraphQLEndPoint: viper.GetString(gitlab.GITLAB_GRAHPQL_ENDPOINT),
		API_Key:         viper.GetString(gitlab.GITLAB_API_KEY),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(proj)
}
