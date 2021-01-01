package main

import (
	"context"
	"fmt"
	"gihub.com/yuzujoe/grpc-microservice-sample/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cloud not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	//doUnary(c)
	//doServerStreaming(c)

	//doClientStreaming(c)
	//doBiDiStreaming(c)

	doUnaryWithDeadLine(c, 5 * time.Second) // should complete
	doUnaryWithDeadLine(c, 1 * time.Second) //  should timeout
}

func doUnary(c greetpb.GreetServiceClient) {
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

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Mark",
		},
	}

	res, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetMany Times RPC: %v", err)
	}
	for {
		stream, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", stream.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	request := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephane",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Nancy",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tom",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("errro while reading stream: %v", err)
	}

	for _, req := range request {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient)  {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	request := []*greetpb.GreetEveryOneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephane",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Nancy",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tom",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
	}

	stream, err := c.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req :=  range request {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
		// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while reciving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()
	// block until everything is done
	<-waitc
}

func doUnaryWithDeadLine(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do an UnaryWithDeadLine RPC...")
	req := &greetpb.GreetWithDeadLineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maarek",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadLine(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline wax exceeded")
			} else {
				fmt.Printf("unexpected error: %v", err)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadLine RPC: %v", err)
		}
		return
	}
	log.Printf("Response from GreetWithDeadLine: %v", res.Result)
}
