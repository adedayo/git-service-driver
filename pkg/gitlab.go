package gitlab

import (
	"context"
	"fmt"

	"github.com/hasura/go-graphql-client"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func GetRepositories(ctx context.Context, gLab GitlabService) (projects []GitLabProject, err error) {

	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString(GITLAB_API_KEY)}))
	client := graphql.NewClient(gLab.GraphQLEndPoint, httpClient)
	var query projectQuery

	variables := map[string]interface{}{
		"endCursor": (*graphql.String)(nil),
	}

	for {
		err = client.Query(ctx, &query, variables)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		} else {
			projects = append(projects, query.Projects.Nodes...)
			if !query.Projects.PageInfo.HasNextPage {
				break
			}
			variables["endCursor"] = graphql.String(query.Projects.PageInfo.EndCursor)
		}
	}

	return
}
