package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/KWARC/gitlab-force-upload/src"
)

func main() {
	fmt.Printf("gitlab-force-upload: use -legal to view licensing and legal information\n")

	if verbose {
		fmt.Printf("verbose: %b\n", verbose)
		fmt.Printf("gitlabURL: %q\n", gitlabURL)
		fmt.Printf("authToken: %q\n", authToken)
		fmt.Printf("folder: %q\n", folder)
		fmt.Printf("dest: %q\n", dest)
		fmt.Printf("-------\n")
	}

	// prepare the remote repository
	if verbose {
		fmt.Println("Preparing repository")
	}

	uri, user, token, err := src.PrepareRepo(authToken, gitlabURL, dest, verbose)
	if err != nil {
		panic(err)
	}

	if verbose {
		fmt.Printf("Created %s for user %s and token %s\n", uri, user.Username, token)
	}

	if folder == "" {
		if verbose {
			fmt.Println("No local folder provided, exiting. ")
		}
		return
	}

	// create the local repository
	if verbose {
		fmt.Println("Making local repository")
	}

	local, err := src.MakeLocalRepo(folder, "Initial commit", user, verbose)
	if err != nil {
		panic(err)
	}

	// and push it
	if verbose {
		fmt.Println("Pushing repository")
	}

	err = src.PushToRemote(local, uri, user, token, verbose)
	if err != nil {
		panic(err)
	}

	if verbose {
		fmt.Println("Done. ")
	}

}

var legal bool
var verbose bool
var gitlabURL string
var authToken string
var folder string
var dest string

func init() {
	flag.BoolVar(&verbose, "legal", false, "Show legal information and exit")
	flag.BoolVar(&verbose, "v", false, "Log more verbose")
	flag.StringVar(&gitlabURL, "url", "https://gitlab.com", "GitLab URL to connect to")
	flag.StringVar(&authToken, "token", "", "Token for GitLab (required)")
	flag.StringVar(&folder, "folder", "", "Folder to upload to GitLab (required)")
	flag.StringVar(&dest, "dest", "", "Destination repository (required)")
	flag.Parse()

	if legal {
		src.Legal()
		os.Exit(0)
	}

	if authToken == "" {
		fmt.Println("Missing -token argument")
		os.Exit(1)
	}

	if !strings.HasSuffix(gitlabURL, "/") {
		gitlabURL = gitlabURL + "/"
	}

	if dest == "" {
		fmt.Println("Missing -dest argument")
		os.Exit(1)
	}
}
