package github

const (
	projectsQuery = `
	{
		%s(login: "%s") {
		  repositories(first: %d%s) {
			nodes {
			  name
			  id
			  isArchived
			  url
			  isDisabled
			}
			totalCount
			pageInfo {
			  endCursor
			  hasNextPage
			}
		  }
		}
	  }	  
	`
)

type GitHubProject struct {
	Name       string
	ID         string
	IsArchived bool
	Url        string
	IsDisabled bool
}

type GitHubProjectSearchResult struct {
	InstanceID             string
	Projects               []GitHubProject
	EndCursor              string
	HasNextPage            bool
	RemainingProjectsCount int64 //how many projects remain after this qeury cursor
}

type projectsQueryResultGH struct {
	Nodes      []GitHubProject
	TotalCount int64
	PageInfo   struct {
		EndCursor   string
		HasNextPage bool
	}
}
