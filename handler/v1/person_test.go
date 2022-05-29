package handler

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/finest08/PubSubPublisher/gen/proto/go/proto/person/v1"
)

// data structure for test cases (using generics for fun and teaching)
type testSpec[T proto.Message] struct {
	name string
	req  T
	code codes.Code
}

// start mock server, register our service, and return client connection
func mockConn() *grpc.ClientConn {
	return mockServer(func(srv *grpc.Server) {
		pb.RegisterPersonServiceServer(srv, &PersonServer{
		})
	})
}

// mock context metatdata
func mockContext() context.Context {
	ctx := context.Background()
	md := metadata.MD{
		"permissions":      []string{"read:person,write:person"},
		"x-token-c-tenant": []string{"12345678-1234-1234-1234-123456789012"},
	}
	return metadata.NewOutgoingContext(ctx, md)
}

func TestCreate(t *testing.T) {
	// create client connection to our server (handler)
	conn := mockConn()
	defer conn.Close()
	client := pb.NewPersonServiceClient(conn)

	// create test data
	tests := []testSpec[*pb.CreateRequest]{
		{
			name: "Valid",
			req: &pb.CreateRequest{
				Person: &pb.Person{
					FirstName: "Peter",
					LastName:   "Griffin",
					Email:  "pete@griffme.com",
				},	
			},
			code: codes.OK,
		},
		// {
		// 	name: "NoEmail",
		// 	req: &pb.CreateRequest{
		// 		Person: &pb.Person{
		// 			FirstName: "Peter",
		// 			LastName:   "Griffin",
		// 			Email:  "",
		// 		},
		// 	},
		// 	code: codes.InvalidArgument,
		// },
		// {
		// 	name: "NoFirstName",
		// 	req: &pb.CreateRequest{
		// 		Person: &pb.Person{
		// 			FirstName: "",
		// 			LastName:   "Griffin",
		// 			Email:  "pete@griffme.com",
		// 		},
		// 	},
		// 	code: codes.InvalidArgument,
		// },
		// {
		// 	name: "NoLastName",
		// 	req: &pb.CreateRequest{
		// 		Person: &pb.Person{
		// 			FirstName: "Peter",
		// 			LastName:   "",
		// 			Email:  "pete@griffme.com",
		// 		},
		// 	},
		// 	code: codes.InvalidArgument,
		// },
	}


// 	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rsp, err := client.Create(mockContext(), test.req)
			if err != nil {
				// get grpc status
				if s, ok := status.FromError(err); ok {
					if s.Code() != test.code {
						t.Fatalf("unexpected error code: %d, %s", s.Code(), s.Message())
					}
				}
				return 
			}

			// no validation errors - check response
			if rsp == nil {
				t.Fatal("response is nil")
			}
			if rsp.Message == "" {
				t.Fatal("Person not found in response")
			}
			// if rsp.Message != test.req.Person.FirstName {
			// 	t.Errorf("Check id is not uuid: %s", rsp.Check.Id)
			// }
			// if rsp.Check.CheckType != test.req.Check.CheckType {
			// 	t.Errorf("Check type is incorrect: %s", rsp.Check.CheckType)
			// }
			// if rsp.Check.Status != test.req.Check.Status {
			// 	t.Errorf("Check status is incorrect: %s", rsp.Check.Status)
			// }
		})
	}
}

// // ----------------------------------------------------------------------------
// // Mock gRPC server
// // ----------------------------------------------------------------------------
func mockServer(fn func(*grpc.Server)) *grpc.ClientConn {
	srv := grpc.NewServer()
	fn(srv)

	lis := bufconn.Listen(1024 * 1024)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	dialler := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialler))
	if err != nil {
		log.Fatal(err)
	}

	return conn
}