// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/finest08/PubSubPublisher/gen/proto/go/proto/greeting/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Person struct {
	Name string
	Occupation string
	Age  string

}


const (
	defaultName = "big wide world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
	
	p = Person{
		Name:"Alice", 
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
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	per, err := c.Person(ctx, &pb.PersonRequest{Name: p.Name, Occupation: p.Occupation, Age: p.Age})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Subscribing service response: %s", r.GetMessage())
	log.Printf("Subscribing service response: %s", per.GetMessage())
}
