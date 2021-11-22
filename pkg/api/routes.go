package api

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	model "github.com/adedayo/git-service-driver/pkg"
	gitutils "github.com/adedayo/git-service-driver/pkg/git"
	"github.com/adedayo/git-service-driver/pkg/gitlab"
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

	configManager = model.MakeConfigManager()
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
			Path:    "/api/github/clone",
			Handler: cloneFromService(GitHub),
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/gitlab/clone",
			Handler: cloneFromService(GitLab),
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/gitlab/discover",
			Handler: discoverGitLab,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/gitlab/integrate",
			Handler: integrateGitLab,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/gitlab/integrations",
			Handler: getGitLabIntegrations,
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

	results := []gitlab.GitLabProject{}
	config := configManager.GetConfig()

	for _, v := range config.GitServices[model.GitLab] {
		proj, err := gitlab.GetRepositories(r.Context(), model.GitService{
			GraphQLEndPoint: v.GraphQLEndPoint,
			API_Key:         v.API_Key,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, proj...)
	}

	json.NewEncoder(w).Encode(results)
}

func integrateGitLab(w http.ResponseWriter, r *http.Request) {
	var detail model.GitService
	if err := json.NewDecoder(r.Body).Decode(&detail); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	detail.GraphQLEndPoint = detail.GraphQLEndPoint + "/api/graphql"
	detail.ID = fmt.Sprintf("%x", sha1.New().Sum([]byte(detail.GraphQLEndPoint)))
	detail.Type = model.GitLab
	config := configManager.GetConfig()
	config.AddService(&detail)
	json.NewEncoder(w).Encode(listIntegrations(model.GitLab))
}

func getGitLabIntegrations(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(listIntegrations(model.GitLab))
}

func listIntegrations(sType model.GitServiceType) []model.GitService {
	config := configManager.GetConfig()
	services := config.GitServices[sType]
	out := []model.GitService{}
	for _, v := range services {
		out = append(out, model.GitService{
			GraphQLEndPoint: v.GraphQLEndPoint,
			ID:              v.ID,
		})
	}

	return out
}
