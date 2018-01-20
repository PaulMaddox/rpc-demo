package main

//go:generate protoc --proto_path=models --twirp_out=twitterarchive --go_out=twitterarchive models/twitterarchive.proto

import (
	"fmt"
	"net/http"
	"os"

	"github.com/paulmaddox/rpc-demo/archiver/twitterarchive"
)

var (

	// Region is the AWS region to archive tweets in
	Region = os.Getenv("AWS_REGION")

	// KinesisStreamName is the ARN of the Amazon Kinesis Stream to archive tweets to
	KinesisStreamName = os.Getenv("KINESIS_STREAM_NAME")
)

func main() {

	// Start a local RPC server to listen for incoming tweets to archive
	fmt.Printf("Mounting http://localhost:8080%s\n", twitterarchive.TwitterArchivePathPrefix)
	mux := http.NewServeMux()

	archiver := twitterarchive.New(Region, KinesisStreamName)
	archiveHandler := twitterarchive.NewTwitterArchiveServer(archiver, nil)
	mux.Handle(twitterarchive.TwitterArchivePathPrefix, archiveHandler)

	http.ListenAndServe(":8080", mux)

}
