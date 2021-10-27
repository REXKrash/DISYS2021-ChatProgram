package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	pb "chat-program/routeguide"

	uuid "github.com/nu7hatch/gouuid"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

var user pb.User

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
	uuid, _ := uuid.NewV4()
	var timestamp int32
	user = pb.User{Uuid: uuid.String(), Name: userName, Timestamp: timestamp}
	stream, err := client.JoinChatServer(ctx, &user)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	SetupCloseHandler(client)
	go func() {
		for {
			sc.Scan()
			var msg = sc.Text()
			if len(msg) > 0 && len(msg) <= 128 {
				_, err := client.SendMessage(ctx, &pb.Message{Sender: userName, Message: msg, Timestamp: timestamp + 1})
				if err != nil {
					log.Fatalln("Failed to send message")
				}
			} else {
				log.Println("Message must be between 1-128 characters")
			}
		}
	}()
	for {
		serverMessage, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive from server: %v", err)
		}
		timestamp = int32(math.Max(float64(timestamp), float64(serverMessage.Timestamp))) + 1

		log.Println(serverMessage.Sender+":", serverMessage.Message, "- timestamp:", serverMessage.Timestamp)
	}
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler(client pb.ChatServiceClient) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		client.LeaveChatServer(context.Background(), &user)
		os.Exit(0)
	}()
}
