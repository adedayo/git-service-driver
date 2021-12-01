package gitlab

import (
	"context"
	"fmt"
	"log"

	model "github.com/adedayo/git-service-driver/pkg"
)

func GetRepositories(ctx context.Context, gLab *model.GitService, pagedSearch *GitLabPagedSearch) (projects []GitLabProject, loc GitLabCursorLocation, err error) {

	for len(projects) < pagedSearch.PageSize {
		query := fmt.Sprintf(projectsQuery, pagedSearch.First, pagedSearch.NextCursor)
		projs, err := QueryProjects(ctx, query, gLab)
		if err != nil {
			log.Printf("Error: %v\n", err)
			break
		} else {
			loc.EndCursor = projs.PageInfo.EndCursor
			loc.HasNextPage = projs.PageInfo.HasNextPage
			projects = append(projects, projs.Nodes...)
			if !loc.HasNextPage {
				break
			}
			pagedSearch.NextCursor = loc.EndCursor
		}
	}

	return
}

type GitLabPagedSearch struct {
	ServiceID  string
	PageSize   int
	First      int //(first: n, ...) in the query
	NextCursor string
}

type GitLabCursorLocation struct {
	EndCursor   string
	HasNextPage bool
}
