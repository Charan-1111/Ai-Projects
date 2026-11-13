package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (c *Clients) CallGithub() {

	owner := os.Getenv("OWNER")
	repo := os.Getenv("REPO")
	filePath := os.Getenv("FILEPATH")

	// forming the api url
	apiUrl := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/%s",
		owner,
		repo,
		filePath,
	)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		apiUrl,
		nil,
	)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Accept", "application/vnd.github.raw+json")
	req.Header.Set("User-Agent", "repository-assistant")

	resp, err := c.clients.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("GitHub returned %s\n", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(string(body))
}
