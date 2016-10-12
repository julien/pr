package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var credentialsPath = os.Getenv("HOME") + "/.pr-token-cache"

func hasCredentials() bool {
	info, err := os.Stat(credentialsPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func getCredentials() (*credentials, error) {
	var username, password string

	fmt.Printf("Please enter your GitHub username: ")
	if _, err := fmt.Scanf("%s", &username); err != nil {
		return nil, err
	}

	fmt.Printf("Please enter your GitHub password: ")
	if _, err := fmt.Scanf("%s", &password); err != nil {
		return nil, err
	}

	if err := requestToken(username, password); err != nil {
		return nil, err
	}

	c := &credentials{username, password}

	if err := saveCredentials(c); err != nil {
		return nil, err
	}

	return c, nil
}

func saveCredentials(c *credentials) error {
	f, err := os.Create(credentialsPath)
	if err != nil {
		fmt.Printf("error: %s\n")
		return err

	}
	defer f.Close()

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f.Write(b)

	return nil
}

func readCredentials() (*credentials, error) {
	data, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}

	var c credentials
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil

}
