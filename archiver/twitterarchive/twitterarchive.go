package twitterarchive

import (
	context "context"
	"fmt"

	"github.com/twitchtv/twirp"
)

type TwitterArchiveServer struct{}

func (t *TwitterArchiveServer) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {

	if len(req.Name) < 1 {
		return nil, twirp.InvalidArgumentError("name", "not set")
	}

	if len(req.Message) < 1 {
		return nil, twirp.InvalidArgumentError("message", "not set")
	}

	fmt.Printf("(archived) %s: %s\n", req.Name, req.Message)

	return &CreateResponse{Ok: true}, nil

}
