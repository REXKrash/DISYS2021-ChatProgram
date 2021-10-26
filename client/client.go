package main

import (
	"bufio"
	"context"
	"log"
	"math"
	"os"
	"time"

	pb "chat-program/routeguide"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)

	var userName string
	if len(os.Args) > 0 {
		userName = os.Args[1]
	} else {
		log.Println("Please enter your username:")
		sc.Scan()
		userName = sc.Text()
	}

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
	timestamp := 0
	go func() {
		for {
			sc.Scan()
			if err := stream.SendMsg(&pb.Message{Sender: userName, Message: sc.Text(), Timestamp: int32(timestamp + 1)}); err != nil {
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
		timestamp = int(math.Max(float64(timestamp), float64(serverMessage.Timestamp))) + 1

		log.Println(serverMessage.Sender+":", serverMessage.Message, "- timestamp:", serverMessage.Timestamp)
	}
}
