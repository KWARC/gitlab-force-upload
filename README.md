# gitlab-force-upload

A Golang script to force upload a folder into a GitLab repository. 


## Usage

```
Usage of gitlab-force-upload:
  -dest string
        Destination repository (required)
  -folder string
        Folder to upload to GitLab (required)
  -legal
        Show legal information and exit
  -token string
        Token for GitLab (required)
  -url string
        GitLab URL to connect to (default "https://gitlab.com")
  -v    Log more verbose
```

Concretly, it performs the following actions:

1. Create a remote gitlab repository, if it does not yet exist
2. Unprotect the main branch
3. Create a new local repository, deleting any older one if it exists
4. Create a single new commit in the repository
5. Force-push this commit to the remote. 

## getting gitlab-force-upload

To get `gitlab-force-upload` you have two options:

- __Build it yourself__. To build `gitlab-force-upload` yourself, you need go 1.9 or newer along with make installed on your machine. After cloning this repository, you can then simply type make and executables will be generated inside the out/ directory.

- __Download a pre-built binary__. You can download a pre-built binary from the [releases page on GitHub](https://github.com/KWARC/gitlab-force-upload/releases/latest/). This page includes releases for Linux, Mac OS X and Windows.

After obtaining the binary (through either of the two means), simply place it in your $PATH. 
It does not depend on any external software (no need for git even).

## License

Released into the public domain, concretely licensed under the terms of the [Unlicense](http://unlicense.org). 