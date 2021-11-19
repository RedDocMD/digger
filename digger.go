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

	opts := &github.SearchOptions{Sort: "created", Order: "asc"}
	query := "deadlock is:issue language:rust"
	issues, _, err := client.Search.Issues(ctx, query, opts)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%d issues found\n", issues.GetTotal())
	if issues.GetIncompleteResults() {
		fmt.Println("Results are incomplete!")
	}
	for _, issue := range issues.Issues {
		fmt.Printf("%s => %s\n", issue.GetHTMLURL(), issue.GetTitle())
	}
}
