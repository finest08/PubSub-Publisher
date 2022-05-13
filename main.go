// Package main implements a client for Person service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/finest08/PubSubPublisher/gen/proto/go/proto/person/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Person struct {
	FirstName string
	LastName  string
	Email	 string
	Occupation string
	Age  string
}

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	
	p = Person{
		FirstName:"Alice",
		LastName:"Springfield",
		Email:"spring@style.me",
		Occupation: "Hairstylist", 
		Age: "25",
	}
)

func main() {
	
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPersonServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	per, err := c.Person(ctx, &pb.PersonRequest{FirstName: p.FirstName, LastName: p.LastName, Email: p.Email, Occupation: p.Occupation, Age: p.Age})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Subscribing service response: %s", per.GetMessage())
}
