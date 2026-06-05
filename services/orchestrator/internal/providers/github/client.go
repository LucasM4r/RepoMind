package github

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/LucasM4r/repomind/internal/providers"
)

type client struct {
	Token      string
	HTTPClient *http.Client
}

func NewClient(token string) providers.Fetcher {
	return &client{
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

func (c *client) FetchRepoZip(ctx context.Context, owner, repo string) (io.ReadCloser, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/main", owner, repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to fetch repo zip: %s", resp.Status)
	}

	return resp.Body, nil
}
