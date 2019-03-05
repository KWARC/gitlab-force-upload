package src

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// PushToRemote pushes a local repository to the remote
func PushToRemote(folder string, local *git.Repository, remote string, user *gitlab.User, token string, verbose bool) (err error) {

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

	err = gitPush(folder, r, remote, user.Username, token, verbose)

	return
}

func gitPush(folder string, r *git.Remote, upstream string, username string, password string, verbose bool) (err error) {
	err = gitPushExternal(folder, upstream, username, password, verbose)
	if err == nil {
		return
	}

	if verbose {
		fmt.Println("  Falling back to internal implementation")
	}

	return gitPushInternal(r, username, password)
}

func gitPushExternal(folder string, upstream string, username string, password string, verbose bool) (err error) {

	// split uri into (protocol, rest)
	uSplit := strings.SplitN(upstream, "://", 2)

	// add username:password@
	remoteString := fmt.Sprintf("%s://%s:%s@%s", uSplit[0], username, password, uSplit[1])

	if verbose {
		fmt.Printf("    Attempting to run `git push --force %q master:master ` in %s\n", remoteString, folder)
	}

	// git add -Av .
	cmd := exec.Command("git", "push", "--force", remoteString)
	cmd.Dir = folder

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdin
	}

	return cmd.Run()
}

func gitPushInternal(r *git.Remote, username string, password string) (err error) {
	// and push with the upstream
	err = r.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/master:refs/heads/master"),
		},
		Auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
		Progress: os.Stdout,
	})

	if err == git.NoErrAlreadyUpToDate {
		fmt.Println(err)
		err = nil
	}

	return
}
