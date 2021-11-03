package gitlab

import "github.com/hasura/go-graphql-client"

type projectQuery struct {
	Projects GitLabProjects `graphql:"projects(first: 5, after: $endCursor)"`
}

type GitLabProjects struct {
	Nodes []GitLabProject

	PageInfo struct {
		EndCursor   graphql.String
		HasNextPage graphql.Boolean
	}
}

type GitLabProject struct {
	Name              graphql.String
	NameWithNamespace graphql.String
	Description       graphql.String
	ID                graphql.ID
	Archived          graphql.Boolean
	SshUrlToRepo      graphql.String
	HttpUrlToRepo     graphql.String
	WebUrl            graphql.String
	Statistics        struct {
		RepositorySize graphql.Float
		StorageSize    graphql.Float
	}
	Repository struct {
		RootRef     graphql.String
		BranchNames []graphql.String `graphql:"branchNames(offset:0, searchPattern: \"*\", limit:1000)"`
	}
	ProjectMembers struct {
		Nodes []struct {
			ID   graphql.String
			User struct {
				ID     graphql.String
				Groups struct {
					Nodes []struct {
						ID   graphql.String
						Name graphql.String
					}
				}
			}
		}
	}
	Group struct {
		ID             graphql.ID
		Name           graphql.String
		Description    graphql.String
		FullName       graphql.String
		EmailsDisabled graphql.Boolean
		Contacts       struct {
			Nodes []struct {
				Email     graphql.String
				FirstName graphql.String
				LastName  graphql.String
			}
		}
		GroupMembers struct {
			Nodes []struct {
				User struct {
					ID               graphql.String
					Username         graphql.String
					GroupMemberships struct {
						Nodes []struct {
							ID    graphql.String
							Group struct {
								ID   graphql.String
								Name graphql.String
							}
						}
					}
				}
			}
		}
		Projects struct {
			Nodes []struct {
				ID           graphql.String
				Name         graphql.String
				SshUrlToRepo graphql.String
			}
		}
	}
}
