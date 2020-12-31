package main

import (
	"context"
	"fmt"
	"gihub.com/yuzujoe/grpc-microservice-sample/calculator/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Start Calc...")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Falied connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewCalculatorServiceClient(conn)
	//doCalc(c)
	//doDecomposition(c)

	doComputeAverage(c)
}

func doCalc(c pb.CalculatorServiceClient) {
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

func doDecomposition(c pb.CalculatorServiceClient) {
	fmt.Println("starting to prime number decomposition")
	req := &pb.PrimeNumberDecompositionRequest{
		Number: 120,
	}
	res, err := c.Decomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition: %v", err)
	}
	for {
		stream, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from Decomposition: %v", stream.GetResult())
	}
}

func doComputeAverage(c pb.CalculatorServiceClient) {
	request := []*pb.ComputeAverageRequest{
		{
			Number: 1,
		},
		{
			Number: 2,
		},
		{
			Number: 3,
		},
		{
			Number: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while reading stream: %v", err)
	}

	for _, req := range request {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receving response from ComputeAcerage: %v", err)
	}
	fmt.Printf("ComputeAvarage Response: %v", res)
}
