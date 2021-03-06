// Original: https://github.com/go-git/go-git/blob/bc1f419cebcf7505db31149fa459e9e3f8260e00/_examples/clone/auth/ssh/main.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	git "github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"golang.org/x/crypto/ssh"
)

func main() {
	CheckArgs("<url>", "<directory>", "<private_key_file>")
	url, directory, privateKeyFile := os.Args[1], os.Args[2], os.Args[3]

	_, err := os.Stat(privateKeyFile)
	if err != nil {
		Warning("read file %s failed: %s\n", privateKeyFile, err.Error())
		return
	}
	privKey, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		Warning("load file %s failed: %s\n", privateKeyFile, err.Error())
		return
	}

	signer, err := ssh.ParsePrivateKey([]byte(privKey))
	if err != nil {
		Warning("parsing ssh private key %s failed: %s\n", privateKeyFile, err.Error())
		return
	}

	sshAuth := gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
	}

	// Clone the given repository to the given directory
	Info("git clone %s ", url)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		Auth:     &sshAuth,
		URL:      url,
		Progress: os.Stdout,
	})
	CheckIfError(err)

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}
