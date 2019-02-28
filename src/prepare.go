package src

import (
	"fmt"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
)

// PrepareRepo makes sure that the repo exists and can be force pushed to
func PrepareRepo(tokenIn string, url string, proEdition bool, repo string, verbose bool) (uri string, user *gitlab.User, token string, err error) {
	// return the token again
	token = tokenIn

	// create a new client
	gl := gitlab.NewClient(nil, token)
	gl.SetBaseURL(url + "api/v4/")

	// get the current user
	user, _, err = gl.Users.CurrentUser()
	if err != nil {
		return
	}
	// get or create the repo
	pro, err := getOrCreate(gl, proEdition, repo, verbose)
	if err != nil {
		return
	}

	// unblock the main branch
	err = unprotectMainBranch(gl, pro, verbose)
	if err != nil {
		return
	}

	// and return the url to the repo
	uri = pro.HTTPURLToRepo
	return
}

// getOrCreate gets or create a new project
func getOrCreate(gl *gitlab.Client, proEdition bool, name string, verbose bool) (pro *gitlab.Project, err error) {
	pro, err = getProject(gl, name)
	if err == nil {
		if verbose {
			fmt.Printf("  Not creating %s: Already exists\n", pro.PathWithNamespace)
		}
		return
	}

	if verbose {
		fmt.Printf("  Creating Project %s. \n", name)
	}
	return createProject(gl, proEdition, name, verbose)
}

// getProject gets a single project
func getProject(gl *gitlab.Client, name string) (pro *gitlab.Project, err error) {
	pro, _, err = gl.Projects.GetProject(name, &gitlab.GetProjectOptions{})
	return
}

// createProject create a new project
func createProject(gl *gitlab.Client, proEdition bool, name string, verbose bool) (pro *gitlab.Project, err error) {
	// split into path and name and prepare options
	repoPath, repoName := splitPath(name)

	// get the namespace
	ns, _, err := gl.Namespaces.GetNamespace(repoPath)
	if err != nil {
		if verbose {
			fmt.Printf("  Cannot create Project %s: Namespace %s does not exist\n", name, repoPath)
		}
		return
	}

	// argument to create the gitlab project
	createArgs := &gitlab.CreateProjectOptions{
		Path:        &repoName,
		NamespaceID: &ns.ID,
	}

	// the ultimate edition requires the 'ApprovalsBeforeMerge' flag
	// or causes a database error that is exposed via the api
	// this is a workaround
	if proEdition {
		zero := 0
		createArgs.ApprovalsBeforeMerge = &zero
	}

	// and create the project in it
	pro, _, err = gl.Projects.CreateProject(createArgs)

	// done
	return
}

// unprotectMainBranch unprotects the main branch opf the given project
func unprotectMainBranch(gl *gitlab.Client, pro *gitlab.Project, verbose bool) (err error) {
	if pro.DefaultBranch != "" {
		if verbose {
			fmt.Printf("  Unprotecting main branch %s of %s\n", pro.DefaultBranch, pro.PathWithNamespace)
		}
		_, _, err = gl.Branches.UnprotectBranch(pro.ID, pro.DefaultBranch)
	} else {
		if verbose {
			fmt.Printf("  Not unprotecting main branch of %s: No main branch\n", pro.PathWithNamespace)
		}
	}
	return
}

// splitPath splits a repository name into path and name
func splitPath(uri string) (path string, name string) {
	idx := strings.LastIndex(uri, "/")
	if idx == -1 {
		return "", uri
	}

	path = uri[:idx]
	name = uri[idx+1:]

	return
}
