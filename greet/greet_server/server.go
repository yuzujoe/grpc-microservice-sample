package main

import (
	"fmt"
	"gihub.com/yuzujoe/grpc-microservice-sample/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {}

func main()  {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listern: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
