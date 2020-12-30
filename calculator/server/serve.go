package main

import (
	"context"
	"fmt"
	"gihub.com/yuzujoe/grpc-microservice-sample/calculator/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {}

func (*server) Calc(ctx context.Context, req *pb.CalculatorRequest) (*pb.CalculatorResponse, error) {
	fmt.Println("Calculator function was invoked with %v", req)
	sum1 := req.GetCalculator().GetSum1()
	sum2 := req.GetCalculator().GetSum2()
	res := &pb.CalculatorResponse{
		Result: sum1 + sum2,
	}
	return res, nil
}

func main()  {
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
