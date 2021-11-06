package api

import gitutils "github.com/adedayo/git-service-driver/pkg/git"

type Config struct {
	ApiPort    int
	GitLabAuth gitutils.GitAuth
	GitHubAuth gitutils.GitAuth
	Local      bool //if set, to bind the api to localhost:port (electron) or simply :port (web service) instead
}
