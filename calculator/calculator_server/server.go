package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"../calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage invoked...")
	sum := int32(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream, %v", err)
			return err
		}
		number := req.GetNumber()
		sum += number
		count++
	}
	return nil
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot invoked %v", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Received a negative number")
	}
	res := &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
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
