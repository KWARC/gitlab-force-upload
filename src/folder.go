package src

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	gitlab "github.com/xanzy/go-gitlab"
	git "gopkg.in/src-d/go-git.v4"
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

	// git add the files
	err = gitAdd(folder, w, verbose)

	// if something happened, return
	if err != nil {
		return
	}

	// and commit
	if verbose {
		fmt.Printf("  Creating new commit %q from %s<%s> \n", commitMessage, user.Name, user.Email)
	}

	// git commit the files
	err = gitCommit(folder, w, commitMessage, user.Name, user.Email, verbose)

	return
}

func gitAdd(folder string, w *git.Worktree, verbose bool) (err error) {
	err = gitAddExternal(folder, verbose)
	if err == nil {
		return
	}

	if verbose {
		fmt.Println("  Falling back to internal implementation", folder)
	}

	return gitAddInternal(folder, w, verbose)
}

func gitAddExternal(folder string, verbose bool) (err error) {
	if verbose {
		fmt.Printf("    Attempting to run `git add -A .` in %s\n", folder)
	}

	// git add -Av .
	cmd := exec.Command("git", "add", "-A", ".")
	cmd.Dir = folder

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdin
	}

	return cmd.Run()
}

func gitAddInternal(folder string, w *git.Worktree, verbose bool) (err error) {
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) (e error) {
		// if an error occured, return
		if err != nil {
			return err
		}

		// skip directories
		if info.IsDir() {
			return
		}

		// make it a relative path
		path, e = filepath.Rel(folder, path)
		if e != nil {
			return
		}

		// if we are starting with '.git/', then we should not add the path
		if strings.HasPrefix(path, ".git"+string(filepath.Separator)) {
			return
		}

		// print what we are adding
		if verbose {
			fmt.Printf("    Adding %q ...\n", path)
		}

		// add the path
		_, e = w.Add(path)

		return
	})
}

func gitCommit(folder string, w *git.Worktree, message string, name string, email string, verbose bool) (err error) {
	err = gitCommitExternal(folder, message, name, email, verbose)
	if err == nil {
		return
	}

	if verbose {
		fmt.Println("  Falling back to internal implementation")
	}

	return gitCommitInternal(folder, w, message, name, email, verbose)
}

func gitCommitExternal(folder string, message string, name string, email string, verbose bool) (err error) {
	authorString := fmt.Sprintf("%s <%s>", name, email)

	if verbose {
		fmt.Printf("    Attempting to run `git commit -m %q --author %q` in %s\n", message, authorString, folder)
	}

	// git add -Av .
	cmd := exec.Command("git", "commit", "-m", "'"+message+"'", "--author", "'"+authorString+"'")
	cmd.Dir = folder

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdin
	}

	return cmd.Run()
}

func gitCommitInternal(folder string, w *git.Worktree, message string, name string, email string, verbose bool) (err error) {
	c, err := w.Commit(message, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})

	if verbose {
		fmt.Printf("  Created commit %s\n", c)
	}

	return
}
