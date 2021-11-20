package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v40/github"
)

func main() {
	ctx := context.Background()
	client := github.NewClient(nil)

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
