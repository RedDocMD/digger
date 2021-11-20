package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	token, found := os.LookupEnv("GITHUB_TOKEN")
	if !found {
		log.Fatalln("Expected GITHUB_TOKEN environment variable")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	allIssues, err := githubSearchApi(ctx, client)
	if err != nil {
		log.Fatalln(err)
	}
	for _, issue := range allIssues {
		fmt.Printf("%s => %s\n", issue.GetHTMLURL(), issue.GetTitle())
	}
}

func githubSearchApi(ctx context.Context, client *github.Client) ([]*github.Issue, error) {
	opts := &github.SearchOptions{
		Sort:  "created",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}
	query := "deadlock is:issue language:rust"
	fmt.Println("Query => ", query)

	var allIssues []*github.Issue
	for {
		issues, resp, err := client.Search.Issues(ctx, query, opts)
		if err != nil {
			return nil, err
		}
		if opts.Page == 1 {
			fmt.Printf("%d issues found\n", issues.GetTotal())
		}
		if issues.GetIncompleteResults() {
			fmt.Println("Results are incomplete!")
		}
		fmt.Printf("Fetched page %d with %d results\n", opts.Page, len(issues.Issues))
		allIssues = append(allIssues, issues.Issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allIssues, nil
}
