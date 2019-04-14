package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failer to connect %v", err)
	}
	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	doUnaryWithDeadline(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Deepak",
			LastName:  "Ahuja",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error is %v", err)
	}
	log.Printf("response is %v", res.Result)
}

func doManyTimes(c greetpb.GreetServiceClient) {
	fmt.Println("Starting streaming from client...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Deepak",
			LastName:  "Ahuja",
		},
	}
	res, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error is %v", err)
	}
	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream, %v", err)
		}
		log.Printf("response is %v", msg.String())
	}
}

func doLongGreet(c greetpb.GreetServiceClient) {
	fmt.Println("Starting client streaming rpc...")
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error is %v", err)
	}
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deepak",
				LastName:  "1",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deepak",
				LastName:  "2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deepak",
				LastName:  "3",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deepak",
				LastName:  "4",
			},
		},
	}
	for _, req := range requests {
		fmt.Println("Sending... ", req.String())
		stream.Send(req)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error is %v", err)
	}
	fmt.Printf("Response of LongGreet from server is %v\n", response)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting client streaming rpc...")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream %v", err)
	}
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deepak",
				LastName:  "1",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deepak",
				LastName:  "2",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deepak",
				LastName:  "3",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deepak",
				LastName:  "4",
			},
		},
	}

	waitc := make(chan struct{})
	go func() {
		for idx, req := range requests {
			fmt.Println("Sending message ", idx, req.String())
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving %v", err)
				close(waitc)
				break
			}
			fmt.Printf("Received %v\n", response.String())
		}
	}()
	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetWithDeadLineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Deepak",
			LastName:  "Ahuja",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := c.GreetWithDeadLine(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalf("Deadline Exceeded")
			} else {
				log.Fatalf("Unexpected error %v", statusErr)
			}
		} else {
			log.Fatalf("Error is %v", statusErr.Message())
		}
	}
	log.Printf("response is %v", res.Result)
}
