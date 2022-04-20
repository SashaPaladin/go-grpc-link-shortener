package main

import (
	"context"
	"google.golang.org/grpc"
	pb "link-shortener/gen/proto"
	"log"
)

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connected: %v", err)
	}

	client := pb.NewUrlManagementClient(conn)

	respCreate, err := client.Create(context.Background(), &pb.OriginalUrl{Url: "http://originalurl"})
	if err != nil {
		log.Fatalf("Failed to created: %v", err)
	}

	log.Printf("Create method: %v", respCreate)

	respGet, err := client.Get(context.Background(), &pb.ShortUrl{Url: respCreate.Url})
	if err != nil {
		log.Fatalf("Failed to geted: %v", err)
	}

	log.Printf("Get method: %v", respGet)
}
