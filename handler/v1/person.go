package handler

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/finest08/PubSubPublisher/gen/proto/go/proto/person/v1"
)

type PersonServer struct {
	Dapr dapr.Client
	pb.UnimplementedPersonServiceServer
	
}

func (p PersonServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	person := req.Person

	// publish event
	if err := p.Dapr.PublishEvent(
		context.Background(),
		"pubsub-publish", "mytopic", person,
		dapr.PublishEventWithContentType("application/json"),
	); err != nil {
		return &pb.CreateResponse{}, status.Errorf(codes.Aborted, "%s", "error publishing event")
	}
	return &pb.CreateResponse{Person: person}, nil
}
