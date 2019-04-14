package main

import (
	"context"
	"fmt"
	"log"

	"../calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failer to connect %v", err)
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNumber:  4,
		SecondNumber: 2,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error is %v", err)
	}
	log.Printf("response is %v", res.SumResult)
}

func doComputeAverage(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting client streaming rpc...")
	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error is %v", err)
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

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SquareRootRequest{
		Number: 1,
	}
	res, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		error, ok := status.FromError(err)
		if ok {
			// actual error from gRpc (user error)
			fmt.Println(error.Message())
			fmt.Println(error.Code())
		} else {
			log.Fatalf("Big Error in SquareRoot %v", err)
		}
	}
	log.Printf("response is %v", res.String())
	req2 := &calculatorpb.SquareRootRequest{
		Number: -1,
	}
	res2, err := c.SquareRoot(context.Background(), req2)
	if err != nil {
		error, ok := status.FromError(err)
		if ok {
			// actual error from gRpc (user error)
			fmt.Println(error.Message())
			fmt.Println(error.Code())
		} else {
			log.Fatalf("Big Error in SquareRoot %v", err)
		}
	}
	log.Printf("response is %v", res2.String())
}
