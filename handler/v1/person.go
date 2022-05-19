package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/finest08/PubSubPublisher/gen/proto/go/proto/person/v1"
)

type PersonServer struct {
	pb.UnimplementedPersonServiceServer
}

func (p PersonServer) Person(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		return
	}
	defer r.Body.Close()

	byt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v", err)))
	}

	// unmarshal to protobuf
	
	req := &pb.PersonRequest{}
	err = protojson.Unmarshal(byt, req)

	rsp, err := p.SendProto(req)
	

	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("%s", rsp.GetMessage())))
	}

}

func (p PersonServer) SendProto(req *pb.PersonRequest) (*pb.PersonResponse, error) {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewPersonServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	per, err := c.Person(ctx, &pb.PersonRequest{FirstName: req.FirstName, LastName: req.LastName, Email: req.Email, Occupation: req.Occupation, Age: req.Age})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Subscribing service response: %s", per.GetMessage())

	return per, nil
}
