package gitutil

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/k0kubun/pp/v3"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	var viewerQuery struct {
		Viewer struct {
			Login githubv4.String
		}
	}

	NewClient(os.Getenv("access_token")).api.Query(context.Background(), &viewerQuery, nil)
	pp.Println(viewerQuery)

	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("access_token")},
	))
	client := github.NewClient(httpClient)
	pp.Println(client.Users.Get(context.Background(), ""))
}
