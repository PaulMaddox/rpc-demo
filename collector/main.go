package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/paulmaddox/rpc-demo/archiver/twitterarchive"
)

var (
	// SearchTerm defines the Twitter search phrase subscribe to
	SearchTerm = os.Getenv("SEARCH_TERM")

	// ArchiveEndpoint is the URL of the archiving service
	ArchiveEndpoint = os.Getenv("ARCHIVE_ENDPOINT")
)

func main() {

	fmt.Printf("SEARCH_TERM: %s\n", SearchTerm)
	fmt.Printf("ARCHIVE_ENDPOINT: %s\n", ArchiveEndpoint)
	fmt.Printf("AWS_REGION: %s\n", os.Getenv("AWS_REGION"))
	fmt.Printf("AWS_ACCESS_KEY_ID: %s\n", os.Getenv("AWS_ACCESS_KEY_ID"))
	fmt.Printf("AWS_SECRET_ACCESS_KEY: %s\n", os.Getenv("AWS_SECRET_ACCESS_KEY"))

	// Search twitter for the term provided by the env var SEARCH_TERM
	params := &twitter.StreamFilterParams{
		Track:         []string{SearchTerm},
		StallWarnings: twitter.Bool(true),
	}

	fmt.Printf("Fetching Twitter authentication tokens from AWS SSM Parameter Store\n")
	auth, err := GetTwitterAuthDetails()
	if err != nil {
		fmt.Printf("Could not retrive Twitter authentication tokens: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Subscribing to tweets mentioning '%s'\n", params.Track)
	config := oauth1.NewConfig(auth.ConsumerKey, auth.ConsumerSecret)
	token := oauth1.NewToken(auth.Token, auth.Secret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	stream, err := client.Streams.Filter(params)
	if err != nil {
		fmt.Printf("Error creating Twitter stream: %s\n", err)
		os.Exit(1)
	}

	// Send the tweet to the archiver service
	ctx := context.Background()
	archiver := twitterarchive.NewTwitterArchiveProtobufClient(ArchiveEndpoint, &http.Client{})

	for msg := range stream.Messages {
		if tweet, ok := msg.(*twitter.Tweet); ok {

			_, err := archiver.Create(ctx, &twitterarchive.CreateRequest{
				Name:    tweet.User.Name,
				Message: tweet.Text,
			})

			if err != nil {
				fmt.Printf("WARNING: Failed to archive tweet: %s\n", err)
			}

		}
	}

}
