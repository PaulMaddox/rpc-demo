package main

//go:generate protoc --proto_path=models --twirp_out=twitterarchive --go_out=twitterarchive models/twitterarchive.proto

import (
	"fmt"
	"net/http"

	"github.com/paulmaddox/rpc-demo/processor/twitterarchive"
)

func main() {

	// Messages
	fmt.Printf("Mounting http://localhost:8080%s\n", twitterarchive.TwitterArchivePathPrefix)
	mux := http.NewServeMux()

	archiveHandler := twitterarchive.NewTwitterArchiveServer(&twitterarchive.TwitterArchiveServer{}, nil)
	mux.Handle(twitterarchive.TwitterArchivePathPrefix, archiveHandler)

	http.ListenAndServe(":8080", mux)

}
