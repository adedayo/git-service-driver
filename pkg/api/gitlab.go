package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/adedayo/checkmate-core/pkg/util"
	model "github.com/adedayo/git-service-driver/pkg"
	"github.com/adedayo/git-service-driver/pkg/gitlab"
)

func discoverGitLab(w http.ResponseWriter, r *http.Request) {

	config := configManager.GetConfig()
	var pagedSearch gitlab.GitLabPagedSearch

	if err := json.NewDecoder(r.Body).Decode(&pagedSearch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gitService, err := config.FindService(pagedSearch.ServiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	proj, loc, err := gitlab.GetRepositories(r.Context(), gitService, &pagedSearch)

	if err != nil {
		log.Printf("Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(gitlab.GitLabProjectSearchResult{
		InstanceID:             gitService.ID,
		Projects:               proj,
		EndCursor:              loc.EndCursor,
		HasNextPage:            loc.HasNextPage,
		RemainingProjectsCount: getCount(loc),
	})
}

func getCount(loc gitlab.GitLabCursorLocation) int64 {
	if loc.HasNextPage {
		if x, err := base64.RawStdEncoding.DecodeString(loc.EndCursor); err == nil {
			var out struct {
				ID string `json:"id"`
			}
			if e := json.Unmarshal(x, &out); e == nil {
				n, err := strconv.ParseInt(out.ID, 10, 0)
				if err != nil {
					log.Printf("%v", e)
					return 0
				}
				return n
			} else {
				log.Printf("%v", e)
			}
		} else {
			log.Printf("%v", err)
		}
	}
	return 0
}

func integrateGitLab(w http.ResponseWriter, r *http.Request) {
	log.Printf("Called the gitlab integrate")
	var detail model.GitService
	if err := json.NewDecoder(r.Body).Decode(&detail); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	detail.InstanceURL = strings.TrimSuffix(strings.TrimSpace(detail.InstanceURL), "/")
	detail.GraphQLEndPoint = detail.InstanceURL + "/api/graphql"
	detail.APIEndPoint = detail.InstanceURL + "/api"
	detail.ID = util.NewRandomUUID().String()
	detail.Type = model.GitLab
	log.Printf("Gitlab detail: %#v", detail)
	// fmt.Printf("Got Integration: %#v\n", detail)
	config := configManager.GetConfig()
	if err := config.AddService(&detail); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(listIntegrations(model.GitLab))
}

func deleteGitLabIntegration(w http.ResponseWriter, r *http.Request) {
	var id struct {
		ID string
	}
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config := configManager.GetConfig()
	delete(config.GitServices[model.GitLab], id.ID)
	configManager.SaveConfig(config)
	json.NewEncoder(w).Encode(listIntegrations(model.GitLab))
}

func getGitLabIntegrations(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(listIntegrations(model.GitLab))
}
