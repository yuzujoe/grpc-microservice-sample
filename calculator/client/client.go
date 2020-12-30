package main

import (
	"context"
	"fmt"
	"gihub.com/yuzujoe/grpc-microservice-sample/calculator/pb"
	"google.golang.org/grpc"
	"log"
)

func main()  {
	fmt.Println("Start Calc...")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Falied connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewCalculatorServiceClient(conn)
	doCalc(c)
}

func doCalc(c pb.CalculatorServiceClient)  {
	fmt.Println("starting to calc RPC...")
	req := &pb.CalculatorRequest{
		Calculator: &pb.Calculator{
			Sum1: 3,
			Sum2: 10,
		},
	}
	res, err := c.Calc(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calc RPC: %v", err)
	}

	log.Printf("Response from Calc: %v", res.Result)
}
