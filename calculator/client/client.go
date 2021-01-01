package main

import (
	"context"
	"fmt"
	"gihub.com/yuzujoe/grpc-microservice-sample/calculator/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	//doComputeAverage(c)
	//doFindMaximum(c)

	doErrorUnary(c)
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

func doFindMaximum(c pb.CalculatorServiceClient)  {
	request := []*pb.FindMaximumRequest{
		{
			Number: 1,
		},
		{
			Number: 5,
		},
		{
			Number: 3,
		},
		{
			Number: 6,
		},
		{
			Number: 2,
		},
		{
			Number: 20,
		},
	}

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatal("Error while opening stream and calling FindMaximum: %v", err)
		return
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range request {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for  {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}

func doErrorUnary(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")

	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c, -2)
}

func doErrorCall(c pb.CalculatorServiceClient, n int32)  {
	res, err := c.SquareRoot(context.Background(), &pb.SquareRootRequest{
		Number: n,
	})
	if err != nil {
		respErr, ok :=status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
			return
		} else {
			log.Fatal("Big Error calling SquareRoot: %v", err)
			return
		}
	}
	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
}
