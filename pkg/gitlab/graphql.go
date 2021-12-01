package gitlab

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

func QueryProjects(ctx context.Context, query string, gitService *model.GitService) (projects projectsQueryResult, err error) {
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gitService.API_Key))
	req.Header.Set("Content-Type", "application/json")
	resp, err := ctxhttp.Do(context.Background(), http.DefaultClient, req)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer resp.Body.Close()
	var out struct {
		Data struct {
			Projects projectsQueryResult
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&out)
	if err == nil {
		projects = out.Data.Projects
	}
	return
}
