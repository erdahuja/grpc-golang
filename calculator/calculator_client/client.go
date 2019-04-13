package main

import (
	"context"
	"log"

	"../calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failer to connect %v", err)
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	doUnary(c)
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
