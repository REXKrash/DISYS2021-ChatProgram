package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	pb "chat-program/routeguide"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	userName := os.Args[1]

	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer connection.Close()
	client := pb.NewChatServiceClient(connection)

	ctx := context.Background()

	log.Println("Sending request to join chat room...")
	stream, err := client.ChatMessage(ctx)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	sc := bufio.NewScanner(os.Stdin)
	go func() {
		for {
			sc.Scan()
			if err := stream.SendMsg(&pb.Message{Sender: userName, Message: sc.Text()}); err != nil {
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
