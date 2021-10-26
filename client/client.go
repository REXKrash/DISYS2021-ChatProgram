package main

import (
	"context"
	"log"
	"time"

	pb "chat-program/routeguide"

	uuid "github.com/nu7hatch/gouuid"
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

	ctx := context.Background()

	log.Println("Sending request to join chat room...")
	uuid, _ := uuid.NewV4()
	stream, err := client.ChatMessage(ctx)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	go func() {
		for {
			if err := stream.SendMsg(&pb.Message{Sender: uuid.String(), Message: "Hello there"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	for {
		serverMessage, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive from server: %v", err)
		}
		log.Println(serverMessage.Sender+":", serverMessage.Message)
	}
}
