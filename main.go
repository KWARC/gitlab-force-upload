package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("to be implemented")
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

	// TODO: Check that folder exists
	if folder == "" {
		fmt.Println("Missing -folder argument")
		os.Exit(1)
	}

	if dest == "" {
		fmt.Println("Missing -dest argument")
		os.Exit(1)
	}
}
