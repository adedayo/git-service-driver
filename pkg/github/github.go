package github

import (
	"context"
	"fmt"
	"log"

	gitutils "github.com/adedayo/checkmate-core/pkg/git"
)

func GetRepositories(ctx context.Context, gHub *gitutils.GitService, pagedSearch *GitHubPagedSearch) (projects []GitHubProject, loc GitHubCursorLocation, err error) {

	if pagedSearch.First < 1 {
		pagedSearch.First = 7 //conservatively push up the number of projects retrieved if they forgot to set First parameter
	}

	accountType := "user"

	if gHub.AccountType == "Organization" {
		accountType = "organization"
	}

	for len(projects) < pagedSearch.PageSize {
		query := fmt.Sprintf(projectsQuery, accountType, gHub.AccountName, pagedSearch.First, formatCursor(pagedSearch.NextCursor))
		projs, err := queryProjects(ctx, query, gHub)
		if err != nil {
			log.Printf("Error: %v\n", err)
			break
		} else {
			loc.EndCursor = projs.PageInfo.EndCursor
			loc.HasNextPage = projs.PageInfo.HasNextPage
			loc.TotalCount = projs.TotalCount
			projects = append(projects, projs.Nodes...)
			if !loc.HasNextPage {
				break
			}
			pagedSearch.NextCursor = loc.EndCursor
		}
	}

	return
}

func formatCursor(cursor string) string {
	if cursor == "" {
		return ""
	}
	return fmt.Sprintf(`, after: "%s"`, cursor)
}

type GitHubPagedSearch struct {
	ServiceID  string
	PageSize   int
	First      int //(first: n, ...) in the query
	NextCursor string
}

type GitHubCursorLocation struct {
	EndCursor   string
	HasNextPage bool
	TotalCount  int64
}
