package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-github/v40/github"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalln("Expected 1 arg: <outfilename>")
	}

	ctx := context.Background()
	token, found := os.LookupEnv("GITHUB_TOKEN")
	if !found {
		log.Fatalln("Expected GITHUB_TOKEN environment variable")
	}
	tp := github.BasicAuthTransport{
		Username: "RedDocMD",
		Password: token,
	}
	client := github.NewClient(tp.Client())

	allIssues, err := githubSearchApi(ctx, client)
	if err != nil {
		log.Fatalln(err)
	}
	for _, issue := range allIssues {
		fmt.Printf("%s => %s\n", issue.GetHTMLURL(), issue.GetTitle())
	}

	outFileName := os.Args[1]
	issueJson, err := json.MarshalIndent(allIssues, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(outFileName, issueJson, fs.ModePerm)
	if err != nil {
		log.Fatalln(err)
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
