package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	gitutils "github.com/adedayo/checkmate-core/pkg/git"

	"github.com/adedayo/checkmate-core/pkg/util"
	"github.com/adedayo/git-service-driver/pkg/github"
)

func integrateGitHub(w http.ResponseWriter, r *http.Request) {
	var detail gitutils.GitService
	if err := json.NewDecoder(r.Body).Decode(&detail); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	detail.InstanceURL = strings.TrimSuffix(strings.TrimSpace(detail.InstanceURL), "/")
	detail.GraphQLEndPoint = "https://api.github.com/graphql"
	detail.APIEndPoint = "https://api.github.com"
	detail.ID = util.NewRandomUUID().String()
	detail.Type = gitutils.GitHub
	// fmt.Printf("Got Integration: %#v\n", detail)
	config, err := configManager.GetConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	config.AddService(&detail)
	json.NewEncoder(w).Encode(listIntegrations(gitutils.GitHub))
}

func getGitHubIntegrations(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(listIntegrations(gitutils.GitHub))
}

func deleteGitHubIntegration(w http.ResponseWriter, r *http.Request) {
	var id struct {
		ID string
	}
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := configManager.GetConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	delete(config.GitServices[gitutils.GitHub], id.ID)
	configManager.SaveConfig(config)
	json.NewEncoder(w).Encode(listIntegrations(gitutils.GitHub))
}

func discoverGitHub(w http.ResponseWriter, r *http.Request) {

	config, err := configManager.GetConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var pagedSearch github.GitHubPagedSearch

	if err := json.NewDecoder(r.Body).Decode(&pagedSearch); err != nil {
		log.Printf("Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gitService, err := config.FindService(pagedSearch.ServiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	proj, loc, err := github.GetRepositories(r.Context(), gitService, &pagedSearch)

	if err != nil {
		log.Printf("Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(github.GitHubProjectSearchResult{
		InstanceID:             gitService.ID,
		Projects:               proj,
		EndCursor:              loc.EndCursor,
		HasNextPage:            loc.HasNextPage,
		RemainingProjectsCount: loc.TotalCount,
	})
}
