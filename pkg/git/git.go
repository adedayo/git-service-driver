package gitutils

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

//Clone a repository,returning the location on disk where the clone is placed
func Clone(repository string, options GitCloneOptions) (string, error) {
	repository = strings.ToLower(repository)
	//git@ is not supported, replace with https://
	if strings.HasPrefix(repository, "git@") {
		repository = strings.Replace(strings.Replace(repository, ":", "/", 1), "git@", "https://", 1)
	}

	dir, err := filepath.Abs(path.Clean(path.Join(options.BaseDir, strings.TrimSuffix(path.Base(repository), ".git"))))

	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			os.RemoveAll(dir)
			log.Printf("Error: %v\n", err)
		}
	}()

	if err = os.Mkdir(dir, 0755); err != nil {
		return "", err
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: repository,
		// Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: options.Auth.User,
			Password: options.Auth.Credential,
		},
	})

	if err != nil {
		return "", err
	}

	if options.CommitHash != "" {
		w, err := repo.Worktree()

		if err != nil {
			return "", err
		}

		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(options.CommitHash),
		})

		if err != nil {
			return "", err
		}
	}

	return dir, nil
}
