package main

import (
	"context"
	"fmt"
	"gihub.com/yuzujoe/grpc-microservice-sample/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cloud not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient)  {
	fmt.Println("Starting to do an Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maarek",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greeting RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}
