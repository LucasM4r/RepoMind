package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

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

	resp, err := c.fetch(ctx, url, strings.TrimSpace(c.Token))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized && strings.TrimSpace(c.Token) != "" {
		resp.Body.Close()
		resp, err = c.fetch(ctx, url, "")
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to fetch repo zip: %s", resp.Status)
	}

	return resp.Body, nil
}

func (c *client) fetch(ctx context.Context, url, token string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "repomind-orchestrator")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return c.HTTPClient.Do(req)
}
