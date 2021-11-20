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

	githubSearchApi(ctx, client)
}

func githubSearchApi(ctx context.Context, client *github.Client) {
	opts := &github.SearchOptions{Sort: "created", Order: "desc"}
	query := "deadlock is:issue language:rust"
	issues, _, err := client.Search.Issues(ctx, query, opts)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Query => ", query)
	fmt.Printf("%d issues found\n", issues.GetTotal())
	if issues.GetIncompleteResults() {
		fmt.Println("Results are incomplete!")
	}
	for _, issue := range issues.Issues {
		fmt.Printf("%s => %s\n", issue.GetHTMLURL(), issue.GetTitle())
	}
}
