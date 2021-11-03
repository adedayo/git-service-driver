package gitutils

type GitAuth struct {
	User, Credential string
}

type GitCloneOptions struct {
	BaseDir    string // directory to clone into
	Auth       GitAuth
	CommitHash string //if set checkout the specified commit
}
