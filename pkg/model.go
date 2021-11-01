package gitlab

var (
	GITLAB_API_KEY          = "gitlab_api_key"
	GITLAB_GRAHPQL_ENDPOINT = "gitlab_graphql_endpoint"
)

type GitlabService struct {
	GraphQLEndPoint, API_Key string
}
