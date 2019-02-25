package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/KWARC/gitlab-force-upload/src"
)

func main() {
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

	_, err = src.MakeLocalRepo(folder, "Initial commit", user, verbose)
	if err != nil {
		panic(err)
	}

	// and push it

}

var verbose bool
var gitlabURL string
var authToken string
var folder string
var dest string

func init() {
	flag.BoolVar(&verbose, "v", false, "Log more verbose")
	flag.StringVar(&gitlabURL, "url", "https://gitlab.com", "GitLab URL to connect to")
	flag.StringVar(&authToken, "token", "", "Token for GitLab (required)")
	flag.StringVar(&folder, "folder", "", "Folder to upload to GitLab (required)")
	flag.StringVar(&dest, "dest", "", "Destionation repository (required)")
	flag.Parse()

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
