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
	"math"
	"net"
)

type server struct{}

func (*server) Calc(ctx context.Context, req *pb.CalculatorRequest) (*pb.CalculatorResponse, error) {
	fmt.Println("Calculator function was invoked with %v", req)
	sum1 := req.GetCalculator().GetSum1()
	sum2 := req.GetCalculator().GetSum2()
	res := &pb.CalculatorResponse{
		Result: sum1 + sum2,
	}
	return res, nil
}

func (*server) Decomposition(req *pb.PrimeNumberDecompositionRequest, stream pb.CalculatorService_DecompositionServer) error {
	fmt.Println("Decomposition function was invoked with %v", req)
	number := req.GetNumber()
	k := int64(2)
	for number > 1 {
		if number%k == 0 {
			stream.Send(&pb.PrimeNumberDecompositionResponse{
				Result: k,
			})
			number = number / k
			fmt.Println(number)
		} else {
			k++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream pb.CalculatorService_ComputeAverageServer) error {
	sum := int64(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err != nil {
			result := float64(sum) / float64(count)
			return stream.SendAndClose(&pb.ComputeAverageResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
	return nil
}

func (*server) FindMaximum(stream pb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request\n")
	maximum := int64(0)
	for  {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reding client stream: %v", err)
			return err
		}
		number := req.GetNumber()
		if number > maximum {
			maximum = number
			err := stream.Send(&pb.FindMaximumResponse{
				Result: maximum,
			})
			if err != nil {
				log.Fatal("Error while sending data to client: %v", err)
				return err
			}
		}
	}
	return nil
}

func (*server) SquareRoot(ctx context.Context, req *pb.SquareRootRequest) (*pb.SquareRootResponse, error)  {
 	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number %v", number),
		)
	}
	return &pb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Server Started...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
