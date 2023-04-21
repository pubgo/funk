package gitutil

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Client is a GitHub client.
type Client struct {
	api *githubv4.Client
}

// NewClient creates a new GitHub client.
func NewClient(token string) *Client {
	var httpClient *http.Client
	httpClient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	client := githubv4.NewClient(httpClient)

	c := &Client{
		api: client,
	}

	return c
}

func (c *Client) queryWithRetry(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	if err := c.api.Query(ctx, q, variables); err != nil {
		if strings.Contains(err.Error(), "abuse-rate-limits") {
			time.Sleep(time.Minute)
			return c.queryWithRetry(ctx, q, variables)
		}

		return err
	}

	return nil
}
