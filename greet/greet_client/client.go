package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"../greetpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failer to connect %v", err)
	}
	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	doManyTimes(c)
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
