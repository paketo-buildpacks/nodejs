package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	fmt.Println("Dispatching")

	var config struct {
		Endpoint string
		Repo     string
		Token    string
	}

	flag.StringVar(&config.Endpoint, "endpoint", "https://api.github.com", "Specifies endpoint for sending dispatch request")
	flag.StringVar(&config.Repo, "repo", "", "Specifies repo for sending dispatch request")
	flag.StringVar(&config.Token, "token", "", "Github Authorization Token")
	flag.Parse()

	file, err := os.Open(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		fail(fmt.Errorf("failed to read $GITHUB_EVENT_PATH: %w", err))
	}

	var event struct {
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
		Release struct {
			Name string `json:"name"`
		} `json:"release"`
	}
	err = json.NewDecoder(file).Decode(&event)
	if err != nil {
		fail(fmt.Errorf("failed to decode $GITHUB_EVENT_PATH: %w", err))
	}

	fmt.Printf("Repository: %s\n", event.Repository.FullName)
	fmt.Printf("Release: %s\n", event.Release.Name)

	var dispatch struct {
		EventType     string `json:"event_type"`
		ClientPayload struct {
			Repo    string `json:"repo"`
			Release string `json:"release"`
		} `json:"client_payload"`
	}

	dispatch.EventType = "update-buildpack-toml"
	dispatch.ClientPayload.Repo = event.Repository.FullName
	dispatch.ClientPayload.Release = event.Release.Name

	payloadData, err := json.Marshal(&dispatch)
	if err != nil {
		fail(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/repos/%s/dispatches", config.Endpoint, config.Repo), bytes.NewBuffer(payloadData))
	if err != nil {
		fail(fmt.Errorf("failed to create dispatch request: %w", err))
	}

	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("request ->\n%s\n", dump)

	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.Token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fail(fmt.Errorf("failed to complete dispatch request: %w", err))
	}

	dump, _ = httputil.DumpResponse(resp, true)
	fmt.Printf("response ->\n%s\n", dump)

	if resp.StatusCode != http.StatusNoContent {
		fail(fmt.Errorf("Error: unexpected response from dispatch request: %s", dump))
	}

	fmt.Println("Success!")
}

func fail(err error) {
	fmt.Printf("Error: %s", err)
	os.Exit(1)
}
