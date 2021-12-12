package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "github.com/adedayo/git-service-driver/pkg"
	"golang.org/x/net/context/ctxhttp"
)

func queryProjects(ctx context.Context, query string, gitService *model.GitService) (projects projectsQueryResultGH, err error) {
	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: map[string]interface{}{},
	}
	var buff bytes.Buffer
	err = json.NewEncoder(&buff).Encode(in)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", gitService.GraphQLEndPoint, &buff)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", gitService.API_Key))
	req.Header.Set("Content-Type", "application/json")
	resp, err := ctxhttp.Do(ctx, http.DefaultClient, req)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer resp.Body.Close()

	var out struct {
		Data struct {
			User struct {
				Repositories projectsQueryResultGH
			}
			Organization struct {
				Repositories projectsQueryResultGH
			}
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&out)
	if err == nil {
		if len(out.Data.Organization.Repositories.Nodes) == 0 {
			projects = out.Data.User.Repositories
		} else {
			projects = out.Data.Organization.Repositories

		}
	}
	return
}
