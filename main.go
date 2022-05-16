// Package main implements a client for Person service.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	pb "github.com/finest08/PubSubPublisher/gen/proto/go/proto/person/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type Person struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Occupation string `json:"occupation"`
	Age        string `json:"age"`
}

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	
	p = Person{
		FirstName: "Alice",
		LastName:"Springfield",
		Email:"spring@style.me",
		Occupation: "Hairstylist", 
		Age: "25",
	}
)

func main() {

	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		middleware.StripSlashes,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "QUERY"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
			// Debug:            true,
		}),
	)

	r.Route("/person", func(r chi.Router) {
		r.Post("/", p.Publish)
	})

	// start server ":"+os.Getenv("PORT")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Print(":"+os.Getenv("PORT"))
		fmt.Print(err)
	}
}


func (p *Person) Publish(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		return
	}

	defer r.Body.Close()
	reqByt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v", err)))
	}

	req := &pb.PersonRequest{}
	err = protojson.Unmarshal(reqByt, req)

	msg, err := p.SendProto(req)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("%s", msg.GetMessage())))
		// w.Write([]byte("done"))
	}
	
}

func (p *Person) SendProto(req *pb.PersonRequest) (*pb.PersonResponse, error) {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
