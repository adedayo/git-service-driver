package gitlab

import (
	"context"
	"fmt"

	model "github.com/adedayo/git-service-driver/pkg"
	"github.com/hasura/go-graphql-client"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func GetRepositories(ctx context.Context, gLab model.GitService, pagedSearch GitLabPagedSearch) (projects []GitLabProject, loc GitLabCursorLocation, err error) {

	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString(GITLAB_API_KEY)}))
	client := graphql.NewClient(gLab.GraphQLEndPoint, httpClient)
	var query projectQuery

	variables := map[string]interface{}{
		"endCursor": (*graphql.String)(nil),
	}

	if pagedSearch.NextCursor != "" {
		variables["endCursor"] = graphql.String(pagedSearch.NextCursor)
	}

	for len(projects) < pagedSearch.PageSize {
		err = client.Query(ctx, &query, variables)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		} else {
			// log.Printf("Projects %v\nCursor: %s\n", query.Projects.Nodes, graphql.String(query.Projects.PageInfo.EndCursor))
			loc.EndCursor = graphql.String(query.Projects.PageInfo.EndCursor)
			loc.HasNextPage = query.Projects.PageInfo.HasNextPage
			projects = append(projects, query.Projects.Nodes...)
			if !loc.HasNextPage {
				break
			}
			variables["endCursor"] = loc.EndCursor
		}
	}

	return
}

type GitLabPagedSearch struct {
	ServiceID  string
	PageSize   int
	NextCursor string
}

type GitLabCursorLocation struct {
	EndCursor   graphql.String
	HasNextPage graphql.Boolean
}
