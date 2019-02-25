package src

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// MakeLocalRepo makes a local repository in the given folder
func MakeLocalRepo(folder string, commitMessage string, user *gitlab.User, verbose bool) (repo *git.Repository, err error) {

	// delete the previous repo
	if err = deleteRepo(folder, verbose); err != nil {
		return
	}

	return createRepository(folder, commitMessage, user, verbose)
}

// deleteRepo deletes the repository in folder, if any
func deleteRepo(folder string, verbose bool) (err error) {
	if _, e := git.PlainOpen(folder); e == nil {
		if verbose {
			fmt.Printf("  Deleting existing repository in %s\n", folder)
		}
		err = os.RemoveAll(filepath.Join(folder, ".git"))
	}
	return
}

// createRepository create a repository on disk
func createRepository(folder string, commitMessage string, user *gitlab.User, verbose bool) (repo *git.Repository, err error) {

	if verbose {
		fmt.Printf("  Creating new repository in %s\n", folder)
	}

	// create the new repository
	repo, err = git.PlainInit(folder, false)
	if err != nil {
		return
	}

	// get the worktree
	w, err := repo.Worktree()
	if err != nil {
		return
	}

	if verbose {
		fmt.Printf("  Adding all files in %s\n", folder)
	}

	status, err := w.Status()

	// run over the status
	for path, file := range status {
		if file.Worktree == git.Untracked {
			if _, err := w.Add(path); err != nil {
				return nil, err
			}
			if verbose {
				fmt.Printf("    Adding %q ...\n", path)
			}
		}
	}

	// and commit
	if verbose {
		fmt.Printf("  Creating new commit %q from %s<%s> \n", commitMessage, user.Name, user.Email)
	}

	c, err := w.Commit(commitMessage, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  user.Name,
			Email: user.Email,
			When:  time.Now(),
		},
	})

	if err == nil && verbose {
		fmt.Printf("  Created commit %s\n", c.String())
	}

	return
}
