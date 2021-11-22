package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/google/go-github/v40/github"
	"github.com/spf13/cobra"
)

var pr bool
var outFileName string
var verboseSearch bool
var ascIssuesOrder bool

var searchCmd = &cobra.Command{
	Use:   "search QUERY",
	Short: "Search through GitHub using its search API",
	RunE: func(_ *cobra.Command, args []string) error {
		ctx := context.Background()
		token, found := os.LookupEnv("GITHUB_TOKEN")
		if !found {
			return fmt.Errorf("expected GITHUB_TOKEN environment variable")
		}
		tp := github.BasicAuthTransport{
			Username: "RedDocMD",
			Password: token,
		}
		client := github.NewClient(tp.Client())

		query := args[0]
		allIssues, err := githubSearchRust(ctx, client, query)
		if err != nil {
			return err
		}
		fmt.Printf("%d issues returned\n", len(allIssues))
		if verboseSearch {
			for _, issue := range allIssues {
				fmt.Printf("%s => %s\n", issue.GetTitle(), issue.GetHTMLURL())
			}
		}

		if outFileName != "" {
			fmt.Printf("Saving to %s ...\n", outFileName)
			issueJson, err := json.MarshalIndent(allIssues, "", "  ")
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(outFileName, issueJson, fs.ModePerm)
			if err != nil {
				return err
			}
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

func initSearchCmd() {
	searchCmd.Flags().BoolVar(&pr, "pr", false, "search for pull-requests instead of issues")
	searchCmd.Flags().StringVarP(&outFileName, "output", "o", "", "file to dump search results to (JSON)")
	searchCmd.Flags().BoolVarP(&verboseSearch, "verbose", "v", false, "print verbose output")
	searchCmd.Flags().BoolVar(&ascIssuesOrder, "asc", false, "order searches in ascending order of date")
}

func githubSearchRust(ctx context.Context, client *github.Client, query string) ([]*github.Issue, error) {
	var order string
	if ascIssuesOrder {
		order = "asc"
	} else {
		order = "desc"
	}
	opts := &github.SearchOptions{
		Sort:  "created",
		Order: order,
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}
	var queryString string
	if !pr {
		queryString = fmt.Sprintf("%s is:issue language:rust", query)
	} else {
		queryString = fmt.Sprintf("%s is:pull-request language:rust", query)
	}
	if verboseSearch {
		fmt.Printf("Query => \"%s\"\n", queryString)
	}

	var allIssues []*github.Issue
	for {
		issues, resp, err := client.Search.Issues(ctx, queryString, opts)
		if err != nil {
			return nil, err
		}
		if opts.Page == 1 {
			fmt.Printf("%d issues found\n", issues.GetTotal())
		}
		if issues.GetIncompleteResults() {
			fmt.Println("Results are incomplete!")
		}
		if verboseSearch {
			fmt.Printf("Fetched page %d with %d results\n", opts.Page, len(issues.Issues))
		}
		allIssues = append(allIssues, issues.Issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allIssues, nil
}
