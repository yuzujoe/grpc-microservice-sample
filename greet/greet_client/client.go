package main

import (
	"fmt"
	"gihub.com/yuzujoe/grpc-microservice-sample/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main()  {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cloud not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Created client: %f", c)
}
