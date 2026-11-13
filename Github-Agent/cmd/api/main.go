package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	owner := "Charan-1111"
	repo := "Rate-Limiter"
	filePath := "go.mod"

	apiURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/%s",
		owner,
		repo,
		filePath,
	)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/vnd.github.raw+json")
	req.Header.Set("User-Agent", "repository-assistant")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("GitHub returned %s", resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(content))
}
