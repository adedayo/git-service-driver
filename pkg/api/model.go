package api

import (
	"net/http"

	gitutils "github.com/adedayo/checkmate-core/pkg/git"
)

type Config struct {
	ApiPort     int
	GitLabAuth  gitutils.GitAuth
	GitHubAuth  gitutils.GitAuth
	Local       bool   //if set, to bind the api to localhost:port (electron) or simply :port (web service) instead
	CodeBaseDir string // the location where code is cloned into
}

type RouteSpec struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
	Methods []string
}
