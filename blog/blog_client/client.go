package main

import (
	"context"
	"log"

	"../blogpb"
	"google.golang.org/grpc"
)

func createBlogRequest(c blogpb.BlogServiceClient) {
	createBlogRequest := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			Title:    "Deepak",
			AuthorId: "100",
			Content:  "Some content",
		},
	}
	res, err := c.CreateBlog(context.Background(), createBlogRequest)
	if err != nil {
		log.Fatalf("Error while creating blog %v", err)
	}
	log.Printf("Blog created is %v\n", res.GetBlog())
	blogID := res.GetBlog().GetId()
	log.Printf("Blog id is %v\n", blogID)
}

func readBlogRequest(c blogpb.BlogServiceClient) {
	readBlogRequest := &blogpb.ReadBlogRequest{
		BlogId: "5cb7551381a0418608db494b",
	}
	res, err := c.ReadBlog(context.Background(), readBlogRequest)
	if err != nil {
		log.Fatalf("Error while reading blog %v", err)
	}
	log.Printf("Blog read is %v", res.Blog)
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failer to connect %v", err)
	}
	defer conn.Close()
	c := blogpb.NewBlogServiceClient(conn)
	createBlogRequest(c)
	readBlogRequest(c)
}
