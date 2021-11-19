package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v40/github"
)

func main() {
	fmt.Println("Hello, World!")
	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	orgs, _, err := client.Organizations.List(context.Background(), "willnorris", nil)
	if err != nil {
		log.Println(err)
	}
	for _, org := range orgs {
		fmt.Println(*org)
	}
}
