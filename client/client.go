package main

import (
	"context"
	"log"
	"time"

	uuid "github.com/nu7hatch/gouuid"

	pb "chat-program/routeguide"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer connection.Close()
	client := pb.NewChatServiceClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("Sending request to join chat room...")
	uuid, _ := uuid.NewV4()
	_, err = client.JoinServer(ctx, &pb.User{Name: "Tue", Uuid: uuid.String()})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
}
