package src

import (
	"fmt"
	"os"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// PushToRemote pushes a local repository to the remote
func PushToRemote(local *git.Repository, remote string, user *gitlab.User, token string, verbose bool) (err error) {

	// Create the remote
	if verbose {
		fmt.Printf("  Creating remote %q pointing to %q\n", "origin", remote)
	}
	r, err := local.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remote},
	})

	if err != nil {
		return
	}

	if verbose {
		fmt.Println("  Pushing repository to remote")
	}

	// and push with the upstream
	err = r.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/master:refs/heads/master"),
		},
		Auth: &http.BasicAuth{
			Username: user.Username,
			Password: token,
		},
		Progress: os.Stdout,
	})

	if err == git.NoErrAlreadyUpToDate {
		fmt.Println(err)
		err = nil
	}

	return

}
