package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"../calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("calculator invoked %v", req)
	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()
	result := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Starting server...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failer to listen %v", err)
	}
	s := grpc.NewServer()
	if err != nil {
		log.Fatalf("Failer to load server %v", err)
	}
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failer to listen %v", err)
	}
	fmt.Println("Starting listening in: 0.0.0.0:50051")
}
