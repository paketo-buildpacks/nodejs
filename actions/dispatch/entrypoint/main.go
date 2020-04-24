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

	webhook, err := ioutil.ReadFile(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("webhook -> \n%s\n", webhook)

	file, err := os.Open(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		fail(fmt.Errorf("failed to read $GITHUB_EVENT_PATH: %w", err))
	}

	var event struct {
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
		Release struct {
			TagName string `json:"tag_name"`
			Assets  []struct {
				BrowserDownloadURL string `json:"browser_download_url"`
			} `json:"assets"`
		} `json:"release"`
	}
	err = json.NewDecoder(file).Decode(&event)
	if err != nil {
		fail(fmt.Errorf("failed to decode $GITHUB_EVENT_PATH: %w", err))
	}

	fmt.Printf("Repository: %s\n", event.Repository.FullName)
	fmt.Printf("Tag: %s\n", event.Release.TagName)

	var dispatch struct {
		EventType     string `json:"event_type"`
		ClientPayload struct {
			Source string `json:"source"`
			URI    string `json:"uri"`
		} `json:"client_payload"`
	}

	dispatch.EventType = "update-buildpack-toml"
	dispatch.ClientPayload.Source = fmt.Sprintf("https://github.com/%s/archive/%s.tar.gz", event.Repository.FullName, event.Release.TagName)
	dispatch.ClientPayload.URI = event.Release.Assets[0].BrowserDownloadURL

	payloadData, err := json.Marshal(&dispatch)
	if err != nil {
		fail(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/repos/%s/dispatches", config.Endpoint, config.Repo), bytes.NewBuffer(payloadData))
	if err != nil {
		fail(fmt.Errorf("failed to create dispatch request: %w", err))
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.Token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fail(fmt.Errorf("failed to complete dispatch request: %w", err))
	}

	if resp.StatusCode != http.StatusNoContent {
		dump, _ := httputil.DumpResponse(resp, true)
		fail(fmt.Errorf("Error: unexpected response from dispatch request: %s", dump))
	}

	fmt.Println("Success!")
}

func fail(err error) {
	fmt.Printf("Error: %s", err)
	os.Exit(1)
}
