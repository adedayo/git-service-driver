package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/adedayo/checkmate-core/pkg/util"
	model "github.com/adedayo/git-service-driver/pkg"
	"github.com/adedayo/git-service-driver/pkg/github"
)

func integrateGitHub(w http.ResponseWriter, r *http.Request) {
	var detail model.GitService
	if err := json.NewDecoder(r.Body).Decode(&detail); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	detail.InstanceURL = strings.TrimSuffix(strings.TrimSpace(detail.InstanceURL), "/")
	detail.GraphQLEndPoint = "https://api.github.com/graphql"
	detail.APIEndPoint = "https://api.github.com"
	detail.ID = util.NewRandomUUID().String()
	detail.Type = model.GitHub
	// fmt.Printf("Got Integration: %#v\n", detail)
	config := configManager.GetConfig()
	config.AddService(&detail)
	json.NewEncoder(w).Encode(listIntegrations(model.GitHub))
}

func getGitHubIntegrations(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(listIntegrations(model.GitHub))
}

func deleteGitHubIntegration(w http.ResponseWriter, r *http.Request) {
	var id struct {
		ID string
	}
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config := configManager.GetConfig()
	delete(config.GitServices[model.GitHub], id.ID)
	configManager.SaveConfig(config)
	json.NewEncoder(w).Encode(listIntegrations(model.GitHub))
}

func discoverGitHub(w http.ResponseWriter, r *http.Request) {

	config := configManager.GetConfig()
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
