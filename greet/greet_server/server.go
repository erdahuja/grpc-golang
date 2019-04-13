package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"../greetpb"
	"google.golang.org/grpc"
)

type server struct {
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("greet invoked %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes invoked %v\n", req.String())
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	for index := 0; index < 10; index++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: "Hello " + strconv.Itoa(index) + " " + firstName + " " + lastName,
		}
		stream.Send(res)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failer to listen %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	fmt.Println("Server listening on... 0.0.0.0:50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failer to listen %v", err)
	}
}
