package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/fatih/color"
)

var (
	noTitleErr = errors.New("A title is required")
	noDescErr  = errors.New("A description is required")
	koColor    = color.New(color.FgWhite, color.BgRed).PrintfFunc()
	okColor    = color.New(color.FgWhite, color.BgGreen).PrintfFunc()
)

func main() {
	log.SetFlags(0)

	var (
		title       = flag.String("t", "", "pull request title")
		description = flag.String("d", "", "pull request description")
		into        = flag.String("i", "master", "target branch (optional)")
		from        = flag.String("f", "", "source branch (optional) (default = current branch)")
		intoRemote  = flag.String("into-remote", "origin", "target remote (optional)")
		fromRemote  = flag.String("from-remote", "origin", "source remote (optional)")
	)
	flag.Parse()

	var err error

	// Check flags
	if *title == "" {
		koColor("%s\n", noTitleErr)
		os.Exit(1)
	}

	if *description == "" {
		koColor("%s\n", noDescErr)
		os.Exit(1)
	}

	if *into == "" {
		koColor("target branch can't be null\n")
		os.Exit(1)
	}

	if *from == "" {
		*from, err = getBranchName()
		if err != nil {
			koColor("%s\n", err)
			os.Exit(1)
		}
	}

	if *intoRemote == "" {
		koColor("target remote can't be null\n")
		os.Exit(1)
	}

	if *fromRemote == "" {
		koColor("source remote can't be null\n")
		os.Exit(1)
	}

	// Check if "git" is present
	if _, err := exec.LookPath("git"); err != nil {
		koColor("%s\n", err)
		os.Exit(1)
	}

	// Check repository
	if _, err := execCommand("git", "status"); err != nil {
		koColor("%s\n", err)
		os.Exit(1)
	}

	// Get remotes
	remotes, err := getRemoteServers()
	if err != nil {
		koColor("%s\n", err)
		os.Exit(1)
	}

	if err := checkRemote(remotes, *intoRemote); err != nil {
		koColor("%s\n", err)
		os.Exit(1)
	}

	if err := checkRemote(remotes, *fromRemote); err != nil {
		koColor("%s\n", err)
		os.Exit(1)
	}

	// Check credentials
	var c *credentials
	if !hasCredentials() {
		c, err = getCredentials()
	} else {
		c, err = readCredentials()
	}
	if err != nil {
		koColor("%s\n", err)
		os.Exit(1)
	}

	opts := pullRequestOpts{
		credentials: c,
		fromBranch:  *from,
		fromOwner:   remotes[*fromRemote].owner,
		fromRepo:    remotes[*fromRemote].repo,
		intoBranch:  *into,
		intoOwner:   remotes[*intoRemote].owner,
		intoRepo:    remotes[*intoRemote].repo,
		description: *description,
		title:       *title,
	}

	resp, err := openPullRequest(opts)
	if err != nil {
		koColor("%s\n", err)
		os.Exit(1)
	}
	okColor("Pull request opened  %s\n", resp.HTML_URL)
}

func execCommand(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	// cmd.Path = path

	var so bytes.Buffer
	cmd.Stdout = &so

	var se bytes.Buffer
	cmd.Stderr = &se

	err := cmd.Run()
	if err != nil {
		log.Printf("%s\n", se.String())
		return nil, err
	}

	return so.Bytes(), nil
}

func checkRemote(remotes map[string]*remote, name string) error {
	if _, ok := remotes[name]; !ok {
		return errors.New(fmt.Sprintf("unknown remote %s", name))
	}
	return nil
}

func getRemoteServers() (map[string]*remote, error) {

	out, err := execCommand("git", "remote", "-v")
	if err != nil {
		return nil, err
	}

	remotes := make(map[string]*remote)
	lineSep := []byte("\n")
	pat := regexp.MustCompile(`(\S+)(\s|\t)+(git@github.com:|https:\/\/github.com\/)(\S+)(\/)(\S+)(\.git\s+)(\()(fetch)(\))`)
	res := bytes.Split(out, lineSep)
	for i := 0; i < len(res); i++ {
		r := pat.FindStringSubmatch(string(res[i]))
		if len(r) == 0 {
			continue
		}

		if _, ok := remotes[r[1]]; !ok {
			remotes[r[1]] = &remote{name: r[1], owner: r[4], repo: r[6]}
		}
	}

	return remotes, nil
}

func getBranchName() (string, error) {

	out, err := execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	f := bytes.TrimRight(out, "\n")
	return string(f), nil
}
