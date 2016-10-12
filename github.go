package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type authResponse struct {
	ID             int      `json:"id,omitempty"`
	URL            string   `json:"url,omitempty"`
	Token          string   `json:"token,omitempty"`
	HashedToken    string   `json:"hashed_token,omitempty"`
	TokenLastEight string   `json:"token_last_eight,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
}

type remote struct {
	name  string
	owner string
	repo  string
}

type pullRequestOpts struct {
	credentials *credentials
	fromBranch  string
	fromOwner   string
	fromRepo    string
	intoBranch  string
	intoOwner   string
	intoRepo    string
	title       string
	description string
}

type pullRequestResponse struct {
	ID       int    `json:"id,omitempty"`
	HTML_URL string `json:"html_url,omitempty"`
}

type githubError struct {
	Resource string `json:"resource,omitempty"`
	Code     string `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
}

type githubErrorResponse struct {
	Message string `json:"message,omitempty"`
	Errors  []githubError
}

func requestToken(username, password string) error {

	data := []byte(`{"note": "pr", "scopes": ["repo"]}`)

	req, err := http.NewRequest("POST", "https://api.github.com/authorizations", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "pr v1.0.0")
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var ar authResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return err
	}

	return nil
}

func openPullRequest(opts pullRequestOpts) (*pullRequestResponse, error) {

	if opts.credentials == nil {
		return nil, errors.New("missing credentials")
	}

	if opts.fromRepo != opts.intoRepo {
		return nil, errors.New(fmt.Sprintf("source repository (%s) doesn't match target repository (%s)"))
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", opts.intoOwner, opts.intoRepo)

	data := []byte(fmt.Sprintf(`{"head": "%s", "base": "%s", "body": "%s", "title": "%s"}`,
		opts.fromOwner+":"+opts.fromBranch,
		opts.intoBranch,
		opts.description,
		opts.title,
	))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "pr v1.0.0")
	req.SetBasicAuth(opts.credentials.Username, opts.credentials.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
		var ger githubErrorResponse
		err := json.Unmarshal(body, &ger)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(ger.Errors[0].Message)
	}

	var prr pullRequestResponse
	err = json.Unmarshal(body, &prr)
	if err != nil {
		return nil, err
	}

	return &prr, nil
}
