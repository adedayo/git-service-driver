package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	gitutils "github.com/adedayo/checkmate-core/pkg/git"
	model "github.com/adedayo/git-service-driver/pkg"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
			Handler: cloneFromService(model.GitHub),
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/gitlab/clone",
			Handler: cloneFromService(model.GitLab),
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/github/discover",
			Handler: discoverGitHub,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/gitlab/discover",
			Handler: discoverGitLab,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/github/integrate",
			Handler: integrateGitHub,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/gitlab/integrate",
			Handler: integrateGitLab,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/github/deleteintegration",
			Handler: deleteGitHubIntegration,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/gitlab/deleteintegration",
			Handler: deleteGitLabIntegration,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/github/integrations",
			Handler: getGitHubIntegrations,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/gitlab/integrations",
			Handler: getGitLabIntegrations,
			Methods: []string{"GET"},
		},
	}
	return routeSpecs
}

func cloneFromService(service model.GitServiceType) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var spec gitutils.RepositoryCloneSpec
		if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		service, err := configManager.GetConfig().GetService(service, spec.ServiceID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		spec.Options.Auth = service.MakeAuth()

		path, err := gitutils.Clone(r.Context(), spec.Repository, &spec.Options)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(path)
	}
}

func listIntegrations(sType model.GitServiceType) []model.GitService {
	config := configManager.GetConfig()
	services := config.GitServices[sType]
	out := []model.GitService{}
	for _, v := range services {
		out = append(out, model.GitService{
			GraphQLEndPoint: v.GraphQLEndPoint,
			ID:              v.ID,
			InstanceURL:     v.InstanceURL,
			Name:            v.Name,
			AccountName:     v.AccountName,
			AccountType:     v.AccountType,
		})
	}

	sort.Sort(gitServiceList(out))
	return out
}

type gitServiceList []model.GitService

func (gs gitServiceList) Len() int {
	return len(gs)
}

func (gs gitServiceList) Less(i, j int) bool {
	return gs[i].Name < gs[j].Name || (gs[i].Name == gs[j].Name && gs[i].InstanceURL < gs[j].InstanceURL)
}

func (gs gitServiceList) Swap(i, j int) {
	gs[i], gs[j] = gs[j], gs[i]
}
